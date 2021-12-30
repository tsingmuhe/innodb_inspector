package tablespace

import (
	"encoding/json"
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
	"innodb_inspector/innodb/inode"
	"innodb_inspector/innodb/sys"
	"innodb_inspector/innodb/xdes"
	"io"
	"os"
)

const (
	initPageNo uint32 = 0

	insertBufferHeaderPageNo      = 3
	insertBufferRootPageNo        = 4
	transactionSystemHeaderPageNo = 5
	firstRollbackSegmentPageNo    = 6
	dataDictionaryHeaderPageNo    = 7

	doubleWriteBufferPageNo1 = 64
	doubleWriteBufferPageNo2 = 192

	rootPageOfFirstIndexPageNo  = 3
	rootPageOfSecondIndexPageNo = 4
)

type PageDesc struct {
	PageNo        uint32
	FilPageOffset uint32
	FilPageType   innodb.PageType
	SpaceId       uint32
	Comments      string
}

type Tablespace struct {
	FilePath string
	PageSize uint

	spaceId uint32
}

func (t *Tablespace) isSysTablespace() bool {
	return t.spaceId == 0
}

func (t *Tablespace) OverView() ([]*PageDesc, error) {
	var result []*PageDesc

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		br := cursor.New(data)

		pageDesc := &PageDesc{
			PageNo:        pageNo,
			FilPageOffset: br.Seek(4).Uint32(),
			FilPageType:   innodb.PageType(br.Seek(24).Uint16()),
			SpaceId:       br.Seek(34).Uint32(),
			Comments:      "",
		}

		if pageNo == initPageNo {
			t.spaceId = pageDesc.SpaceId
		}

		if t.isSysTablespace() {
			if pageNo == initPageNo {
				pageDesc.Comments = "system tablespace"
			}

			if pageNo == insertBufferHeaderPageNo {
				pageDesc.Comments = "insert buffer header"
			}

			if pageNo == insertBufferRootPageNo {
				pageDesc.Comments = "insert buffer root"
			}

			if pageNo == transactionSystemHeaderPageNo {
				pageDesc.Comments = "transaction system header"
			}

			if pageNo == firstRollbackSegmentPageNo {
				pageDesc.Comments = "first rollback segment"
			}

			if pageNo == dataDictionaryHeaderPageNo {
				pageDesc.Comments = "data dictionary header"
			}

			if pageNo >= doubleWriteBufferPageNo1 && pageNo < doubleWriteBufferPageNo2 {
				pageDesc.Comments = "double write buffer block"
			}
		} else {
			if pageNo == rootPageOfFirstIndexPageNo && pageDesc.FilPageType == innodb.PageTypeIndex {
				pageDesc.Comments = "root page of first index"
			}

			if pageNo == rootPageOfSecondIndexPageNo && pageDesc.FilPageType == innodb.PageTypeIndex {
				pageDesc.Comments = "root page of second index"
			}
		}

		result = append(result, pageDesc)
		return false, nil
	})

	return result, nil
}

func (t *Tablespace) PageDetail(targetPageNo uint32) ([]byte, error) {
	var result []byte

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		if pageNo == initPageNo {
			t.spaceId = innodb.ResolveSpaceId(data)
		}

		if pageNo == targetPageNo {
			result, _ = json.MarshalIndent(t.doPageDetail(pageNo, data), "", "  ")
			return true, nil
		}

		return false, nil
	})

	return result, nil
}

func (t *Tablespace) doPageDetail(pageNo uint32, data []byte) interface{} {
	pageType := innodb.ResolvePageType(data)

	if pageType == innodb.PageTypeFspHdr {
		return xdes.NewFSPHDRPage(data)
	}

	if pageType == innodb.PageTypeXdes {
		return xdes.NewXDESPage(data)
	}

	if pageType == innodb.PageTypeInode {
		return inode.NewInodePage(data)
	}

	if pageType == innodb.PageTypeSys && t.isSysTablespace() && pageNo == insertBufferHeaderPageNo {
		return sys.NewIBufHeaderPage(data)
	}

	if pageType == innodb.PageTypeSys && t.isSysTablespace() && pageNo == firstRollbackSegmentPageNo {
		return sys.NewSysRsegHeaderPage(data)
	}

	if pageType == innodb.PageTypeSys && t.isSysTablespace() && pageNo == dataDictionaryHeaderPageNo {
		return sys.NewDictionaryHeaderPage(data)
	}

	return &innodb.PageOverview{
		FILHeader:  innodb.NewFILHeader(data),
		FILTrailer: innodb.NewFILTrailer(data),
	}
}

func (t *Tablespace) ForeachPage(handle func(uint32, []byte) (bool, error)) error {
	f, err := os.Open(t.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, t.PageSize)
	pageNo := initPageNo

	for {
		switch n, err := f.Read(buf); true {
		case n > 0:
			breakLoop, err := handle(pageNo, buf)
			if err != nil {
				return err
			}

			if breakLoop {
				break
			}
		case n == 0 && err == io.EOF: // EOF
			return nil
		case err != nil:
			return err
		}

		pageNo++
	}
}

func NewTablespace(filePath string) *Tablespace {
	return &Tablespace{
		FilePath: filePath,
		PageSize: innodb.DefaultInnodbPageSize,
	}
}
