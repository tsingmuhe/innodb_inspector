package page

type RsegHeader struct {
	MaxSize     uint32
	HistorySize uint32
	HistoryList *FlstBaseNode
	FsegEntry   *FsegEntry
}

type SysRsegHeaderPage struct {
	*BasePage
}

func (t *SysRsegHeaderPage) RsegHeader() *RsegHeader {
	c := t.CursorAtBodyStart()
	return &RsegHeader{
		MaxSize:     c.Uint32(),
		HistorySize: c.Uint32(),
		HistoryList: c.FlstBaseNode(),
		FsegEntry: &FsegEntry{
			FsegHdrSpace:  c.Uint32(),
			FsegHdrPageNo: c.Uint32(),
			FsegHdrOffset: c.Uint16(),
		},
	}
}
