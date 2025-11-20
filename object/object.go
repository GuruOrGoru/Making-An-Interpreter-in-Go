package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ         = "INTEGER"
	BOOLEAN_OBJ         = "BOOLEAN"
	NULL_OBJ            = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJ           = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// INTEGER

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// INTEGER END

// -------------- //

// BOOLEAN

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// BOOLEAN END

// -------------- //

// NULL

type Null struct {
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// NULL END

// ---------- //

// RETURN VALUE

type ReturnValue struct {
	Value Object
}

func (rb *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

func (rb *ReturnValue) Inspect() string {
	return rb.Value.Inspect()
}

// RETURN VALUE END

// ---------- //

// ERROR

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// ERROR END

// ---------- //
