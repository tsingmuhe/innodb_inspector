package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type InodePage struct {
	*BasePage
}

func (t *InodePage) InodeList() *page.FlstNode {
	inodeList := page.NewFlstNode(page.FilHeaderSize, t.pageBytes)
	if IsUndefinedPageNo(inodeList.Pre.PageNo) {
		return nil
	}
	return inodeList
}

func (t *InodePage) InodeEntry() []*page.InodeEntry {
	return page.NewInodeEntries(t.pageBytes)
}

func (t *InodePage) HexEditorTags() []*page.HexEditorTag {
	var result []*page.HexEditorTag
	result = append(result, t.FilHeader().HexEditorTag()...)

	inodeList := t.InodeEntry()
	for _, inode := range inodeList {
		result = append(result, inode.HexEditorTag())
	}

	result = append(result, t.FILTrailer().HexEditorTag())
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
