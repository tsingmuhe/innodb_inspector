package page

type CompactRecordInfo struct {
	buf *Buf

	free bool

	RecordType       uint8
	NextRecordOffset int16
	Position         uint32
}

func NewCompactRecordInfo(position uint32, pageBytes []byte, free bool) *CompactRecordInfo {
	from := position - 5
	to := position
	buf := NewBuf(from, pageBytes[from:to], to-1)

	b := buf.Bytes(3)
	return &CompactRecordInfo{
		buf:              buf,
		free:             free,
		RecordType:       b[2] & 7,
		NextRecordOffset: buf.Int16(),
		Position:         position,
	}
}

func (t *CompactRecordInfo) NextPosition() uint32 {
	return uint32(int64(t.Position) + int64(t.NextRecordOffset))
}

func (t *CompactRecordInfo) HexEditorTag() *HexEditorTag {
	color := "lime"
	if t.free {
		color = "Chocolate"
	}
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   color,
		Caption: "CompactRecordInfo",
	}
}
