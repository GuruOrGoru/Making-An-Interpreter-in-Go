# Goru Verbal Interpreter

A simple interpreter for the Goru Verbal programming language, implemented in Go. This project is based on the book "Writing an Interpreter in Go" by Thorsten Ball, but features Hindi keywords for a unique twist.

## Features

Goru Verbal is a dynamically-typed programming language with the following features:

- **Variables**: Declare variables with `manau`
- **Functions**: Define functions with `karya`
- **Conditionals**: Use `yadi` (if) and `natra` (else)
- **Booleans**: `satya` (true) and `jhuth` (false)
- **Return statements**: `firta`
- **Arithmetic operations**: `+`, `-`, `*`, `/`
- **Comparison operators**: `<`, `>`, `==`, `!=`
- **Integer literals**

## Installation

Clone the repository and ensure you have Go installed (version 1.24.6 or later).

```bash
git clone https://github.com/guruorgoru/goru-verbal-interpreter.git
cd goru-verbal-interpreter
```

## Building

To build the interpreter:

```bash
make build
```

This creates an executable at `bin/app`.

## Running

To run the interpreter in REPL mode:

```bash
make run
```

Or directly:

```bash
go run main.go
```

## Usage

The interpreter provides an interactive REPL (Read-Eval-Print Loop). Type commands and press Enter to execute them.

### Commands

- `help` - Show available commands
- `exit` - Exit the interpreter

### Examples

```goru-verbal
>> manau x = 5;
>> manau y = 10;
>> manau add = karya(a, b) { a + b; };
>> add(x, y);
15
>> yadi (x < y) { firta satya; } natra { firta jhuth; };
satya
```

## Project Structure

- `ast/` - Abstract Syntax Tree definitions
- `eval/` - Expression evaluator
- `lexer/` - Lexical analyzer
- `object/` - Runtime object system
- `parser/` - Parser for the language
- `token/` - Token definitions

## Testing

Run the test suite:

```bash
make test
```

Or:

```bash
go test ./...
```

## Contributing

project is open source. Check the license file for details.
Update on the interpreter: Just need to add more datatypes support
