package innodb

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"innodb_inspector/pkg/innodb/page"
	"innodb_inspector/pkg/innodb/tablespace"
	"os"
)

func resolveFspHeaderSpaceId(f *os.File) uint32 {
	b := make([]byte, 4)
	f.ReadAt(b, 34)
	return binary.BigEndian.Uint32(b)
}

func parsePage(fspHeaderSpaceId, pageNo uint32, pageBits []byte) Page {
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
		case tablespace.DataDictionaryHeaderPageNo:
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
	default:
		return basePage
	}
}

type PageDesc struct {
	PageNo    uint32
	PageType  page.Type
	SpaceId   uint32
	PageNotes string
}

func OverView(name string, pageSize int) ([]*PageDesc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fspHeaderSpaceId := resolveFspHeaderSpaceId(f)
	return scanPage(f, fspHeaderSpaceId, pageSize)
}

func scanPage(f *os.File, fspHeaderSpaceId uint32, pageSize int) ([]*PageDesc, error) {
	buf := make([]byte, pageSize)
	s := bufio.NewScanner(f)
	s.Buffer(buf, pageSize)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			return 0, nil, nil
		}

		if len(data) >= pageSize {
			return pageSize, data[0:pageSize], nil
		}

		return 0, nil, nil
	})

	var pds []*PageDesc

	var i uint32
	for i = 0; s.Scan(); i++ {
		b := s.Bytes()
		if len(b) != pageSize {
			return nil, errors.New("invalid page bytes")
		}

		pg := parsePage(fspHeaderSpaceId, i, b)
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
		case tablespace.InsertBufferHeaderPageNo:
			return "insert buffer header"
		case tablespace.InsertBufferRootPageNo:
			return "insert buffer root page"
		case tablespace.TransactionSystemHeaderPageNo:
			return "transaction system header"
		case tablespace.FirstRollbackSegmentPageNo:
			return "first rollback segment"
		case tablespace.DataDictionaryHeaderPageNo:
			return "data dictionary header"
		}

		if pageNo >= tablespace.DoubleWriteBufferPageNo1 && pageNo < tablespace.DoubleWriteBufferPageNo2 {
			return "double write buffer block"
		}
	} else {
		switch pageNo {
		case tablespace.RootPageOfFirstIndexPageNo:
			return "root page of first index"
		case tablespace.RootPageOfSecondIndexPageNo:
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

func PageDetail(name string, targetPageNo, pageSize uint32, exportPath string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fspHeaderSpaceId := resolveFspHeaderSpaceId(f)

	b := make([]byte, pageSize)
	off := int64(targetPageNo) * int64(pageSize)
	_, err = f.ReadAt(b, off)
	if err != nil {
		return "", err
	}
	
	pg := parsePage(fspHeaderSpaceId, targetPageNo, b)
	if exportPath != "" {
		file, _ := os.Create(exportPath)
		file.Write(b)

		tags, _ := json.Marshal(pg.HexEditorTags())
		tagFile, _ := os.Create(exportPath + ".tags")
		tagFile.WriteString(string(tags))
	}

	return pg.String(), nil
}
