package page

const (
	InodeEntrySize         = 8 + 4 + 16 + 16 + 16 + 4 + (4 * 32)
	InodeEntryCountPerPage = 85
)

type InodeEntry struct {
	buf *Buf

	FsegId           uint64
	FsegNotFullNUsed uint32
	FsegFree         *FlstBaseNode
	FsegNotFull      *FlstBaseNode
	FsegFull         *FlstBaseNode
	FsegMagicN       uint32
	FsegFragArr      []uint32
}

func NewInodeEntries(pageBytes []byte) []*InodeEntry {
	from := FilHeaderSize + FlstNodeSize

	var inodeEntries []*InodeEntry
	for i := 0; i < InodeEntryCountPerPage; i++ {
		to := from + InodeEntrySize

		buf := NewBuf(from, pageBytes[from:to], to-1)

		inodeEntry := &InodeEntry{
			buf:              buf,
			FsegId:           buf.Uint64(),
			FsegNotFullNUsed: buf.Uint32(),
			FsegFree:         buf.FlstBaseNode(),
			FsegNotFull:      buf.FlstBaseNode(),
			FsegFull:         buf.FlstBaseNode(),
			FsegMagicN:       buf.Uint32(),
		}

		for j := 0; j < 32; j++ {
			pageNo := buf.Uint32()
			if pageNo != UndefinedPageNo {
				inodeEntry.FsegFragArr = append(inodeEntry.FsegFragArr, pageNo)
			}
		}

		if inodeEntry.FsegId > 0 {
			inodeEntries = append(inodeEntries, inodeEntry)
		}

		from = to
	}

	return inodeEntries

}

func (t *InodeEntry) HexEditorTag() *HexEditorTag {
	id := (t.buf.from - FilHeaderSize - FlstNodeSize) / InodeEntrySize

	color := "yellow"
	if id%2 == 1 {
		color = "lime"
	}

	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   color,
		Caption: "InodeEntry",
	}
}
