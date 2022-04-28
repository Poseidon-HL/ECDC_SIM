package data_center

import "ECDC_SIM/internal/pkg/util"

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

type LRCChunkPlace int8

const (
	DataChunk LRCChunkPlace = iota
	LocalChunkParity
	GlobalChunkParity
	NotDefined
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

func (ecf *ErasureCodeConf) LRCChunkType(offset int) (LRCChunkPlace, int) {
	if util.ContainsInt(ecf.LRCLocalChunkParity, offset) {
		return LocalChunkParity, 0
	}
	if util.ContainsInt(ecf.LRCGlobalChunkParity, offset) {
		return GlobalChunkParity, 0
	}
	for gid, dataChunks := range ecf.LRCDataChunkOffset {
		if util.ContainsInt(dataChunks, offset) {
			return DataChunk, gid
		}
	}
	return NotDefined, 0
}
