package sprint

type Status struct {
	Active          bool
	BlockHeight     int64
	CollectionQueue []BatchLogViewer
	ValidationQueue []BatchLogViewer
}
