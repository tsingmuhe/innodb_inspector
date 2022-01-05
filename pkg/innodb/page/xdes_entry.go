package page

const (
	XDESStateFree     XDESState = 1
	XDESStateFreeFrag XDESState = 2
	XDESStateFullFrag XDESState = 3
	XDESStateFseg     XDESState = 4
)

type XDESState uint32

func (t XDESState) String() string {
	switch t {
	case XDESStateFree:
		return "free"
	case XDESStateFreeFrag:
		return "free_frag"
	case XDESStateFullFrag:
		return "full_frag"
	case XDESStateFseg:
		return "fseg"
	}

	return "unknown"
}

func (t XDESState) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

const (
	XDESEntrySize         = 8 + 12 + 4 + 16
	XDESEntryCountPerPage = 256
)

type XDESEntry struct {
	buf *Buf

	SegmentId uint64
	FlstNode  *FlstNode
	State     XDESState
	Bitmap    []byte
}

func NewXDESEntrys(pageBytes []byte) []*XDESEntry {
	from := FilHeaderSize + FSPHeaderSize

	var xdesEntries []*XDESEntry

	for i := 0; i < XDESEntryCountPerPage; i++ {
		to := from + XDESEntrySize

		buf := NewBuf(from, pageBytes[from:to], to-1)
		xdesEntry := &XDESEntry{
			buf:       buf,
			SegmentId: buf.Uint64(),
			FlstNode:  buf.FlstNode(),
			State:     XDESState(buf.Uint32()),
			Bitmap:    buf.Bytes(16),
		}

		if xdesEntry.State > 0 {
			xdesEntries = append(xdesEntries, xdesEntry)
		}

		from = to
	}

	return xdesEntries
}

func (t *XDESEntry) HexEditorTag() *HexEditorTag {
	id := (t.buf.from - FilHeaderSize - FSPHeaderSize) / XDESEntrySize
	color := "yellow"
	if id%2 == 1 {
		color = "lime"
	}

	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   color,
		Caption: "XDESEntry",
	}
}
