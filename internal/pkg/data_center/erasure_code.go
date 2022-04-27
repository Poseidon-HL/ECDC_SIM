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

	// LRC相关参数设置
	L                    int
	LRCDataChunkOffset   [][]int
	LRCLocalChunkParity  []int
	LRCGlobalChunkParity []int
}

func (ecf *ErasureCodeConf) CheckParams() bool {
	if ecf.K < 1 || ecf.N <= ecf.K {
		return false
	}
	if ecf.CodeType == LRC && ecf.L == 0 {
		return false
	}
	return true
}
