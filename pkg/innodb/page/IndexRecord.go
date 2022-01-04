package page

type IndexRecord struct {
	OffSet uint32

	RecordType       uint
	NextRecordOffset int16
}

func (t *IndexRecord) NextRecord() uint32 {
	return uint32(int32(t.OffSet) + int32(t.NextRecordOffset))
}
