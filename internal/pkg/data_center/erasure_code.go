package data_center

type ErasureCodeType int8

const (
	RS ErasureCodeType = iota
	LRC
)

type ChunkPlaceType int8

const (
	FLAT ChunkPlaceType = iota
	HIERARCHICAL
)

type ErasureCodeConf struct {
	CodeType       ErasureCodeType
	ChunkPlaceType ChunkPlaceType
	N              int
	K              int
}
