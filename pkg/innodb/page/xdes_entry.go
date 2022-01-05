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
	id uint32

	SegmentId uint64
	FlstNode  *FlstNode
	State     XDESState
	Bitmap    Bits
}

func NewXDESEntry(id uint32) *XDESEntry {
	return &XDESEntry{
		id: id,
	}
}

func (t *XDESEntry) HexEditorTag() *HexEditorTag {
	from := FSPHeaderPosition + FSPHeaderSize + t.id*XDESEntrySize
	color := "yellow"
	if t.id%2 == 1 {
		color = "lime"
	}

	return &HexEditorTag{
		From:    from,
		To:      from + XDESEntrySize - 1,
		Color:   color,
		Caption: "XDESEntry",
	}
}
