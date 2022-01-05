package page

const (
	RedundantSupremumPosition = RedundantInfimumPosition + RedundantInfimumSize
	RedundantSupremumSize     = 7 + 8
)

type RedundantSupremum struct {
	buf *Buf

	Supremum string
}

func NewRedundantSupremum(pageBytes []byte) *RedundantSupremum {
	from := RedundantSupremumPosition
	to := from + RedundantSupremumSize

	buf := NewBuf(from, pageBytes[from:to], to-1)
	buf.Bytes(7)
	return &RedundantSupremum{
		buf:      buf,
		Supremum: string(buf.Bytes(8)),
	}
}

func (t *RedundantSupremum) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "lime",
		Caption: "RedundantSupremum",
	}
}

const (
	CompactSupremumPosition = CompactInfimumPosition + CompactInfimumSize
	CompactSupremumSize     = 3 + 2 + 8
)

type CompactSupremum struct {
	buf *Buf

	RecordType       uint8
	NextRecordOffset int16
	Position         uint32
	Supremum         string
}

func NewCompactSupremum(pageBytes []byte) *CompactSupremum {
	from := CompactSupremumPosition
	to := from + CompactSupremumSize

	buf := NewBuf(from, pageBytes[from:to], to-1)

	buf.Byte()
	b := buf.Bytes(2)
	return &CompactSupremum{
		buf:              buf,
		RecordType:       b[1] & 7,
		NextRecordOffset: buf.Int16(),
		Position:         CompactSupremumPosition + 5,
		Supremum:         string(buf.Bytes(8)),
	}
}

func (t *CompactSupremum) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "lime",
		Caption: "CompactSupremum",
	}
}
