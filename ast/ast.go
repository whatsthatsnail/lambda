package ast

import (
	"fmt"
)

// ---------- Visitor interface: --------- //

type Visitor interface {
	VisitDefinition(def Definition) interface{}

	VisitAbstraction(abs Abstraction) interface{}
	VisitApplication(app Application) interface{}
	VisitIdentifier(id Identifier) interface{}
}

// ---------- Node types: --------- //

type Definition struct {
	Id   Identifier
	Term Term
}

func (def Definition) Accept(v Visitor) interface{} {
	return v.VisitDefinition(def)
}

func (def Definition) String() string {
	return fmt.Sprintf("let %s = %s\n", def.Id, def.Term)
}

type Term interface {
	Accept(v Visitor) interface{}
}

type Abstraction struct {
	Param Term
	Body  Term
}

func (abs Abstraction) Accept(v Visitor) interface{} {
	return v.VisitAbstraction(abs)
}

func (abs Abstraction) String() string {
	return fmt.Sprintf("(λ%s. %s)", abs.Param, abs.Body)
}

type Application struct {
	Left  Term
	Right Term
}

func (app Application) Accept(v Visitor) interface{} {
	return v.VisitApplication(app)
}

func (app Application) String() string {
	return fmt.Sprintf("(%s %s)", app.Left, app.Right)
}

type Identifier struct {
	Name string
}

func (id Identifier) Accept(v Visitor) interface{} {
	return v.VisitIdentifier(id)
}

func (id Identifier) String() string {
	return id.Name
}
