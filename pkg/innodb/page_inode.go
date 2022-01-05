package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type InodePage struct {
	*BasePage
}

func (t *InodePage) InodeList() *page.FlstNode {
	c := t.PageCursorAtBodyStart()
	inodeList := c.FlstNode()
	if IsUndefinedPageNo(inodeList.Pre.PageNo) {
		return nil
	}
	return inodeList
}

func (t *InodePage) InodeEntry() []*page.InodeEntry {
	c := t.PageCursorAt(page.FilHeaderSize + page.FlstNodeSize)

	var inodeEntries []*page.InodeEntry
	for i := 0; i < 85; i++ {
		inodeEntry := page.NewInodeEntry(uint32(i))
		inodeEntry.FsegId = c.Uint64()
		inodeEntry.FsegNotFullNUsed = c.Uint32()
		inodeEntry.FsegFree = c.FlstBaseNode()
		inodeEntry.FsegNotFull = c.FlstBaseNode()
		inodeEntry.FsegFull = c.FlstBaseNode()
		inodeEntry.FsegMagicN = c.Uint32()

		for j := 0; j < 32; j++ {
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

func (t *InodePage) HexEditorTags() []*page.HexEditorTag {
	var result []*page.HexEditorTag
	result = append(result, t.FilHeader().HexEditorTag())

	inodeList := t.InodeEntry()
	for _, inode := range inodeList {
		result = append(result, inode.HexEditorTag())
	}

	result = append(result, t.FILTrailer().HexEditorTag(len(t.pageBits)))
	return result
}

func (t *InodePage) String() string {
	type Page struct {
		FILHeader    *page.FILHeader
		InodeList    *page.FlstNode
		InodeEntries []*page.InodeEntry
		FILTrailer   *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:    t.FilHeader(),
		InodeList:    t.InodeList(),
		InodeEntries: t.InodeEntry(),
		FILTrailer:   t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
