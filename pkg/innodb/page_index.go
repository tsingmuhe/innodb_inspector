package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type IndexPage struct {
	*BasePage
}

func (t *IndexPage) IndexHeader() *page.IndexHeader {
	return page.NewIndexHeader(t.pageBytes)
}

func (t *IndexPage) FSegHeader() *page.FSegHeader {
	return page.NewFSegHeader(t.pageBytes)
}

func (t *IndexPage) RedundantInfimum() *page.RedundantInfimum {
	return page.NewRedundantInfimum(t.pageBytes)
}

func (t *IndexPage) CompactInfimum() *page.CompactInfimum {
	return page.NewCompactInfimum(t.pageBytes)
}

func (t *IndexPage) RedundantSupremum() *page.RedundantSupremum {
	return page.NewRedundantSupremum(t.pageBytes)
}

func (t *IndexPage) CompactSupremum() *page.CompactSupremum {
	return page.NewCompactSupremum(t.pageBytes)
}

func (t *IndexPage) IsCompact() bool {
	indexHeader := t.IndexHeader()
	return indexHeader.IsCompact == 1
}

func (t *IndexPage) CompactRecordInfos() []*page.CompactRecordInfo {
	infimum := t.CompactInfimum()
	supremum := t.CompactSupremum()

	np := infimum.NextPosition()
	if np == supremum.Position {
		return nil
	}

	var compactRecordInfos []*page.CompactRecordInfo

	for np != supremum.Position {
		ri := page.NewCompactRecordInfo(np, t.pageBytes, false)
		compactRecordInfos = append(compactRecordInfos, ri)
		np = ri.NextPosition()
	}

	return compactRecordInfos
}

func (t *IndexPage) FreeCompactRecordInfos() []*page.CompactRecordInfo {
	indexheader := t.IndexHeader()
	np := uint32(indexheader.Free)
	if np == 0 {
		return nil
	}

	var compactRecordInfos []*page.CompactRecordInfo

	for {
		ri := page.NewCompactRecordInfo(np, t.pageBytes, true)
		compactRecordInfos = append(compactRecordInfos, ri)
		if ri.NextRecordOffset == 0 {
			break
		}
		np = ri.NextPosition()
	}

	return compactRecordInfos
}

func (t *IndexPage) HexEditorTags() []*page.HexEditorTag {
	var tags []*page.HexEditorTag
	tags = append(tags, t.FilHeader().HexEditorTag())
	tags = append(tags, t.IndexHeader().HexEditorTag())
	tags = append(tags, t.FSegHeader().HexEditorTag())

	if t.IsCompact() {
		tags = append(tags, t.CompactInfimum().HexEditorTag())
		tags = append(tags, t.CompactSupremum().HexEditorTag())

		compactRecordInfos := t.CompactRecordInfos()
		for _, compactRecordInfo := range compactRecordInfos {
			tags = append(tags, compactRecordInfo.HexEditorTag())
		}

		creeCompactRecordInfos := t.FreeCompactRecordInfos()
		for _, compactRecordInfo := range creeCompactRecordInfos {
			tags = append(tags, compactRecordInfo.HexEditorTag())
		}
	} else {
		tags = append(tags, t.RedundantInfimum().HexEditorTag())
		tags = append(tags, t.RedundantSupremum().HexEditorTag())
	}

	tags = append(tags, t.FILTrailer().HexEditorTag())
	return tags
}

func (t *IndexPage) String() string {
	type Page struct {
		FILHeader          *page.FILHeader
		IndexHeader        *page.IndexHeader
		FSegHeader         *page.FSegHeader
		Infimum            interface{}
		Supremum           interface{}
		CompactRecordInfos []*page.CompactRecordInfo
		FILTrailer         *page.FILTrailer
	}

	if t.IsCompact() {
		b, _ := json.MarshalIndent(&Page{
			FILHeader:          t.FilHeader(),
			IndexHeader:        t.IndexHeader(),
			FSegHeader:         t.FSegHeader(),
			Infimum:            t.CompactInfimum(),
			Supremum:           t.CompactSupremum(),
			CompactRecordInfos: t.CompactRecordInfos(),
			FILTrailer:         t.FILTrailer(),
		}, "", "  ")
		return string(b)
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   t.FilHeader(),
		IndexHeader: t.IndexHeader(),
		FSegHeader:  t.FSegHeader(),
		Infimum:     t.RedundantInfimum(),
		Supremum:    t.RedundantSupremum(),
		FILTrailer:  t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
