package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type XdesPage struct {
	*BasePage
}

func (f *XdesPage) XDESEntry() []*page.XDESEntry {
	c := f.PageCursorAt(150)

	var xdesEntries []*page.XDESEntry

	for i := 0; i < 256; i++ {
		xdesEntry := page.NewXDESEntry(uint32(i))
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

func (f *XdesPage) HexEditorTags() []*page.HexEditorTag {
	var result []*page.HexEditorTag
	result = append(result, f.FilHeader().HexEditorTag())

	xdesList := f.XDESEntry()
	for _, xdes := range xdesList {
		result = append(result, xdes.HexEditorTag())
	}

	result = append(result, f.FILTrailer().HexEditorTag(len(f.pageBits)))
	return result
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
