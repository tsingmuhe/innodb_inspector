package page

type IBufHeaderPage struct {
	*BasePage
}

func (t *IBufHeaderPage) FsegEntry() *FsegEntry {
	c := t.CursorAtBodyStart()
	return &FsegEntry{
		FsegHdrSpace:  c.Uint32(),
		FsegHdrPageNo: c.Uint32(),
		FsegHdrOffset: c.Uint16(),
	}
}
