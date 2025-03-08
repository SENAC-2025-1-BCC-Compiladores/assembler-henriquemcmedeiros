package lexer

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	TOKEN_SECTION  = "SECTION"
	TOKEN_EOF      = "EOF"
	TOKEN_INSTR    = "INSTRUCTION"
	TOKEN_NUMBER   = "NUMBER"
	TOKEN_VAR      = "VARIABLE"
	TOKEN_DEFINE   = "DEFINE"
	TOKEN_UNKNOWN  = "UNKNOWN"
)

var (
	Instructions = map[string]uint8{
		"NOP": 0x00, "STA": 0x10, "LDA": 0x20, "ADD": 0x30,
		"OR": 0x40, "AND": 0x50, "NOT": 0x60, "JMP": 0x80,
		"JN": 0x90, "JZ": 0xA0, "HLT": 0xF0,
	}

	Define = map[string]bool{
		"DB": true, "DS": false, "ORG": true,
	}

	varRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

type Token struct {
	Tipo  string
	Valor string
}

func isInstruction(lexema string) bool {
	_, existe := Instructions[lexema]
	return existe
}

func isDefine(lexema string) bool {
	_, existe := Define[lexema]
	return existe
}

func isNumber(lexema string) bool {
	if _, err := strconv.ParseInt(lexema, 16, 64); err == nil {
		return true
	}
	return false
}

func isVariable(lexema string) bool {
	return varRegex.MatchString(lexema)
}

func lexer(lexema string) Token {
	switch {
	case strings.HasPrefix(lexema, "."):
		return Token{Tipo: TOKEN_SECTION, Valor: strings.TrimPrefix(lexema, ".")}
	case isInstruction(lexema):
		return Token{Tipo: TOKEN_INSTR, Valor: lexema}
	case isDefine(lexema):
		return Token{Tipo: TOKEN_DEFINE, Valor: lexema}
	case isNumber(lexema):
		return Token{Tipo: TOKEN_NUMBER, Valor: lexema}
	case isVariable(lexema):
		return Token{Tipo: TOKEN_VAR, Valor: lexema}
	default:
		log.Printf("Token desconhecido: %s", lexema)
		return Token{Tipo: TOKEN_UNKNOWN, Valor: lexema}
	}
}

func GetTokens(caminhoArquivo string) (tokens []Token) {
	arquivo, err := os.ReadFile(caminhoArquivo)
	if err != nil {
		log.Fatalf("Não foi possível ler o arquivo: %v", err)
	}

	linhas := strings.Split(string(arquivo), "\n")

	for _, linha := range linhas {
		linha = strings.Split(linha, ";")[0] // Remove comentários
		re := regexp.MustCompile(`\s+`)
		lexemas := re.Split(strings.TrimSpace(linha), -1)

		for _, lexema := range lexemas {
			if lexema != "" {
				tokens = append(tokens, lexer(lexema))
			}
		}
	}

	tokens = append(tokens, Token{Tipo: TOKEN_EOF, Valor: ""})

	return
}
