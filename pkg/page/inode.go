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

type FsegEntry struct {
	FsegHdrSpace  uint32 //4
	FsegHdrPageNo uint32 //4
	FsegHdrOffset uint16 //2
}

type InodePage struct {
	*BasePage
}

func (t *InodePage) InodeEntry() []*InodeEntry {
	c := t.CursorAt(50)

	var inodeEntries []*InodeEntry
	for i := 0; i < 85; i++ {
		inodeEntry := &InodeEntry{
			FsegId:           c.Uint64(),
			FsegNotFullNUsed: c.Uint32(),
			FsegFree:         c.FlstBaseNode(),
			FsegNotFull:      c.FlstBaseNode(),
			FsegFull:         c.FlstBaseNode(),
			FsegMagicN:       c.Uint32(),
		}

		for i := 0; i < 32; i++ {
			pageNo := c.Uint32()
			if !IsUndefinedPageNo(pageNo) {
				inodeEntry.FsegFragArr = append(inodeEntry.FsegFragArr, pageNo)
			}
		}

		if inodeEntry.FsegId > 0 {
			inodeEntries = append(inodeEntries, inodeEntry)
		}
	}

	return inodeEntries
}
