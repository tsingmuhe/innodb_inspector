package innodb

import (
	"encoding/binary"
	"errors"
	"innodb_inspector/pkg/innodb/page"
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
	case page.FilPageTypeBlob:
		return &BlobPage{
			BasePage: basePage,
		}
	default:
		return basePage
	}
}

func PageDetail(filePath string, targetPageNo, pageSize uint32) (string, error) {
	dbFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer dbFile.Close()

	fspHeaderSpaceId, err := resolveFspHeaderSpaceId(dbFile)
	if err != nil {
		return "", err
	}

	buf := make([]byte, pageSize)
	n, _ := dbFile.ReadAt(buf, int64(targetPageNo*pageSize))
	if n <= 0 {
		return "", errors.New("bad file")
	}

	pg := ParsePage(fspHeaderSpaceId, targetPageNo, buf)
	return pg.String(), nil
}

type PageDesc struct {
	PageNo    uint32
	PageType  page.Type
	SpaceId   uint32
	PageNotes string
}

func OverView(filePath string, pageSize uint32) ([]*PageDesc, error) {
	dbFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()

	fspHeaderSpaceId, err := resolveFspHeaderSpaceId(dbFile)
	if err != nil {
		return nil, err
	}

	dbFileStat, _ := dbFile.Stat()
	dbFileSize := dbFileStat.Size()
	pageCount := uint32(dbFileSize) / pageSize

	buf := make([]byte, pageSize)
	var pds []*PageDesc

	for pageNo := uint32(0); pageNo < pageCount; pageNo++ {
		n, _ := dbFile.ReadAt(buf, int64(pageNo*pageSize))
		if n <= 0 {
			return nil, errors.New("bad file")
		}

		pg := ParsePage(fspHeaderSpaceId, pageNo, buf)
		pds = append(pds, &PageDesc{
			PageNo:    pg.PageNo(),
			PageType:  pg.Type(),
			SpaceId:   pg.SpaceId(),
			PageNotes: pageNotes(pg),
		})
	}

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

		return ""
	}

	switch pageNo {
	case page.RootPageOfFirstIndexPageNo:
		return "root page of first index"
	case page.RootPageOfSecondIndexPageNo:
		return "root page of second index"
	}

	return ""
}

func resolveFspHeaderSpaceId(dbFile *os.File) (uint32, error) {
	buf := make([]byte, 4)
	_, err := dbFile.ReadAt(buf, 34)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}
