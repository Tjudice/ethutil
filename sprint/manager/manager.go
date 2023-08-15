package manager

import (
	"context"
	"errors"
	"math/big"

	"gfx.cafe/open/ghost"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tjudice/ethutil/sprint"
	"github.com/upper/db/v4"
	"tuxpa.in/a/zlog/log"
)

type Event interface {
	// Copy should initialize a new event for populating
	Copy() Event
	// Process should utilize the event topics and data to populate the event's fields
	// furthermore, process should return an address or list of addresses to add to the
	// allowed address list
	Process(l FilterUpdater, raw *ghost.ErigonLog) (bool, error)
	// Table should return the table name for the event
	Table() string
	// Event Hash (topic0) Of Event
	EventHash() common.Hash
	// Log Index
	GetLogIndex() int
}

type EventUploader interface {
	GetActiveAddresses(ctx context.Context, d db.Session, block int) ([]common.Address, error)
	UploadEventSet(ctx context.Context, d db.Session, eventInfo []*EventInfo) error
	ValidateEventSet(ctx context.Context, d db.Session, startBlock, endBlock int, eventInfo []*EventInfo) (bool, error)
	DeleteEventRange(ctx context.Context, d db.Session, startBlock, endBlock int) error
}

var ErrDuplicateEvent = errors.New("cannot add event: duplicate event hash")

type Manager struct {
	poller         HeadTracker
	events         map[common.Hash]Event
	eventGroups    map[string]*EventGroup
	allowAddresses []common.Address
	uploader       EventUploader
}

type EventGroup struct {
	EventHashes []common.Hash
	Addresses   []common.Address
}

type HeadTracker interface {
	CurrentBlock() uint64
}

type FilterUpdater interface {
	// Adds a local address filter. This is used to filter out events that are not
	// relevant to the current sprint
	AddAddress(address common.Address)
}

func NewManager(poller HeadTracker, uploader EventUploader) *Manager {
	return &Manager{
		poller:      poller,
		events:      make(map[common.Hash]Event),
		eventGroups: make(map[string]*EventGroup),
		uploader:    uploader,
	}
}

func (m *Manager) GetStageFilters(startBlock, endBlock int) []ethereum.FilterQuery {
	fqs := make([]ethereum.FilterQuery, 0, len(m.eventGroups))
	for _, grp := range m.eventGroups {
		fqs = append(fqs, ethereum.FilterQuery{
			Addresses: grp.Addresses,
			Topics: [][]common.Hash{
				grp.EventHashes,
			},
			FromBlock: big.NewInt(int64(startBlock)),
			ToBlock:   big.NewInt(int64(endBlock)),
		})
	}
	return fqs
}

func (m *Manager) AddEventFilter(groupID string, addrs []common.Address, evs ...Event) error {
	if _, ok := m.eventGroups[groupID]; ok {
		return errors.New("cannot add event group: duplicate group id")
	}
	hashes := make([]common.Hash, 0, len(evs))
	for _, e := range evs {
		// ensures we dont already have an event with the same event sig
		if _, ok := m.events[e.EventHash()]; ok {
			return ErrDuplicateEvent
		}
		m.events[e.EventHash()] = e
		hashes = append(hashes, e.EventHash())
	}
	m.eventGroups[groupID] = &EventGroup{
		EventHashes: hashes,
		Addresses:   addrs,
	}
	return nil
}

func (m *Manager) AllowAddresses(addrs ...common.Address) {
	// these addresses wont be filtered out ever
	m.allowAddresses = append(m.allowAddresses, addrs...)
}

func (m *Manager) CurrentBlock() uint64 {
	// chain head, used for scheduling
	return m.poller.CurrentBlock()
}

