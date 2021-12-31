package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type XdesPage struct {
	*BasePage
}

func (f *XdesPage) XDESEntry() []*page.XDESEntry {
	c := f.CursorAt(150)

	var xdesEntries []*page.XDESEntry

	for i := 0; i < 256; i++ {
		xdesEntry := &page.XDESEntry{
			SegmentId: c.Uint64(),
			FlstNode:  c.FlstNode(),
			State:     page.XDESState(c.Uint32()),
			Bitmap:    c.Bytes(16),
		}

		if xdesEntry.State > 0 {
			xdesEntries = append(xdesEntries, xdesEntry)
		}
	}

	return xdesEntries
}

func (f *XdesPage) String() string {
	type Page struct {
		FILHeader   *page.FILHeader
		XDESEntries []*page.XDESEntry
		FILTrailer  *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   f.FilHeader(),
		XDESEntries: f.XDESEntry(),
		FILTrailer:  f.FILTrailer(),
	}, "", "  ")
	return string(b)
}
