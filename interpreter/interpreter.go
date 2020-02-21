package interpreter

import "lambda/ast"

type Interpreter struct{}

// ---------- Visit methods: ---------- //

func (i Interpreter) VisitAbstraction(abs ast.Abstraction) interface{} {
	// Evaluate both the parameter and body.
	pValue := abs.Param.Accept(i).(ast.Term)
	bValue := abs.Body.Accept(i).(ast.Term)

	// Return an abstraction with it's terms evaluated to values.
	return ast.Abstraction{pValue, bValue}
}

func (i Interpreter) VisitApplication(app ast.Application) interface{} {
	//lValue := app.Left.Accept(i)
	//rValue := app.Right.Accept(i)

	//value := substitute(rValue, lValue.(ast.Abstraction).Body)
	return nil
}

func (i Interpreter) VisitIdentifier(id ast.Identifier) interface{} {
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
		return isValue(term.(ast.Abstraction).Param) && isValue(term.(ast.Abstraction).Body)

	default:
		return false
	}
}
