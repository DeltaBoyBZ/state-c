package token

type TokenCategory int8 
type NumberCategory int8

const (
    TOKEN_INT TokenCategory = iota
    TOKEN_FLOAT
    TOKEN_CHAR
    TOKEN_ARRAY
    TOKEN_IDENTIFIER
    TOKEN_EQUALS
    TOKEN_LITERAL_INT
    TOKEN_LITERAL_FLOAT
    TOKEN_LITERAL_STRING
    TOKEN_PREFIX
    TOKEN_TABLE
    TOKEN_DUMP
    TOKEN_ENUM
)

const (
    NUMBER_INT NumberCategory = iota
    NUMBER_FLOAT
)

var TokenString = map[TokenCategory]string {
    TOKEN_INT: "int",
    TOKEN_FLOAT: "float",
    TOKEN_CHAR: "char", 
    TOKEN_EQUALS: "=",
    TOKEN_ARRAY: "array",
    TOKEN_PREFIX: "prefix",
    TOKEN_TABLE: "table",
    TOKEN_DUMP: "dump",
    TOKEN_ENUM: "enum",
}

type Token struct {
    Category TokenCategory
    Source string
    IntVal int
    FloatVal float32
    StringVal string
}

