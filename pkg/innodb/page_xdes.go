package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type XdesPage struct {
	*BasePage
}

func (t *XdesPage) XDESEntry() []*page.XDESEntry {
	return page.NewXDESEntrys(t.pageBytes)
}

func (t *XdesPage) HexEditorTags() []*page.HexEditorTag {
	var tags []*page.HexEditorTag
	tags = append(tags, t.FilHeader().HexEditorTag())

	xdesList := t.XDESEntry()
	for _, xdes := range xdesList {
		tags = append(tags, xdes.HexEditorTag())
	}

	tags = append(tags, t.FILTrailer().HexEditorTag())
	return tags
}

func (t *XdesPage) String() string {
	type Page struct {
		FILHeader   *page.FILHeader
		XDESEntries []*page.XDESEntry
		FILTrailer  *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   t.FilHeader(),
		XDESEntries: t.XDESEntry(),
		FILTrailer:  t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
