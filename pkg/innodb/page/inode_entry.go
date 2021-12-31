package page

type InodeEntry struct {
	FsegId           uint64
	FsegNotFullNUsed uint32
	FsegFree         *FlstBaseNode
	FsegNotFull      *FlstBaseNode
	FsegFull         *FlstBaseNode
	FsegMagicN       uint32
	FsegFragArr      []uint32
}
