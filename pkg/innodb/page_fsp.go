package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type FspHdrPage struct {
	*BasePage
}

func (f *FspHdrPage) FSPHeader() *page.FSPHeader {
	c := f.CursorAtBodyStart()
	return &page.FSPHeader{
		SpaceId:       c.Uint32(),
		Unused:        c.Uint32(),
		Size:          c.Uint32(),
		FreeLimit:     c.Uint32(),
		Flags:         c.Bytes(4),
		FreeFragNUsed: c.Uint32(),
		Free:          c.FlstBaseNode(),
		FreeFrag:      c.FlstBaseNode(),
		FullFrag:      c.FlstBaseNode(),
		NextSegId:     c.Uint64(),
		FullInodes:    c.FlstBaseNode(),
		FreeInodes:    c.FlstBaseNode(),
	}
}

func (f *FspHdrPage) XDESEntry() []*page.XDESEntry {
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

func (f *FspHdrPage) String() string {
	type Page struct {
		FILHeader   *page.FILHeader
		FSPHeader   *page.FSPHeader
		XDESEntries []*page.XDESEntry
		FILTrailer  *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   f.FilHeader(),
		FSPHeader:   f.FSPHeader(),
		XDESEntries: f.XDESEntry(),
		FILTrailer:  f.FILTrailer(),
	}, "", "  ")
	return string(b)
}
