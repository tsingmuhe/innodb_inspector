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
	FilePath string
	PageSize uint
}

func (t *Tablespace) OverView() ([]page.Page, error) {
	var pages []page.Page

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		pg := page.NewBasePage(pageNo, data)
		pages = append(pages, pg)
		return false, nil
	})

	return pages, nil
}

func (t *Tablespace) PageDetail(targetPageNo uint32) (string, error) {
	var pageDetail string

	t.ForeachPage(func(pageNo uint32, data []byte) (bool, error) {
		if pageNo == targetPageNo {
			pg := page.NewBasePage(pageNo, data)
			pageDetail = pg.String()
			return true, nil
		}

		return false, nil
	})

	return pageDetail, nil
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
		PageSize: page.DefaultSize,
	}
}
