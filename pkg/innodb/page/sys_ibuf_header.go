package page

import "encoding/json"

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

func (t *IBufHeaderPage) String() string {
	type Page struct {
		FILHeader  *FILHeader
		FsegEntry  *FsegEntry
		FILTrailer *FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:  t.FilHeader(),
		FsegEntry:  t.FsegEntry(),
		FILTrailer: t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
