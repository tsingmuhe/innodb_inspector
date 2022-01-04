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
		inodeEntry := &page.InodeEntry{
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
