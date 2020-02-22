package interpreter

import (
	"lambda/ast"
	//"fmt"
)

type interpreter struct{
	tree ast.Term
}

// Interpreter constructor
func NewInterpreter(tree ast.Term) interpreter {
	i := interpreter{tree}
	return i
}

func (i interpreter) Evaluate() interface{} {
	// First, index all free variables by their depth (with no shift offset)
	result := indexFree(i.tree, 0)
	
	return result.Accept(i)
}

// ---------- Visit methods: ---------- //

func (i interpreter) VisitAbstraction(abs ast.Abstraction) interface{} {
	// Evaluate both the parameter and body.
	pValue := abs.Param
	bValue := abs.Body.Accept(i).(ast.Term)

	// Return an abstraction with it's terms evaluated to values.
	return ast.Abstraction{pValue, bValue}
}

func (i interpreter) VisitApplication(app ast.Application) interface{} {
	lValue := app.Left.Accept(i)
	rValue := app.Right.Accept(i)

	value := substitute(lValue.(ast.Term), rValue.(ast.Term))
	return value
}

func (i interpreter) VisitIdentifier(id ast.Identifier) interface{} {
	// Identifiers are, by default, values. So they all return itself.
	return id
}

// ---------- Helper methods: ---------- //

// Do we even need this?
func isValue(term ast.Term) bool {
	switch term.(type) {

	// Identifiers are values
	case ast.Identifier:
		return true

	// Abstractions with only identifiers in the parameter and body are values
	case ast.Abstraction:
		return isValue(term.(ast.Abstraction).Body)

	default:
		return false
	}
}

// Shifter visitor shifts free variables' indexes by x.
type shifter struct{
	x int
	tree ast.Term
}

func (s shifter) VisitAbstraction(abs ast.Abstraction) interface{} {
	return ast.Abstraction{abs.Param, abs.Body.Accept(s).(ast.Term)}
}

func (s shifter) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(s).(ast.Term)
	right := app.Right.Accept(s).(ast.Term)

	return ast.Application{left, right}
}

func (s shifter) VisitIdentifier(id ast.Identifier) interface{} {
	if id.Free {
		return ast.Identifier{id.Token, id.Index + s.x, id.Free}
	} else {
		return id
	}
}

// Use the shifter visitor to shift all free variables in an AST by 'x'.
func shift(node ast.Term, x int) ast.Term {
	shifter := shifter{x, node}
	return node.Accept(shifter).(ast.Term)
}

// Find the abstraction depth of a given identifier.
type indexer struct {
	depth int
	shift int
}

func (i *indexer) VisitAbstraction(abs ast.Abstraction) interface{} {
	// Increase depth when entering abstraction.
	i.depth++

	body := abs.Body.Accept(i).(ast.Term)

	// Decrease depth when ending abstraction.
	i.depth--

	return ast.Abstraction{abs.Param, body}
}

func (i *indexer) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(i).(ast.Term)
	right := app.Right.Accept(i).(ast.Term)
	return ast.Application{left, right}
}

func (i *indexer) VisitIdentifier(id ast.Identifier) interface{} {
	if id.Free {
		id.Index = i.depth + i.shift
	}
	return id
}

// Use the indexer visitor to index all free variables based on their abstraction depth.
func indexFree(node ast.Term, shift int) ast.Term {
	indexer := &indexer{0, shift}
	result := node.Accept(indexer)
	return result.(ast.Term)
}

type substituter struct {
	value ast.Term
}

func (s substituter) VisitAbstraction(abs ast.Abstraction) interface{} {
	return ast.Abstraction{abs.Param, abs.Body.Accept(s).(ast.Term)}
}

func (s substituter) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(s).(ast.Term)
	right := app.Right.Accept(s).(ast.Term)
	return ast.Application{left, right}
}

func (s substituter) VisitIdentifier(id ast.Identifier) interface{} {
	if id.Index == 0 {
		return s.value
	} else {
		return id
	}
}


// Substitute all variables in 'node' that have a De Bruijn index of zero by 'value'.
// Substitute is defined as the shifted down result of replacing all occurrences of x in node with value shifted up
func substitute(value ast.Term, node ast.Term) ast.Term {
	// First, shift up the value.
	value = shift(value, 1)
	
	// Replace all variables with index 0 with the new value.
	substituter := substituter{value}
	result := node.Accept(substituter).(ast.Term)

	// Re-index free variables and shift them down by one.
	result = indexFree(result, -1)

	return result
}