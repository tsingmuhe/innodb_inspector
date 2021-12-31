package innodb

import (
	"encoding/binary"
	"errors"
	page2 "innodb_inspector/pkg/innodb/page"
	"os"
)

func ParsePage(fspHeaderSpaceId, pageNo uint32, pageBits []byte) page2.Page {
	basePage := page2.NewBasePage(fspHeaderSpaceId, pageNo, pageBits)

	pageType := basePage.Type()
	switch pageType {
	case page2.FilPageTypeFspHdr:
		return &page2.FspHdrXdesPage{
			BasePage: basePage,
		}
	case page2.FilPageInode:
		return &page2.InodePage{
			BasePage: basePage,
		}
	case page2.FilPageTypeSys:
		switch pageNo {
		case page2.InsertBufferHeaderPageNo:
			return &page2.IBufHeaderPage{
				BasePage: basePage,
			}
		case page2.FirstRollbackSegmentPageNo:
			return &page2.SysRsegHeaderPage{
				BasePage: basePage,
			}
		case page2.DataDictionaryHeaderPageNo:
			return &page2.DictionaryHeaderPage{
				BasePage: basePage,
			}
		default:
			return basePage
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
	PageType  page2.Type
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

func pageNotes(pg page2.Page) string {
	pageNo := pg.PageNo()

	if pg.IsSysTablespace() {
		switch pageNo {
		case 0:
			return "system tablespace"
		case page2.InsertBufferHeaderPageNo:
			return "insert buffer header"
		case page2.InsertBufferRootPageNo:
			return "insert buffer root page"
		case page2.TransactionSystemHeaderPageNo:
			return "transaction system header"
		case page2.FirstRollbackSegmentPageNo:
			return "first rollback segment"
		case page2.DataDictionaryHeaderPageNo:
			return "data dictionary header"
		}

		if pageNo >= page2.DoubleWriteBufferPageNo1 && pageNo < page2.DoubleWriteBufferPageNo2 {
			return "double write buffer block"
		}

		return ""
	}

	switch pageNo {
	case page2.RootPageOfFirstIndexPageNo:
		return "root page of first index"
	case page2.RootPageOfSecondIndexPageNo:
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
