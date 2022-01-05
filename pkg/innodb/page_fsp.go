package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type FspHdrPage struct {
	*BasePage
}

func (f *FspHdrPage) FSPHeader() *page.FSPHeader {
	c := f.PageCursorAtBodyStart()
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
	c := f.PageCursorAt(150)

	var xdesEntries []*page.XDESEntry

	for i := uint32(0); i < 256; i++ {
		xdesEntry := page.NewXDESEntry(i)
		xdesEntry.SegmentId = c.Uint64()
		xdesEntry.FlstNode = c.FlstNode()
		xdesEntry.State = page.XDESState(c.Uint32())
		xdesEntry.Bitmap = c.Bytes(16)

		if xdesEntry.State > 0 {
			xdesEntries = append(xdesEntries, xdesEntry)
		}
	}

	return xdesEntries
}

func (f *FspHdrPage) HexEditorTags() []*page.HexEditorTag {
	var result []*page.HexEditorTag
	result = append(result, f.FilHeader().HexEditorTag())
	result = append(result, f.FSPHeader().HexEditorTag())

	xdesList := f.XDESEntry()
	for _, xdes := range xdesList {
		result = append(result, xdes.HexEditorTag())
	}

	result = append(result, f.FILTrailer().HexEditorTag(len(f.pageBits)))
	return result
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
