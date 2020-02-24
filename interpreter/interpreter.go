package interpreter

import (
	"fmt"
	"lambda/ast"
)

type interpreter struct {
	tree ast.Term
}

// Interpreter constructor
func NewInterpreter(tree ast.Term) interpreter {
	i := interpreter{tree}
	return i
}

func (i interpreter) Evaluate() interface{} {
	// First, index all free variables by their depth (with no shift offset)
	fmt.Println(i.tree)

	return i.tree.Accept(i)
}

// ---------- Interpreter visit methods: ---------- //

func (i interpreter) VisitAbstraction(abs ast.Abstraction) interface{} {
	// Evaluate both the parameter and body.
	pValue := abs.Param.Accept(i).(ast.Term)
	bValue := abs.Body.Accept(i).(ast.Term)

	// Return an abstraction with it's terms evaluated to values.
	return ast.Abstraction{pValue, bValue}
}

func (i interpreter) VisitApplication(app ast.Application) interface{} {
	lValue := app.Left.Accept(i)
	rValue := app.Right.Accept(i)
	value := ast.Application{lValue.(ast.Term), rValue.(ast.Term)}

	// If the expression can be evaluated further, begin evaluating.
	isValue := isValue{}
	if !value.Accept(isValue).(bool) {
		// Alpha conversion to resolve duplicate variables
		alpha := &alpha{}
		value = alpha.conv(value)

		// Beta reduction to evaluate the application
		switch v := value.Left.(type) {
		case ast.Identifier:
			return ast.Application{v, value.Right}
		case ast.Abstraction:
			beta := &beta{v.Param.(ast.Identifier).Name, value.Right}
			result := v.Body.Accept(beta)
			return result
		case ast.Application:
			return v.Accept(i)
		}
	}

	return value
}

func (i interpreter) VisitIdentifier(id ast.Identifier) interface{} {
	// Identifiers are, by default, values. So they all return itself.
	return id
}

// ---------- Alpha Conversion Visitor: ---------- //

type alpha struct {
	variables []ast.Identifier
	left      bool
}

func (a *alpha) VisitAbstraction(abs ast.Abstraction) interface{} {
	pValue := abs.Param.Accept(a).(ast.Term)
	bValue := abs.Body.Accept(a).(ast.Term)

	return ast.Abstraction{pValue, bValue}
}

func (a *alpha) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(a).(ast.Term)
	right := app.Right.Accept(a).(ast.Term)
	return ast.Application{left, right}
}

func (a *alpha) VisitIdentifier(id ast.Identifier) interface{} {
	if a.left {
		a.variables = append(a.variables, id)
		return id
	}

	for _, k := range a.variables {
		if k.Name == id.Name {
			id.Name = fmt.Sprintf("%s%s", id.Name, "'")
			return id
		}
	}

	return id
}

func (a *alpha) conv(app ast.Application) ast.Application {
	// The left flag sets the visitor to "append mode" and only appends variables.
	a.left = true
	lValue := app.Left.Accept(a)

	// With the left flag false, the visitor actually converts duplicate variables.
	a.left = false
	rValue := app.Right.Accept(a)

	return ast.Application{lValue.(ast.Term), rValue.(ast.Term)}
}

// ---------- Beta Reduction Visitor: ---------- //
type beta struct {
	parameter string
	value     ast.Term
}

func (b *beta) VisitAbstraction(abs ast.Abstraction) interface{} {
	pValue := abs.Param.Accept(b).(ast.Term)
	bValue := abs.Body.Accept(b).(ast.Term)

	return ast.Abstraction{pValue, bValue}
}

func (b *beta) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(b).(ast.Term)
	right := app.Right.Accept(b).(ast.Term)
	return ast.Application{left, right}
}

func (b *beta) VisitIdentifier(id ast.Identifier) interface{} {
	if id.Name == b.parameter {
		return b.value
	}

	return id
}

// ---------- isValue visitor ---------- //

// Returns true if the expression is a value and cannot be evaluated further.
type isValue struct {}

func (v isValue) VisitAbstraction(abs ast.Abstraction) interface{} {
	return abs.Body.Accept(v)
}

func (v isValue) VisitApplication(app ast.Application) interface{} {
	var left bool
	switch l := app.Left.(type) {
	case ast.Identifier:
		left = true
	case ast.Abstraction:
		left = false
	case ast.Application:
		left = l.Accept(v).(bool)
	}
	right := app.Right.Accept(v).(bool)

	if (left && right) {
		return true
	}
	
	return false
}

func (v isValue) VisitIdentifier(id ast.Identifier) interface{} {
	return true
}