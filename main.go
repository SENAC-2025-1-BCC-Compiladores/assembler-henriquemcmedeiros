package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"FullGPT/lexer"
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
)

type Assembler struct {
	Tokens []lexer.Token
	PC     uint8
	Output []uint8
	Labels map[string]uint8
}

func NewAssembler(tokens []lexer.Token) *Assembler {
	return &Assembler{
		Tokens: tokens,
		PC:     0,
		Output: make([]uint8, 0),
		Labels: make(map[string]uint8),
	}
}

func (a *Assembler) FirstPass() error {
	for i := 0; i < len(a.Tokens); i++ {
		token := a.Tokens[i]

		switch token.Tipo {
		case TOKEN_VAR:
			a.Labels[token.Valor] = a.PC
		case TOKEN_INSTR:
			a.PC += 2
		case TOKEN_NUMBER:
			a.PC += 2
		case TOKEN_DEFINE:
			switch token.Valor {
			case "DB":
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após DB")
				}
				a.PC += 2
			case "ORG":
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após ORG")
				}
				value, err := parseNumber(a.Tokens[i].Valor)
				if err != nil {
					return fmt.Errorf("número inválido após ORG: %s", a.Tokens[i].Valor)
				}
				a.PC = uint8(value)
			}
		}
	}
	return nil
}

func (a *Assembler) SecondPass() error {
	for i := 0; i < len(a.Tokens); i++ {
		token := a.Tokens[i]

		switch token.Tipo {
		case TOKEN_INSTR:
			opcode, ok := Instructions[token.Valor]
			if !ok {
				return fmt.Errorf("instrução desconhecida: %s", token.Valor)
			}
			a.Output = append(a.Output, opcode, 0x00)
		case TOKEN_NUMBER:
			value, err := parseNumber(token.Valor)
			if err != nil {
				return fmt.Errorf("número inválido: %s", token.Valor)
			}
			a.Output = append(a.Output, uint8(value), 0x00)
		case TOKEN_VAR:
			address, ok := a.Labels[token.Valor]
			if !ok {
				return fmt.Errorf("label não definido: %s", token.Valor)
			}
			a.Output = append(a.Output, address, 0x00)
		case TOKEN_DEFINE:
			switch token.Valor {
			case "DB":
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após DB")
				}
				value, err := parseNumber(a.Tokens[i].Valor)
				if err != nil {
					return fmt.Errorf("número inválido após DB: %s", a.Tokens[i].Valor)
				}
				a.Output = append(a.Output, uint8(value), 0x00)
			case "ORG":
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após ORG")
				}
			}
		}
	}
	return nil
}

func parseNumber(s string) (uint64, error) {
	if strings.HasPrefix(s, "0x") {
		return strconv.ParseUint(s[2:], 16, 8)
	}
	return strconv.ParseUint(s, 10, 8)
}

func (a *Assembler) WriteMEM(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := []uint8{0x03, 0x4E, 0x44, 0x52}
	output := append(header, a.Output...)

	for len(output) < 516 {
		output = append(output, 0x00)
	}

	_, err = file.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run main.go <arquivo.asm>")
	}

	asmFile := os.Args[1]
	if _, err := os.Stat(asmFile); os.IsNotExist(err) {
		log.Fatalf("Arquivo de entrada não encontrado: %s", asmFile)
	}

	tokens := lexer.GetTokens(asmFile)

	assembler := NewAssembler(tokens)

	err := assembler.FirstPass()
	if err != nil {
		log.Fatalf("Erro na primeira passagem: %v", err)
	}

	err = assembler.SecondPass()
	if err != nil {
		log.Fatalf("Erro na segunda passagem: %v", err)
	}

	err = assembler.WriteMEM("output.mem")
	if err != nil {
		log.Fatalf("Erro ao escrever arquivo .mem: %v", err)
	}

	fmt.Println("Arquivo .mem binário gerado com sucesso!")
}