func (m *Manager) Insert(ctx context.Context, s db.Session, startBlock, endBlock int, events *sprint.EventBatch) error {
	// gets all active addresses to check events against
	addrs, err := m.uploader.GetActiveAddresses(ctx, s, startBlock)
	if err != nil {
		return err
	}
	addrs = append(addrs, m.allowAddresses...)
	// filters out events that are not relevant to the current sprint
	filtered, err := m.filterQueriedEvents(addrs, events)
	if err != nil {
		return err
	}
	// uploads set of post-processed events in the same transaction
	if err := m.uploader.UploadEventSet(ctx, s, filtered); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Validate(ctx context.Context, s db.Session, startBlock, endBlock int, events *sprint.EventBatch) (bool, error) {
	// gets all active addresses to check events against
	// we must perform validation in the exact same way as we do
	// collecting, otherwise we cannot discern the validaty of a given range
	addrs, err := m.uploader.GetActiveAddresses(ctx, s, startBlock)
	if err != nil {
		return false, err
	}
	addrs = append(addrs, m.allowAddresses...)
	filtered, err := m.filterQueriedEvents(addrs, events)
	if err != nil {
		return false, err
	}
	// checks the newly filtered range (ideally a distance behind head that is safe from chain reorgianizations)
	// and returns whether or not the newly queried range is identicaly to the previously inserted range
	isValid, err := m.uploader.ValidateEventSet(ctx, s, startBlock, endBlock, filtered)
	if err != nil {
		return false, err
	}
	// ranges are equal, so don't reinsert
	if isValid {
		return false, nil
	}
	// validation fail, delete and reupload
	log.Debug().Int("start", startBlock).Int("end", endBlock).Msg("reorg detected, deleting and reuploading")
	err = m.uploader.DeleteEventRange(ctx, s, startBlock, endBlock)
	if err != nil {
		return false, err
	}
	return true, m.uploader.UploadEventSet(ctx, s, filtered)
}

type EventInfo struct {
	EventLog        Event
	TransactionInfo *sprint.TransactionInfo
	Block           int
	Timestamp       uint64
}

func (m *Manager) filterQueriedEvents(addrs []common.Address, events *sprint.EventBatch) ([]*EventInfo, error) {
	filterUpdater := newFilterTracker(addrs)
	out := make([]*EventInfo, 0, len(events.Events))
	for _, e := range events.Events {
		// paranoid check
		if len(e.Topics) == 0 {
			continue
		}
		// another paranoid check
		event, ok := m.events[e.Topics[0]]
		if !ok {
			continue
		}
		// checks if addresses is in set of valid event contract sources
		if _, ok := filterUpdater.addrs[e.Address]; !ok {
			continue
		}
		// user event inteface to create new instance to unmarshall log into
		newEvent := event.Copy()
		// check to make sure block data is present
		block, ok := events.Blocks[int(e.BlockNumber)]
		if !ok {
			return nil, errors.New("missing block")
		}
		// check to make sure tx data is present in block
		tx, ok := block.TxMap[e.TxHash]
		if !ok {
			return nil, errors.New("missing tx")
		}
		// set block timestamp, done this way so we dont have to align timestamp and other event info
		e.Timestamp = block.Timestamp.Uint64()
		// marshals log info into event struct, determines if the event is valid
		// and allows the event interface to update the addresses filter
		// in the case of a new contract event source being added intra-range
		ok, err := newEvent.Process(filterUpdater, e)
		if err != nil {
			return nil, err
		}
		// bad event, dont insert
		if !ok {
			continue
		}
		out = append(out, &EventInfo{
			EventLog:        newEvent,
			TransactionInfo: tx,
			Block:           int(e.BlockNumber),
			Timestamp:       block.Timestamp.Uint64(),
		})
	}
	return out, nil
}

func newFilterTracker(addrs []common.Address) *filterTracker {
	f := &filterTracker{
		addrs: make(map[common.Address]struct{}, len(addrs)),
	}
	for _, a := range addrs {
		f.addrs[a] = struct{}{}
	}
	return f
}

type filterTracker struct {
	addrs map[common.Address]struct{}
}

func (f *filterTracker) AddAddress(addr common.Address) {
	f.addrs[addr] = struct{}{}
}
