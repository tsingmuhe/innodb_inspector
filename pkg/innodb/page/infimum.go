package page

const (
	RedundantInfimumPosition = FilHeaderSize + IndexHeaderSize + FSegHeaderSize
	RedundantInfimumSize     = 7 + 8
)

type RedundantInfimum struct {
	buf *Buf

	Infimum string
}

func NewRedundantInfimum(pageBytes []byte) *RedundantInfimum {
	from := RedundantInfimumPosition
	to := from + RedundantInfimumSize

	buf := NewBuf(from, pageBytes[from:to], to-1)

	buf.Bytes(7)
	return &RedundantInfimum{
		buf:     buf,
		Infimum: string(buf.Bytes(8)),
	}
}

func (t *RedundantInfimum) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "lime",
		Caption: "RedundantInfimum",
	}
}

const (
	CompactInfimumPosition uint32 = FilHeaderSize + IndexHeaderSize + FSegHeaderSize
	CompactInfimumSize            = 3 + 2 + 8
)

type CompactInfimum struct {
	buf *Buf

	RecordType       uint8
	NextRecordOffset int16
	Position         uint32
	Infimum          string
}

func NewCompactInfimum(pageBytes []byte) *CompactInfimum {
	from := CompactInfimumPosition
	to := from + CompactInfimumSize

	buf := NewBuf(from, pageBytes[from:to], to-1)

	buf.Byte()
	b := buf.Bytes(2)
	return &CompactInfimum{
		buf:              buf,
		RecordType:       b[1] & 7,
		NextRecordOffset: buf.Int16(),
		Position:         CompactInfimumPosition + 5,
		Infimum:          string(buf.Bytes(8)),
	}
}

func (t *CompactInfimum) NextPosition() uint32 {
	return uint32(int64(t.Position) + int64(t.NextRecordOffset))
}

func (t *CompactInfimum) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "lime",
		Caption: "CompactInfimum",
	}
}
