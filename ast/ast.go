package ast

import (
	"fmt"
	"lambda/lexer"
)

// ---------- Visitor interface: --------- //

type Visitor interface {
	VisitAbstraction(abs Abstraction) interface{}
	VisitApplication(app Application) interface{}
	VisitIdentifier(id Identifier) interface{}
}

// ---------- Node types: --------- //

type Term interface {
	Accept(v Visitor) interface{}
}

type Abstraction struct {
	Param string
	Body  Term
}

func (abs Abstraction) Accept(v Visitor) interface{} {
	return v.VisitAbstraction(abs)
}

func (abs Abstraction) String() string {
	return fmt.Sprintf("{Abstraction: Param: %s, Body: %s}", abs.Param, abs.Body)
}

type Application struct {
	Left  Term
	Right Term
}

func (app Application) Accept(v Visitor) interface{} {
	return v.VisitApplication(app)
}

func (app Application) String() string {
	return fmt.Sprintf("{Application: Left: %s, Right: %s}", app.Left, app.Right)
}

type Identifier struct {
	Id    lexer.Token
	Index int
}

func (id Identifier) Accept(v Visitor) interface{} {
	return v.VisitIdentifier(id)
}

func (id Identifier) String() string {
	return fmt.Sprintf("{Identifier: %s, Index: %v}", id.Id.Lexeme, id.Index)
}
