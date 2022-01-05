package page

const (
	FSPHeaderPosition uint32 = FilHeaderSize
	FSPHeaderSize     uint32 = 4 + 4 + 4 + 4 + 4 + 4 + 16 + 16 + 16 + 8 + 16 + 16
)

type FSPHeader struct {
	*Buf

	SpaceId       uint32
	Unused        uint32
	Size          uint32
	FreeLimit     uint32
	Flags         []byte `json:"-"`
	FreeFragNUsed uint32
	Free          *FlstBaseNode
	FreeFrag      *FlstBaseNode
	FullFrag      *FlstBaseNode
	NextSegId     uint64
	FullInodes    *FlstBaseNode
	FreeInodes    *FlstBaseNode
}

func NewFSPHeader(pageBytes []byte) *FSPHeader {
	from := FSPHeaderPosition
	to := FSPHeaderPosition + FSPHeaderSize

	buf := NewBuf(from, pageBytes[from:to], to-1)
	return &FSPHeader{
		Buf:           buf,
		SpaceId:       buf.Uint32(),
		Unused:        buf.Uint32(),
		Size:          buf.Uint32(),
		FreeLimit:     buf.Uint32(),
		Flags:         buf.Bytes(4),
		FreeFragNUsed: buf.Uint32(),
		Free:          buf.FlstBaseNode(),
		FreeFrag:      buf.FlstBaseNode(),
		FullFrag:      buf.FlstBaseNode(),
		NextSegId:     buf.Uint64(),
		FullInodes:    buf.FlstBaseNode(),
		FreeInodes:    buf.FlstBaseNode(),
	}
}

func (t *FSPHeader) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.from,
		To:      t.to,
		Color:   "orange",
		Caption: "FSPHeader",
	}
}
