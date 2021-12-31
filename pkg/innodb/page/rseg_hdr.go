package page

type RsegHeader struct {
	MaxSize     uint32
	HistorySize uint32
	HistoryList *FlstBaseNode
	FsegEntry   *FsegEntry
}
