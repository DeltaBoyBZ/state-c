package build

import (
    "fmt"
    "delta/statesy/data/block"
)

func BuildAssembly(blocks []block.BlockID) string {
    source := "section .data\n"
    for i := 0 ; i < len(blocks) ; i++ {
        blockID := blocks[i]
        switch block.Categories[blockID] {
            case block.BLOCK_INT:
                source += fmt.Sprintf("global %s\n%s: dd %d\n", block.Labels[blockID], block.Labels[blockID], block.IntVals[blockID])
                break
            case block.BLOCK_FLOAT:
                source += fmt.Sprintf("global %s\n%s: dd %f\n", block.Labels[blockID], block.Labels[blockID], block.FloatVals[blockID]) 
                break
            case block.BLOCK_ARRAY_INT:
                source += fmt.Sprintf("global %s\n%s_ELEM: dd ", block.Labels[blockID], block.Labels[blockID])
                numElem := len(block.IntArrays[blockID])
                for index, elem := range block.IntArrays[blockID] {
                    source += fmt.Sprintf("%d", elem)
                    if index != numElem - 1 { source += ", " }
                }
                source += "\n"
                source += fmt.Sprintf("%s: dq %s_ELEM\n", block.Labels[blockID], block.Labels[blockID]);
                break 
            case block.BLOCK_ARRAY_FLOAT:
                source += fmt.Sprintf("global %s\n%s_ELEM: dd ", block.Labels[blockID], block.Labels[blockID])
                numElem := len(block.FloatArrays[blockID])
                for index, elem := range block.FloatArrays[blockID] {
                    source += fmt.Sprintf("%f", elem)
                    if index != numElem - 1 { source += ", " }
                }
                source += "\n"
                source += fmt.Sprintf("%s: dq %s_ELEM\n", block.Labels[blockID], block.Labels[blockID]);
                break 
            case block.BLOCK_FIELD_INT:
                source += fmt.Sprintf("global %s\n%s_ELEM: dd ", block.Labels[blockID], block.Labels[blockID])
                numElem := len(block.IntArrays[blockID])
                for index, elem := range block.IntArrays[blockID] {
                    source += fmt.Sprintf("%d", elem)
                    if index != numElem - 1 { source += ", " }
                }
                source += "\n"
                source += fmt.Sprintf("%s: dq %s_ELEM\n", block.Labels[blockID], block.Labels[blockID])
                break
            case block.BLOCK_FIELD_FLOAT:
                source += fmt.Sprintf("global %s\n%s_ELEM: dd ", block.Labels[blockID], block.Labels[blockID])
                numElem := len(block.FloatArrays[blockID])
                for index, elem := range block.FloatArrays[blockID] {
                    source += fmt.Sprintf("%f", elem)
                    if index != numElem - 1 { source += ", " }
                }
                source += "\n"
                source += fmt.Sprintf("%s: dq %s_ELEM\n", block.Labels[blockID], block.Labels[blockID])
                break
            case block.BLOCK_DUMP:
                source += fmt.Sprintf("global %s\nglobal %s_LEN\n%s_DATA: incbin \"%s\"\n%s: dq %s_DATA\n%s_LEN: dd %d\n", block.Labels[blockID], block.Labels[blockID], block.Labels[blockID], block.DumpPaths[blockID], block.Labels[blockID], block.Labels[blockID], block.Labels[blockID], block.FileLengths[blockID]) 
                break
            case block.BLOCK_ENUM_VAL:
                source += fmt.Sprintf("global %s\n%s: dd %d\n", block.Labels[blockID], block.Labels[blockID], block.IntVals[blockID])
                break
            default:
        }
    }
    return source
}

func BuildHeader(blocks []block.BlockID) string {
    source := "#pragma once\n\n"
    for i := 0 ; i < len(blocks) ; i++ {
        blockID := blocks[i]
        switch block.Categories[blockID] {
            case block.BLOCK_INT:
                source += fmt.Sprintf("extern int %s;\n", block.Labels[blockID])
                break
            case block.BLOCK_FLOAT:
                source += fmt.Sprintf("extern float %s;\n", block.Labels[blockID])
                break
            case block.BLOCK_ARRAY_INT:
                source += fmt.Sprintf("extern int* %s;\n", block.Labels[blockID])
                break 
            case block.BLOCK_ARRAY_FLOAT:
                source += fmt.Sprintf("extern float* %s;\n", block.Labels[blockID])
                break
            case block.BLOCK_FIELD_INT:
                source += fmt.Sprintf("extern int* %s;\n", block.Labels[blockID])
                break
            case block.BLOCK_FIELD_FLOAT:
                source += fmt.Sprintf("extern float* %s;\n", block.Labels[blockID])
                break
            case block.BLOCK_DUMP:
                source += fmt.Sprintf("extern void* %s;\nextern int %s_LEN;\n", block.Labels[blockID], block.Labels[blockID])
                break
            case block.BLOCK_ENUM_VAL:
                source += fmt.Sprintf("extern int %s;\n", block.Labels[blockID])
                break
            default:
        }
    }
    return source
}

