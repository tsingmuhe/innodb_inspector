package innodb

import (
	"encoding/binary"
	"innodb_inspector/pkg/innodb/page"
	"io"
	"os"
)

func ParsePage(fspHeaderSpaceId, pageNo uint32, pageBits []byte) Page {
	basePage := NewBasePage(fspHeaderSpaceId, pageNo, pageBits)

	pageType := basePage.Type()
	switch pageType {
	case page.FilPageTypeFspHdr:
		return &FspHdrPage{
			BasePage: basePage,
		}
	case page.FilPageTypeXdes:
		return &XdesPage{
			BasePage: basePage,
		}
	case page.FilPageInode:
		return &InodePage{
			BasePage: basePage,
		}
	case page.FilPageTypeSys:
		switch pageNo {
		case page.InsertBufferHeaderPageNo:
			return &IBufHeaderPage{
				BasePage: basePage,
			}
		case page.FirstRollbackSegmentPageNo:
			return &SysRsegHeaderPage{
				BasePage: basePage,
			}
		case page.DataDictionaryHeaderPageNo:
			return &DictionaryHeaderPage{
				BasePage: basePage,
			}
		default:
			return basePage
		}
	case page.FilPageIndex:
		return &IndexPage{
			BasePage: basePage,
		}
	case page.FilPageTypeBlob:
		return &BlobPage{
			BasePage: basePage,
		}
	default:
		return basePage
	}
}

func PageDetail(filePath string, targetPageNo, pageSize uint32) (string, error) {
	fspHeaderSpaceId := resolveFspHeaderSpaceId(filePath, pageSize)

	var result string

	forEachPage(filePath, pageSize, func(pageNo uint32, bytes []byte) (bool, error) {
		if pageNo == targetPageNo {
			pg := ParsePage(fspHeaderSpaceId, targetPageNo, bytes)
			result = pg.String()
			return true, nil
		}
		return false, nil
	})

	return result, nil
}

type PageDesc struct {
	PageNo    uint32
	PageType  page.Type
	SpaceId   uint32
	PageNotes string
}

func OverView(filePath string, pageSize uint32) ([]*PageDesc, error) {
	fspHeaderSpaceId := resolveFspHeaderSpaceId(filePath, pageSize)

	var pds []*PageDesc

	forEachPage(filePath, pageSize, func(pageNo uint32, bytes []byte) (bool, error) {
		pg := ParsePage(fspHeaderSpaceId, pageNo, bytes)
		pds = append(pds, &PageDesc{
			PageNo:    pg.PageNo(),
			PageType:  pg.Type(),
			SpaceId:   pg.SpaceId(),
			PageNotes: pageNotes(pg),
		})
		return false, nil
	})

	return pds, nil
}

func pageNotes(pg Page) string {
	pageNo := pg.PageNo()

	if pg.IsSysTablespace() {
		switch pageNo {
		case 0:
			return "system tablespace"
		case page.InsertBufferHeaderPageNo:
			return "insert buffer header"
		case page.InsertBufferRootPageNo:
			return "insert buffer root page"
		case page.TransactionSystemHeaderPageNo:
			return "transaction system header"
		case page.FirstRollbackSegmentPageNo:
			return "first rollback segment"
		case page.DataDictionaryHeaderPageNo:
			return "data dictionary header"
		}

		if pageNo >= page.DoubleWriteBufferPageNo1 && pageNo < page.DoubleWriteBufferPageNo2 {
			return "double write buffer block"
		}
	} else {
		switch pageNo {
		case page.RootPageOfFirstIndexPageNo:
			return "root page of first index"
		case page.RootPageOfSecondIndexPageNo:
			return "root page of second index"
		}
	}

	if val, ok := pg.(*IndexPage); ok {
		if val.IsCompact() {
			return "compact format"
		}
		return "redundant format"
	}

	return ""
}

func resolveFspHeaderSpaceId(filePath string, pageSize uint32) uint32 {
	var fspHeaderSpaceId uint32

	forEachPage(filePath, pageSize, func(pageNo uint32, bytes []byte) (bool, error) {
		fspHeaderSpaceId = binary.BigEndian.Uint32(bytes[34:38])
		pg := NewBasePage(0, pageNo, bytes)
		fspHeaderSpaceId = pg.SpaceId()
		return true, nil
	})

	return fspHeaderSpaceId
}

func forEachPage(filePath string, pageSize uint32, handle func(uint32, []byte) (bool, error)) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, pageSize)
	var pageNo uint32

	for {
		switch n, err := f.Read(buf); true {
		case n > 0:
			breakLoop, err := handle(pageNo, buf)
			if err != nil {
				return err
			}

			if breakLoop {
				return nil
			}
		case n == 0 && err == io.EOF: // EOF
			return nil
		case err != nil:
			return err
		}

		pageNo++
	}
}
