package innodb

import (
	"innodb_inspector/pkg/page"
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

type Tablespace struct {
	filePath string
	pageSize uint
}

func (t *Tablespace) FSPHeaderSpaceId() uint32 {
	var fspHeaderSpaceId uint32

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		c := page.NewCursor(data, 34)
		fspHeaderSpaceId = c.Uint32()
		return true, nil
	})

	return fspHeaderSpaceId
}

func (t *Tablespace) OverView() ([]*PageDesc, error) {
	var pds []*PageDesc
	fspHeaderSpaceId := t.FSPHeaderSpaceId()

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		pg := page.Parse(fspHeaderSpaceId, pageNo, data)
		pds = append(pds, &PageDesc{
			PageNo:    pg.PageNo(),
			PageType:  pg.Type(),
			SpaceId:   pg.SpaceId(),
			PageNotes: t.PageNotes(pg),
		})
		return false, nil
	})

	return pds, nil
}

func (t *Tablespace) PageNotes(pg page.Page) string {
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

func (t *Tablespace) PageDetail(targetPageNo uint32) (string, error) {
	var pageDetail string
	fspHeaderSpaceId := t.FSPHeaderSpaceId()

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		if pageNo == targetPageNo {
			pg := page.Parse(fspHeaderSpaceId, pageNo, data)
			pageDetail = pg.String()
			return true, nil
		}

		return false, nil
	})

	return pageDetail, nil
}

func (t *Tablespace) ForeachPage(handle func(uint32, []byte) (bool, error)) error {
	f, err := os.Open(t.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, t.pageSize)
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

type PageDesc struct {
	PageNo    uint32
	PageType  page.Type
	SpaceId   uint32
	PageNotes string
}

func NewTablespace(filePath string) *Tablespace {
	return &Tablespace{
		filePath: filePath,
		pageSize: page.DefaultSize,
	}
}
