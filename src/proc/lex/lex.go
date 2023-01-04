package lex

import (
    //"fmt"
    "os"
    "unicode"
    "strconv"
    "delta/statesy/data/token"
)

func sourceEq (source []byte, index int, x string) bool {
    return string(source[index:index+len(x)]) == x
}

func skipWhitespace(source []byte, index_ptr *int) {
    sourceLen := len(source)
    for *index_ptr < sourceLen && unicode.IsSpace(rune(source[*index_ptr])) {
        *index_ptr++ 
    }
}

func isIdentifierByte(x byte) bool {
    if unicode.IsLetter(rune(x)) { return true }
    if x == '_' { return true }
    return false
}

func isNumberByte(x byte) bool {
    if unicode.IsDigit(rune(x)) { return true }
    if x == '.' { return true }
    return false
}

func checkFixedToken (source []byte, index_ptr *int, category token.TokenCategory, tokens []token.Token) ([]token.Token, bool) {
    if sourceEq(source, *index_ptr, token.TokenString[category]) {
        var newToken token.Token
        newToken.Category = category
        newToken.Source   = token.TokenString[category]
        tokens = append(tokens, newToken)
        *index_ptr += len(token.TokenString[category])
        skipWhitespace(source, index_ptr)
        return tokens, true
    } 
    return tokens, false
}

func LexFile (filename string) []token.Token {
    source, _ := os.ReadFile(filename)
    sourceIndex := 0 
    sourceLen   := len(source)
    tokens := make([]token.Token, 0, 32)
    skipWhitespace(source, &sourceIndex)
    for sourceIndex < sourceLen {
        initLen := len(tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_INT, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_FLOAT, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_EQUALS, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_ARRAY, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_PREFIX, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_TABLE, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_DUMP, tokens)
        tokens, _  = checkFixedToken(source, &sourceIndex, token.TOKEN_ENUM, tokens)
        if len(tokens) > initLen { continue }
        if unicode.IsLetter(rune(source[sourceIndex])) {
            idLen := 0
            for sourceIndex + idLen < sourceLen && isIdentifierByte(source[sourceIndex+idLen]) { 
                idLen++
            }
            identifierSource := source[sourceIndex:sourceIndex+idLen]
            var newToken token.Token
            newToken.Category = token.TOKEN_IDENTIFIER
            newToken.Source   = string(identifierSource)
            tokens = append(tokens, newToken)
            sourceIndex += idLen
            skipWhitespace(source, &sourceIndex)
            continue
        } 
        if unicode.IsDigit(rune(source[sourceIndex])) {
            numberCategory := token.NUMBER_INT
            litLen := 0
            for sourceIndex + litLen < sourceLen && isNumberByte(source[sourceIndex+litLen]) {
                if rune(source[sourceIndex+litLen]) == '.' {
                    numberCategory = token.NUMBER_FLOAT
                }
                litLen++
            }
            var newToken token.Token
            newToken.Source = string(source[sourceIndex:sourceIndex+litLen])
            switch numberCategory {
                case token.NUMBER_INT:
                    newToken.Category = token.TOKEN_LITERAL_INT
                    newToken.IntVal, _ = strconv.Atoi(newToken.Source)
                    break
                case token.NUMBER_FLOAT:
                    newToken.Category = token.TOKEN_LITERAL_FLOAT
                    x, _ := strconv.ParseFloat(newToken.Source, 32)
                    newToken.FloatVal = float32(x)
                    break
                default:
            }
            tokens = append(tokens, newToken)
            sourceIndex += litLen
            skipWhitespace(source, &sourceIndex)
            continue
        } 
        if source[sourceIndex] == '"' {
            sourceIndex++
            strLen := 0
            for source[sourceIndex+strLen] != '"' {
                strLen++
            }
            var newToken token.Token 
            newToken.Category = token.TOKEN_LITERAL_STRING
            newToken.StringVal = string(source[sourceIndex:sourceIndex+strLen]) 
            newToken.Source = "\"" + newToken.StringVal + "\""
            tokens = append(tokens, newToken)
            sourceIndex += strLen + 1
        }
        sourceIndex++
    }
    return tokens
}

