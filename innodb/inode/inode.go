package inode

import (
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
)

type InodeEntry struct {
	FsegId           uint64               //8 该Inode归属的Segment ID，若值为0表示该slot未被使用
	FsegNotFullNUsed uint32               //4 FSEG_NOT_FULL链表上被使用的Page数量
	FsegFree         *innodb.FlstBaseNode //16 完全没有被使用并分配给该Segment的Extent链表
	FsegNotFull      *innodb.FlstBaseNode //16 至少有一个page分配给当前Segment的Extent链表，全部用完时，转移到FSEG_FULL上，全部释放时，则归还给当前表空间FSP_FREE链表
	FsegFull         *innodb.FlstBaseNode //16 分配给当前segment且Page完全使用完的Extent链表
	FsegMagicN       uint32               //4 Magic Number
	FsegFragArr      []uint32             //4*32
}

func NewInodeEntry(data []byte) []*InodeEntry {
	br := cursor.New(data)
	br.Seek(50)

	var result []*InodeEntry
	for i := 0; i < 85; i++ {
		inodeEntry := &InodeEntry{
			FsegId:           br.Uint64(),
			FsegNotFullNUsed: br.Uint32(),
			FsegFree:         innodb.NewFlstBaseNode(br.Bytes(16)),
			FsegNotFull:      innodb.NewFlstBaseNode(br.Bytes(16)),
			FsegFull:         innodb.NewFlstBaseNode(br.Bytes(16)),
			FsegMagicN:       br.Uint32(),
			FsegFragArr:      NewFsegFragArr(br.Bytes(128)),
		}

		if inodeEntry.FsegId > 0 {
			result = append(result, inodeEntry)
		}
	}

	return result
}

func NewFsegFragArr(data []byte) []uint32 {
	br := cursor.New(data)

	var result []uint32
	for i := 0; i < 32; i++ {
		pageNo := br.Uint32()
		if !innodb.IsUndefinedPageNo(pageNo) {
			result = append(result, pageNo)
		}
	}

	return result
}

type InodePage struct {
	FILHeader         *innodb.FILHeader
	FsegInodePageNode *innodb.FlstNode
	InodeEntries      []*InodeEntry
	FILTrailer        *innodb.FILTrailer
}

func NewInodePage(data []byte) *InodePage {
	return &InodePage{
		FILHeader:         innodb.NewFILHeader(data),
		FsegInodePageNode: innodb.NewFlstNode(cursor.New(data).Seek(38).Bytes(12)),
		InodeEntries:      NewInodeEntry(data),
		FILTrailer:        innodb.NewFILTrailer(data),
	}
}
