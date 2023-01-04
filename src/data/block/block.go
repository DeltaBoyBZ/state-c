package block

type FieldType int8

const (
    FIELD_TYPE_INT FieldType = iota
    FIELD_TYPE_FLOAT
)

type BlockID int 
type BlockCategory int8
type FieldID int
type FloatFieldID int

const (
    BLOCK_INT BlockCategory = iota
    BLOCK_FLOAT
    BLOCK_ARRAY_INT
    BLOCK_ARRAY_FLOAT
    BLOCK_FIELD_INT
    BLOCK_FIELD_FLOAT
    BLOCK_DUMP
    BLOCK_ENUM_VAL
)

var MinFreeBlockID BlockID = 0
var MinFreeFieldID FieldID = 0

var Categories map[BlockID]BlockCategory = map[BlockID]BlockCategory{}
var Labels map[BlockID]string = map[BlockID]string{}
var IntVals map[BlockID]int = map[BlockID]int{}
var FloatVals map[BlockID]float32 = map[BlockID]float32{}
var IntArrays map[BlockID][]int = map[BlockID][]int{}
var FloatArrays map[BlockID][]float32 = map[BlockID][]float32{}
var DumpPaths map[BlockID]string = map[BlockID]string{}
var FileLengths map[BlockID]int = map[BlockID]int{}

var FieldTypes map[BlockID]FieldType = map[BlockID]FieldType{}

func CreateBlock () BlockID {
    MinFreeBlockID++
    return MinFreeBlockID - 1
}

func CreateField () FieldID {
    MinFreeFieldID++
    return MinFreeFieldID - 1 
}

