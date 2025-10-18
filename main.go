package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
	"github.com/guruorgoru/goru-verbal-interpreter/token"
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
			fmt.Println("Exiting Pokedex. Goodbye!")
			os.Exit(0)
		case HELP:
			fmt.Println("Available commands:")
			fmt.Println("  help       - Show this help message")
			fmt.Println("  exit       - Exit the Pokedex")
		default:
			l := lexer.New(input)
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		}

	}
}
