package interpreter

import (
	"fmt"
	"lambda/ast"
	"lambda/environment"
)

type interpreter struct {
	definitions []ast.Definition
	expression  ast.Term
	env         environment.Environment
}

// Interpreter constructor
func NewInterpreter(defs []ast.Definition, expr ast.Term) interpreter {
	i := interpreter{}
	i.definitions = defs
	i.expression = expr
	i.env = environment.NewEnvironment()
	return i
}

func (i interpreter) Evaluate() interface{} {

	// Define all definitions
	for _, def := range i.definitions {
		def.Accept(i)
	}

	// Evaluate the expression.
	result := i.expression.Accept(i).(ast.Term)

	// Remove any ' from identifier names.
	result = reverseAlpha(result)

	// See if the result has a definition in the environment.
	key, ok := i.env.Lookup(result)
	if ok {
		return key
	}

	return i.expression.Accept(i)
}

// ---------- Interpreter visit methods: ---------- //

func (i interpreter) VisitDefinition(def ast.Definition) interface{} {
	i.env.Define(def.Id.Name, def.Term)
	return nil
}

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
	expr := ast.Application{lValue.(ast.Term), rValue.(ast.Term)}

	// If the expression can be evaluated further, begin evaluating.
	result := i.substitute(expr.Left, expr.Right)

	// Values evaluate to themselves, so only continue evaluating non-values
	if result != alpha(expr) {
		return result.(ast.Term).Accept(i)
	}

	return expr
}

func (i interpreter) VisitIdentifier(id ast.Identifier) interface{} {
	// Identifiers are, by default, values. So they all return itself.

	// Except for identifiers bound to a definition.
	term, ok := i.env.Get(id.Name)
	if ok {
		return term
	}

	return id
}

func (i interpreter) substitute(left ast.Term, right ast.Term) interface{} {
	// Alpha conversion to resolve duplicate variables
	alphaExpr := alpha(ast.Application{left, right})

	// Beta reduction to evaluate the application
	switch left := alphaExpr.Left.(type) {
	case ast.Identifier:
		return ast.Application{left, alphaExpr.Right}
	case ast.Abstraction:
		beta := &beta{left.Param.(ast.Identifier).Name, alphaExpr.Right}
		betaExpr := left.Body.Accept(beta).(ast.Term)
		return betaExpr.Accept(i)
	case ast.Application:
		return ast.Application{left.Accept(i).(ast.Term), alphaExpr.Right}
	}

	return nil
}

// ---------- Alpha Conversion Visitor: ---------- //

type alphaVisitor struct {
	variables []ast.Identifier
	left      bool
}

func (a *alphaVisitor) VisitDefinition(def ast.Definition) interface{} {
	return nil
}

func (a *alphaVisitor) VisitAbstraction(abs ast.Abstraction) interface{} {
	pValue := abs.Param.Accept(a).(ast.Term)
	bValue := abs.Body.Accept(a).(ast.Term)

	return ast.Abstraction{pValue, bValue}
}

func (a *alphaVisitor) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(a).(ast.Term)
	right := app.Right.Accept(a).(ast.Term)
	return ast.Application{left, right}
}

func (a *alphaVisitor) VisitIdentifier(id ast.Identifier) interface{} {
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

func alpha(app ast.Application) ast.Application {
	a := &alphaVisitor{}

	// The left flag sets the visitor to "append mode" and only appends variables.
	a.left = true
	lValue := app.Left.Accept(a)

	// With the left flag false, the visitor actually converts duplicate variables.
	a.left = false
	rValue := app.Right.Accept(a)

	return ast.Application{lValue.(ast.Term), rValue.(ast.Term)}
}

// ---------- Alpha Reversal Visitor: ---------- //

type alphaReverser struct{}

func (a alphaReverser) VisitDefinition(def ast.Definition) interface{} {
	return nil
}

func (a alphaReverser) VisitAbstraction(abs ast.Abstraction) interface{} {
	pValue := abs.Param.Accept(a).(ast.Term)
	bValue := abs.Body.Accept(a).(ast.Term)

	return ast.Abstraction{pValue, bValue}
}

func (a alphaReverser) VisitApplication(app ast.Application) interface{} {
	left := app.Left.Accept(a).(ast.Term)
	right := app.Right.Accept(a).(ast.Term)
	return ast.Application{left, right}
}

func (a alphaReverser) VisitIdentifier(id ast.Identifier) interface{} {
	name := id.Name
	for name[len(name)-1:] == "'" {
		name = name[:len(name)-1]
	}

	return ast.Identifier{name}
}

// Essentially undoes the alpha conversion process.
func reverseAlpha(term ast.Term) ast.Term {
	a := alphaReverser{}

	return term.Accept(a).(ast.Term)
}

// ---------- Beta Reduction Visitor: ---------- //
type beta struct {
	parameter string
	value     ast.Term
}

func (b *beta) VisitDefinition(def ast.Definition) interface{} {
	return nil
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
