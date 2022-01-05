package tablespace

const (
	InsertBufferHeaderPageNo      = 3
	InsertBufferRootPageNo        = 4
	TransactionSystemHeaderPageNo = 5
	FirstRollbackSegmentPageNo    = 6
	DataDictionaryHeaderPageNo    = 7

	DoubleWriteBufferPageNo1 = 64
	DoubleWriteBufferPageNo2 = 192
)
