package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type FspHdrPage struct {
	*BasePage
}

func (f *FspHdrPage) FSPHeader() *page.FSPHeader {
	return page.NewFSPHeader(f.pageBytes)
}

func (f *FspHdrPage) XDESEntry() []*page.XDESEntry {
	return page.NewXDESEntrys(f.pageBytes)
}

func (f *FspHdrPage) HexEditorTags() []*page.HexEditorTag {
	var tags []*page.HexEditorTag
	tags = append(tags, f.FilHeader().HexEditorTag()...)
	tags = append(tags, f.FSPHeader().HexEditorTag())

	xdesList := f.XDESEntry()
	for _, xdes := range xdesList {
		tags = append(tags, xdes.HexEditorTag())
	}

	tags = append(tags, f.FILTrailer().HexEditorTag())
	return tags
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
