package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type IndexPage struct {
	*BasePage
}

func (t *IndexPage) IndexHeader() *page.IndexHeader {
	c := t.CursorAtBodyStart()
	return &page.IndexHeader{
		NDirSlots:  c.Uint16(),
		HeapTop:    c.Uint16(),
		NHeap:      c.Uint16(),
		Free:       c.Uint16(),
		Garbage:    c.Uint16(),
		LastInsert: c.Uint16(),
		Direction:  c.Uint16(),
		NDirection: c.Uint16(),
		NRecs:      c.Uint16(),
		MaxTrxID:   c.Uint64(),
		Level:      c.Uint16(),
		IndexId:    c.Uint64(),
	}
}

func (t *IndexPage) FSegHeader() *page.FSegHeader {
	c := t.CursorAt(74)
	return &page.FSegHeader{
		Leaf: &page.FsegEntry{
			FsegHdrSpace:  c.Uint32(),
			FsegHdrPageNo: c.Uint32(),
			FsegHdrOffset: c.Uint16(),
		},
		NoLeaf: &page.FsegEntry{
			FsegHdrSpace:  c.Uint32(),
			FsegHdrPageNo: c.Uint32(),
			FsegHdrOffset: c.Uint16(),
		},
	}
}

func (t *IndexPage) Infimum() *page.Infimum {
	if t.IsCompact() {
		c := t.CursorAt(94)
		return &page.Infimum{
			InfoFlags: c.Bytes(5),
			Infimum:   string(c.Bytes(8)),
		}
	}

	c := t.CursorAt(94)
	return &page.Infimum{
		InfoFlags: c.Bytes(7),
		Infimum:   string(c.Bytes(8)),
	}
}

func (t *IndexPage) Supremum() *page.Supremum {
	if t.IsCompact() {
		c := t.CursorAt(107)
		return &page.Supremum{
			InfoFlags: c.Bytes(5),
			Supremum:  string(c.Bytes(8)),
		}
	}
	
	c := t.CursorAt(109)
	return &page.Supremum{
		InfoFlags: c.Bytes(7),
		Supremum:  string(c.Bytes(8)),
	}
}

func (t *IndexPage) IsCompact() bool {
	c := t.CursorAt(42)
	flag := c.Bytes(1)[0]
	return (flag >> 7) == 1
}

func (t *IndexPage) String() string {
	type Page struct {
		FILHeader   *page.FILHeader
		IndexHeader *page.IndexHeader
		FSegHeader  *page.FSegHeader
		Infimum     *page.Infimum
		Supremum    *page.Supremum
		FILTrailer  *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   t.FilHeader(),
		IndexHeader: t.IndexHeader(),
		FSegHeader:  t.FSegHeader(),
		Infimum:     t.Infimum(),
		Supremum:    t.Supremum(),
		FILTrailer:  t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
