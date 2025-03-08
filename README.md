### Henrique Marques de Carvalho Medeiros

# Assembler Simples

Este é um projeto de um assembler simples escrito em Go. Ele converte um arquivo de código assembly em um arquivo binário `.mem`.

## Como Compilar e Rodar

### Pré-requisitos

- Go instalado na máquina (versão 1.16 ou superior).

### Passos para Compilar e Executar

1. **Clone o repositório** (se aplicável) ou navegue até o diretório onde o arquivo `main.go` está localizado.

2. **Compile o programa** usando o comando abaixo:

   ```bash
   go build main.go
   ```
   
3. **Rode o Programa**
   
    ```bash
     ./main <arquivo.asm>
     ```
4. **Saída**:
    - O assembler gerará um arquivo output.mem no mesmo diretório onde foi executado.
    - Este arquivo contém o código binário correspondente ao assembly fornecido.
