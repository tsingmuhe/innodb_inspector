package innodb

import "innodb_inspector/pkg/innodb/page"

type BlobPage struct {
	*BasePage
}

func (t *BlobPage) BlobHeader() *page.BlobHeader {
	c := t.PageCursorAtBodyStart()
	return &page.BlobHeader{
		Length:   c.Uint32(),
		NextPage: c.Uint32(),
	}
}
