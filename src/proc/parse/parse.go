package parse

import (
    "fmt"
    "os"
    "delta/statesy/data/token"
    "delta/statesy/data/block"
)

func getIntValByLabel(label string, prefix string, blocks []block.BlockID) int {
    for _, blockID := range blocks {
        if prefix + label == block.Labels[blockID] { return block.IntVals[blockID] }
    }
    return 0 
}

func getFloatValByLabel(label string, prefix string, blocks []block.BlockID) float32 {
    for _, blockID := range blocks {
        if prefix + label == block.Labels[blockID] { return block.FloatVals[blockID] }
    }
    return 0 
}

func ParseTokens(tokens []token.Token) []block.BlockID {
    tokenIndex := 0
    numTokens  := len(tokens)
    blocks := make([]block.BlockID, 0, 32)
    var prefix string = ""
    for tokenIndex < numTokens {
        if tokens[tokenIndex].Category == token.TOKEN_PREFIX {
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected an identifier.")
                tokenIndex++
                continue
            }
            prefix = tokens[tokenIndex].Source
            tokenIndex++
        } else if tokens[tokenIndex].Category == token.TOKEN_DUMP {
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected an identifier.")
                tokenIndex++
                continue
            }
            dumpLabel := tokens[tokenIndex].Source 
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_LITERAL_STRING {
                fmt.Println("Expected a string literal.")
                tokenIndex++
                continue
            }
            dumpPath := tokens[tokenIndex].StringVal
            tokenIndex++
            newBlockID := block.CreateBlock()
            block.Categories[newBlockID] = block.BLOCK_DUMP
            block.Labels[newBlockID] = prefix + dumpLabel
            block.DumpPaths[newBlockID] = dumpPath
            contents, err := os.ReadFile(dumpPath) 
            if err != nil {
                fmt.Printf("The file `%s` cannot be read.", dumpPath)
                tokenIndex++
                continue
            }
            block.FileLengths[newBlockID] = len(contents)
            blocks = append(blocks, newBlockID)
        } else if tokens[tokenIndex].Category == token.TOKEN_INT {
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER  {
                fmt.Println("Expected an identifier.")
                tokenIndex++
                continue
            }
            blockLabel := tokens[tokenIndex].Source
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_EQUALS {
                fmt.Println("Expected an `=`.")
                tokenIndex++
                continue
            }
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_LITERAL_INT && tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected an integer literal or identifier.")
                tokenIndex++
                continue
            }
            newBlockID := block.CreateBlock()
            block.Categories[newBlockID] = block.BLOCK_INT
            if tokens[tokenIndex].Category == token.TOKEN_LITERAL_INT {
                block.IntVals[newBlockID]  = tokens[tokenIndex].IntVal
            } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {

                block.IntVals[newBlockID]   = getIntValByLabel(tokens[tokenIndex].Source, prefix, blocks)
            }
            block.Labels[newBlockID]     = prefix + blockLabel
            blocks = append(blocks, newBlockID)
            tokenIndex++
        } else if tokens[tokenIndex].Category == token.TOKEN_FLOAT {
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected an identifier.")
                tokenIndex++
                continue
            }
            blockLabel := tokens[tokenIndex].Source
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_EQUALS {
                fmt.Println("Expected an `=`.")
                tokenIndex++
                continue
            }
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_LITERAL_FLOAT && tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected a float literal or identifier.") 
                tokenIndex++
                continue
            }
            newBlockID := block.CreateBlock()
            block.Categories[newBlockID] = block.BLOCK_FLOAT
            if tokens[tokenIndex].Category == token.TOKEN_LITERAL_FLOAT {
                block.FloatVals[newBlockID]    = tokens[tokenIndex].FloatVal
            } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                block.FloatVals[newBlockID]    = getFloatValByLabel(tokens[tokenIndex].Source, prefix, blocks)
            }
            block.Labels[newBlockID]     = prefix + blockLabel
            blocks = append(blocks, newBlockID)
            tokenIndex++
        } else if tokens[tokenIndex].Category == token.TOKEN_ARRAY {
            var elemType token.NumberCategory
            tokenIndex++
            switch tokens[tokenIndex].Category {
                case token.TOKEN_INT: 
                    elemType = token.NUMBER_INT 
                    break
                case  token.TOKEN_FLOAT:
                    elemType = token.NUMBER_FLOAT
                    break
                default:
                    fmt.Println("Unrecognised type.")
                    tokenIndex++
                    continue
            }
            tokenIndex++
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected identifier.")
                tokenIndex++
                continue
            }
            fieldLabel := tokens[tokenIndex].Source
            tokenIndex++
            if elemType == token.NUMBER_INT {
                elems := make([]int, 0, 32)
                for tokenIndex < len(tokens) && (tokens[tokenIndex].Category == token.TOKEN_LITERAL_INT || tokens[tokenIndex].Category  == token.TOKEN_IDENTIFIER) {
                    var x int
                    if tokens[tokenIndex].Category == token.TOKEN_LITERAL_INT {
                        x = tokens[tokenIndex].IntVal
                    } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                        x = getIntValByLabel(tokens[tokenIndex].Source, prefix, blocks)
                    }
                    elems = append(elems, x)
                    tokenIndex++
                }
                newBlockID := block.CreateBlock()
                block.Categories[newBlockID] = block.BLOCK_ARRAY_INT
                block.Labels[newBlockID] = prefix + fieldLabel
                block.IntArrays[newBlockID] = elems
                blocks = append(blocks, newBlockID)
            } else if(elemType == token.NUMBER_FLOAT) {
                elems := make([]float32, 0, 32)
                for tokenIndex < len(tokens) && (tokens[tokenIndex].Category == token.TOKEN_LITERAL_FLOAT || tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER) {
                    var x float32 
                    if tokens[tokenIndex].Category == token.TOKEN_LITERAL_FLOAT {
                        x = tokens[tokenIndex].FloatVal
                    } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                        x = getFloatValByLabel(tokens[tokenIndex].Source, prefix, blocks)
                    }
                    elems = append(elems, x)
                    tokenIndex++
                }
                newBlockID := block.CreateBlock()
                block.Categories[newBlockID] = block.BLOCK_ARRAY_FLOAT
                block.Labels[newBlockID] = prefix + fieldLabel
                block.FloatArrays[newBlockID] = elems
                blocks = append(blocks, newBlockID)
            }
        } else if tokens[tokenIndex].Category == token.TOKEN_TABLE {
            tokenIndex++ 
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Printf("Expected identifier: %s\n.", tokens[tokenIndex].Source)
                continue
            }
            tableName := tokens[tokenIndex].Source
            fieldIDs := make([]block.BlockID, 0, 8)
            tokenIndex++
            specArgs := true
            for specArgs {
                switch tokens[tokenIndex].Category {
                    case token.TOKEN_INT:
                        tokenIndex++
                        if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                            fmt.Println("Expected an identifier.")
                            tokenIndex++
                            continue
                        }
                        newBlockID := block.CreateBlock() 
                        block.Labels[newBlockID] = prefix + tableName + "_" + tokens[tokenIndex].Source
                        block.FieldTypes[newBlockID] = block.FIELD_TYPE_INT
                        block.IntArrays[newBlockID] = make([]int, 0, 32)
                        block.Categories[newBlockID] = block.BLOCK_FIELD_INT
                        fieldIDs = append(fieldIDs, newBlockID)
                        blocks = append(blocks, newBlockID)
                        tokenIndex++ 
                        break
                    case token.TOKEN_FLOAT:
                        tokenIndex++
                        if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                            fmt.Println("Expected an identifier.")
                            tokenIndex++
                            continue
                        }
                        newBlockID := block.CreateBlock() 
                        block.Labels[newBlockID] = prefix + tableName + "_" + tokens[tokenIndex].Source
                        block.FieldTypes[newBlockID] = block.FIELD_TYPE_FLOAT
                        block.FloatArrays[newBlockID] = make([]float32, 0, 32)
                        block.Categories[newBlockID] = block.BLOCK_FIELD_FLOAT
                        fieldIDs = append(fieldIDs, newBlockID)
                        blocks = append(blocks, newBlockID)
                        tokenIndex++ 
                        break
                    default:
                        specArgs = false
                }
            }
            if len(fieldIDs) == 0 { continue }
            tokenIndex0 := tokenIndex
            for tokenIndex < len(tokens) && (tokens[tokenIndex].Category == token.TOKEN_LITERAL_INT || tokens[tokenIndex].Category == token.TOKEN_LITERAL_FLOAT || tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER) {
                fieldIndex := (tokenIndex - tokenIndex0) % len(fieldIDs)
                fieldID := fieldIDs[fieldIndex]
                switch block.FieldTypes[fieldID] {
                    case block.FIELD_TYPE_INT:
                        var x int
                        if tokens[tokenIndex].Category == token.TOKEN_LITERAL_INT {
                            x = tokens[tokenIndex].IntVal 
                        } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                            x = getIntValByLabel(tokens[tokenIndex].Source, prefix, blocks)
                        } else {
                            fmt.Println("Expected an integer literal or identifier.")
                            tokenIndex++
                            continue
                        }
                        block.IntArrays[fieldID] = append(block.IntArrays[fieldID], x)
                        break
                    case block.FIELD_TYPE_FLOAT:
                        var x float32 
                        if tokens[tokenIndex].Category == token.TOKEN_LITERAL_FLOAT {
                            x = tokens[tokenIndex].FloatVal
                        } else if tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                            x = getFloatValByLabel(tokens[tokenIndex].Source, prefix, blocks)
                        } else {
                            fmt.Println("Expected a float literal or identifier.")
                            tokenIndex++
                            continue
                        }
                        block.FloatArrays[fieldID] = append(block.FloatArrays[fieldID], x)
                        break
                    default:
                }
                tokenIndex++
            }

        } else if tokens[tokenIndex].Category == token.TOKEN_ENUM {
            tokenIndex++
            runningVal := 0
            if tokens[tokenIndex].Category != token.TOKEN_IDENTIFIER {
                fmt.Println("Expected identifier.") 
                tokenIndex++
                continue
            }
            enumName := tokens[tokenIndex].Source 
            tokenIndex++
            for tokenIndex < len(tokens) && tokens[tokenIndex].Category == token.TOKEN_IDENTIFIER {
                newBlockID := block.CreateBlock()  
                block.Categories[newBlockID] = block.BLOCK_ENUM_VAL
                block.Labels[newBlockID] = prefix + enumName + "_" + tokens[tokenIndex].Source
                block.IntVals[newBlockID] = runningVal
                blocks = append(blocks, newBlockID)
                runningVal++
                tokenIndex++
            }
        } else {
            fmt.Printf("Unexpected token: %s\n", tokens[tokenIndex].Source)
            tokenIndex++
        }
    }
    return blocks
}

