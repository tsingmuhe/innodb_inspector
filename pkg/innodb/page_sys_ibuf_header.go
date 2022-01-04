package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type IBufHeaderPage struct {
	*BasePage
}

func (t *IBufHeaderPage) FsegEntry() *page.FsegEntry {
	c := t.PageCursorAtBodyStart()
	return &page.FsegEntry{
		FsegHdrSpace:  c.Uint32(),
		FsegHdrPageNo: c.Uint32(),
		FsegHdrOffset: c.Uint16(),
	}
}

func (t *IBufHeaderPage) String() string {
	type Page struct {
		FILHeader  *page.FILHeader
		FsegEntry  *page.FsegEntry
		FILTrailer *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:  t.FilHeader(),
		FsegEntry:  t.FsegEntry(),
		FILTrailer: t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
