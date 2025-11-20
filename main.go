package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guruorgoru/goru-verbal-interpreter/eval"
	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
	"github.com/guruorgoru/goru-verbal-interpreter/parser"
)

const (
	EXIT   = "exit"
	HELP   = "help"
	PROMPT = ">> "
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(PROMPT)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case EXIT:
			fmt.Println("Exiting interpreter. Goodbye!")
			os.Exit(0)
		case HELP:
			fmt.Println("Available commands:")
			fmt.Println("  help       - Show this help message")
			fmt.Println("  exit       - Exit the Pokedex")
		default:
			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				log.Println(p.Errors())
				continue
			}

			evaluated := eval.Eval(program)
			if evaluated != nil {
				fmt.Println(evaluated.Inspect())
			}
		}

	}
}
