package page

import "encoding/binary"

const (
	IndexHeaderPosition = FilHeaderSize
	IndexHeaderSize     = 2 + 2 + 2 + 2 + 2 + 2 + 2 + 2 + 2 + 8 + 2 + 8
)

type IndexHeader struct {
	buf *Buf

	NDirSlots  uint16
	HeapTop    uint16
	IsCompact  uint8
	NHeap      uint16
	Free       uint16
	Garbage    uint16
	LastInsert uint16
	Direction  uint16
	NDirection uint16
	NRecs      uint16
	MaxTrxID   uint64
	Level      uint16
	IndexId    uint64
}

func NewIndexHeader(pageBytes []byte) *IndexHeader {
	from := FilHeaderSize
	to := from + IndexHeaderSize

	buf := NewBuf(from, pageBytes[from:to], to-1)

	indexHeader := &IndexHeader{
		buf: buf,
	}

	indexHeader.NDirSlots = buf.Uint16()
	indexHeader.HeapTop = buf.Uint16()

	flag := buf.Bytes(2)
	f1 := flag[0]
	indexHeader.IsCompact = f1 >> 7

	b1 := f1 & 127
	b2 := flag[1]
	indexHeader.NHeap = binary.BigEndian.Uint16([]byte{b1, b2})

	indexHeader.Free = buf.Uint16()
	indexHeader.Garbage = buf.Uint16()
	indexHeader.LastInsert = buf.Uint16()
	indexHeader.Direction = buf.Uint16()
	indexHeader.NDirection = buf.Uint16()
	indexHeader.NRecs = buf.Uint16()
	indexHeader.MaxTrxID = buf.Uint64()
	indexHeader.Level = buf.Uint16()
	indexHeader.IndexId = buf.Uint64()

	return indexHeader
}

func (t *IndexHeader) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "orange",
		Caption: "IndexHeader",
	}
}
