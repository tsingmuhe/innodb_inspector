package page

type DictionaryHeader struct {
	DictHdrRowId      uint64 //8
	DictHdrTableId    uint64 //8
	DictHdrIndexId    uint64 //8
	DictHdrMaxSpaceId uint32 //4
	DictHdrMixIdLow   uint32 //4

	DictHdrTables   uint32 //4
	DictHdrTableIds uint32 //4
	DictHdrColumns  uint32 //4
	DictHdrIndexes  uint32 //4
	DictHdrFields   uint32 //4

	//skip 4
	FsegEntry *FsegEntry
}
