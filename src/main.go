package main

import (
    //"fmt"
    "os"
    "delta/statesy/data/cmd" 
    "delta/statesy/proc/lex"
    "delta/statesy/proc/parse"
    "delta/statesy/proc/build"
)

func main () {
    cmd.SourceFilename = os.Args[1]
    //fmt.Println("lexing . . .")
    tokens := lex.LexFile(cmd.SourceFilename)
    //fmt.Println("parsing . . .")
    blocks := parse.ParseTokens(tokens)
    //fmt.Println("building . . .")
    assembly := build.BuildAssembly(blocks)
    os.WriteFile(cmd.SourceFilename + ".asm", []byte(assembly), 0666)
    header   := build.BuildHeader(blocks)
    os.WriteFile(cmd.SourceFilename + ".h", []byte(header), 0666)
}

