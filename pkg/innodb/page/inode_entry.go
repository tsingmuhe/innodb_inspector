package page

const (
	InodeEntrySize         = 8 + 4 + 16 + 16 + 16 + 4 + (4 * 32)
	InodeEntryCountPerPage = 85
)

type InodeEntry struct {
	FsegId           uint64
	FsegNotFullNUsed uint32
	FsegFree         *FlstBaseNode
	FsegNotFull      *FlstBaseNode
	FsegFull         *FlstBaseNode
	FsegMagicN       uint32
	FsegFragArr      []uint32
}
