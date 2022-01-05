package innodb

import (
	"encoding/json"
	"fmt"
	"innodb_inspector/pkg/innodb/page"
	"os"
)

type IndexPage struct {
	*BasePage
}

func (t *IndexPage) IndexHeader() *page.IndexHeader {
	c := t.PageCursorAtBodyStart()
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
	c := t.PageCursorAt(74)
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

func (t *IndexPage) IsCompact() bool {
	c := t.PageCursorAt(42)
	flag := c.Byte()
	return (flag >> 7) == 1
}

func (t *IndexPage) Infimum() *page.CompactInfimum {
	if t.IsCompact() {
		c := t.PageCursorAt(94)
		bs := c.Bytes(3)
		recordType := bs[2]
		return &page.CompactInfimum{
			OffSet:           99,
			InfoFlags:        bs,
			RecordType:       uint(recordType & byte(7)),
			NextRecordOffset: c.Int16(),
			Infimum:          string(c.Bytes(8)),
		}
	}

	return nil
}

func (t *IndexPage) Supremum() *page.CompactSupremum {
	if t.IsCompact() {
		c := t.PageCursorAt(107)
		bs := c.Bytes(3)
		recordType := bs[2]
		return &page.CompactSupremum{
			OffSet:           112,
			InfoFlags:        bs,
			RecordType:       uint(recordType & byte(7)),
			NextRecordOffset: c.Int16(),
			Supremum:         string(c.Bytes(8)),
		}
	}
	return nil
}

func (t *IndexPage) Records() []uint32 {
	infimum := t.Infimum()
	supremum := t.Supremum()

	var result []uint32
	result = append(result, infimum.OffSet)
	result = append(result, supremum.OffSet)

	nr := infimum.NextRecord()
	if nr == supremum.OffSet {
		return result
	}

	for {
		result = append(result, nr)
		ir := t.IndexRecord(nr)
		if ir.NextRecord() == supremum.OffSet {
			break
		}
		nr = ir.NextRecord()
	}

	return result
}

func (t *IndexPage) IndexRecord(offset uint32) *page.IndexRecord {
	c := t.PageCursorAt(offset - 3)
	recordType := c.Byte()
	return &page.IndexRecord{
		OffSet:           offset,
		RecordType:       uint(recordType & byte(7)),
		NextRecordOffset: c.Int16(),
	}
}

func (t *IndexPage) String() string {
	fmt.Println(t.Records())

	file, _ := os.Create("output.bin")
	file.Write(t.pageBits)

	type Page struct {
		FILHeader   *page.FILHeader
		IndexHeader *page.IndexHeader
		FSegHeader  *page.FSegHeader
		Infimum     *page.CompactInfimum
		Supremum    *page.CompactSupremum
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
