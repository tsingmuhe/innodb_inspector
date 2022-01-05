package page

const (
	InodeEntrySize         = 8 + 4 + 16 + 16 + 16 + 4 + (4 * 32)
	InodeEntryCountPerPage = 85
)

type InodeEntry struct {
	id uint32

	FsegId           uint64
	FsegNotFullNUsed uint32
	FsegFree         *FlstBaseNode
	FsegNotFull      *FlstBaseNode
	FsegFull         *FlstBaseNode
	FsegMagicN       uint32
	FsegFragArr      []uint32
}

func NewInodeEntry(id uint32) *InodeEntry {
	return &InodeEntry{
		id: id,
	}
}

func (t *InodeEntry) HexEditorTag() *HexEditorTag {
	from := FilHeaderSize + FlstNodeSize + t.id*InodeEntrySize

	color := "yellow"
	if t.id%2 == 1 {
		color = "lime"
	}

	return &HexEditorTag{
		From:    from,
		To:      from + InodeEntrySize - 1,
		Color:   color,
		Caption: "InodeEntry",
	}
}
