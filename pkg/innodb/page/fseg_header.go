package page

const (
	FSegHeaderPosition = FilHeaderSize + IndexHeaderSize
	FSegHeaderSize     = FsegEntrySize + FsegEntrySize
)

type FSegHeader struct {
	buf *Buf

	Leaf   *FsegEntry
	NoLeaf *FsegEntry
}

func NewFSegHeader(pageBytes []byte) *FSegHeader {
	from := FilHeaderSize + IndexHeaderSize
	to := from + FSegHeaderSize

	buf := NewBuf(from, pageBytes[from:to], to-1)

	return &FSegHeader{
		buf: buf,

		Leaf: &FsegEntry{
			FsegHdrSpace:  buf.Uint32(),
			FsegHdrPageNo: buf.Uint32(),
			FsegHdrOffset: buf.Uint16(),
		},
		NoLeaf: &FsegEntry{
			FsegHdrSpace:  buf.Uint32(),
			FsegHdrPageNo: buf.Uint32(),
			FsegHdrOffset: buf.Uint16(),
		},
	}
}

func (t *FSegHeader) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "yellow",
		Caption: "FSegHeader",
	}
}
