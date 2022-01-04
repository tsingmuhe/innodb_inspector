package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type SysRsegHeaderPage struct {
	*BasePage
}

func (t *SysRsegHeaderPage) RsegHeader() *page.RsegHeader {
	c := t.PageCursorAtBodyStart()
	return &page.RsegHeader{
		MaxSize:     c.Uint32(),
		HistorySize: c.Uint32(),
		HistoryList: c.FlstBaseNode(),
		FsegEntry: &page.FsegEntry{
			FsegHdrSpace:  c.Uint32(),
			FsegHdrPageNo: c.Uint32(),
			FsegHdrOffset: c.Uint16(),
		},
	}
}

func (t *SysRsegHeaderPage) String() string {
	type Page struct {
		FILHeader  *page.FILHeader
		RsegHeader *page.RsegHeader
		FILTrailer *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:  t.FilHeader(),
		RsegHeader: t.RsegHeader(),
		FILTrailer: t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
