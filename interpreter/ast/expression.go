package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

//=============================================================================
// Binary Expression
//=============================================================================

type BinaryOperator int // BinaryOperator is an enum for binary operators.

const (
	BinaryOperatorAdd BinaryOperator = iota // e.g 1 + 2
	BinaryOperatorSub BinaryOperator = iota // e.g 1 - 2
	BinaryOperatorMul BinaryOperator = iota // e.g 1 * 2
	BinaryOperatorDiv BinaryOperator = iota // e.g 1 / 2
	BinaryOperatorMod BinaryOperator = iota // e.g 1 % 2
	BinaryOperatorEq  BinaryOperator = iota // e.g 1 == 2
	BinaryOperatorNe  BinaryOperator = iota // e.g 1 != 2
	BinaryOperatorLt  BinaryOperator = iota // e.g 1 < 2
	BinaryOperatorLe  BinaryOperator = iota // e.g 1 <= 2
	BinaryOperatorGt  BinaryOperator = iota // e.g 1 > 2
	BinaryOperatorGe  BinaryOperator = iota // e.g 1 >= 2
)

// BinaryExpression represents a binary expression.
// EXAMPLE: 1 + 2
// Left: 1
// Operator: +
// Right: 2
type BinaryExpression struct {
	Left     Node
	Operator BinaryOperator
	Right    Node
	NodePosition
}

func (b *BinaryExpression) Expression() {}
func (b *BinaryExpression) Type() NodeType {
	return NodeTypeBinaryExpression
}

func (b *BinaryExpression) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	var left = b.Left.EvalSelect(builder, addError)
	var right = b.Right.EvalSelect(builder, addError)

	switch left.(type) {
	case string:
		var l = left.(string)
		switch b.Operator {
		case BinaryOperatorLe:
			return builder.LessEqualThan(l, right)
		case BinaryOperatorLt:
			return builder.LessThan(l, right)
		case BinaryOperatorGe:
			return builder.GreaterEqualThan(l, right)
		case BinaryOperatorGt:
			return builder.GreaterThan(l, right)
		case BinaryOperatorEq:
			return builder.Equal(l, right)
		case BinaryOperatorNe:
			return builder.NotEqual(l, right)
		default:
			addError(fmt.Errorf("Unknown binary operator: %d", b.Operator))
		}
	default:
		addError(fmt.Errorf("Expected string got %T", left))
	}
	return ""
}

func (bin *BinaryExpression) EvalInsert(transpiler *Transpiler) interface{} {
	panic("implement me")
}

//=============================================================================
// Expression Statement
//=============================================================================

// ExpressionStatement represents an expression statement.
// Example (Binary):    1 + 2
// Example (Logical):   1 and 2
// Example (Pair):      priority:<High
// Example (Pair):      priority:High
// Example (Literal):   "high"
// Example (Tag):       +HOME
type ExpressionStatement struct {
	Expr Expression // binary, logical, tag, pair, literal
	NodePosition
}

func (stmt *ExpressionStatement) Type() NodeType {
	return NodeTypeExpressionStatement
}

func (stmt *ExpressionStatement) Statement()  {}
func (stmt *ExpressionStatement) Expression() {}

func (stmt *ExpressionStatement) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return stmt.Expr.EvalSelect(builder, addError)
}

func (stmt *ExpressionStatement) EvalInsert(transpiler *Transpiler) interface{} {
	return stmt.Expr.EvalInsert(transpiler)
}

//=============================================================================
// Literal Expression
//=============================================================================

type LiteralKind int // LiteralType is an enum for literal types.

const (
	LiteralKindString LiteralKind = iota // e.g "buy dog"
	LiteralKindNumber LiteralKind = iota // e.g 5
)

// Literal represents a literal value in the AST.
// Example (string): "buy dog"
// Example (number): 5
type Literal struct {
	Value string
	Kind  LiteralKind
	NodePosition
}

func (l *Literal) Type() NodeType {
	return NodeTypeLiteral
}

func (l *Literal) Expression() {}

func (l *LiteralKind) String() string {
	switch *l {
	case LiteralKindString:
		return "String"
	case LiteralKindNumber:
		return "Number"
	}
	return "Unknown"
}

func (l *Literal) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return l.ToValue(nil)
}

func (l *Literal) ToValue(transpiler *Transpiler) interface{} {
	if l.Kind == LiteralKindString {
		return l.Value
	}
	if strings.Contains(l.Value, ".") {
		var fl, err = strconv.ParseFloat(l.Value, 64)
		if err != nil {
			return transpiler.AddError(fmt.Errorf("Failed to parse float: %s %w ", l.Value, err), l)
		}
		return fl
	}
	var in, err = strconv.ParseInt(l.Value, 10, 64)
	if err != nil {
		return transpiler.AddError(fmt.Errorf("Failed to parse int: %s %w ", l.Value, err), l)
	}
	return in
}

func (lit *Literal) EvalInsert(transpiler *Transpiler) interface{} {
	transpiler.AddValue(lit.ToValue(transpiler))
	return nil
}

//=============================================================================
// Logical Expression
//=============================================================================

type LogicalOperator int // LogicalOperator is an enum for logical operators.

const (
	LogicalOperatorAnd LogicalOperator = iota // e.g 1 and 2
	LogicalOperatorOr  LogicalOperator = iota // e.g 1 or 2
)

// LogicalExpression represents a logical expression.
// EXAMPLE (and):   1 and 2
// EXAMPLE (or):    1 or 2
type LogicalExpression struct {
	Left     Node
	Operator LogicalOperator
	Right    Node
	NodePosition
}

func (l *LogicalExpression) Type() NodeType {
	return NodeTypeLogicalExpression
}

func (l *LogicalExpression) Expression() {}

func (l *LogicalExpression) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	var left = l.Left.EvalSelect(builder, addError)
	var right = l.Right.EvalSelect(builder, addError)

	switch left.(type) {
	case string:
		switch l.Operator {
		case LogicalOperatorAnd:
			return builder.And(left.(string), right.(string))
		case LogicalOperatorOr:
			return builder.Or(left.(string), right.(string))
		default:
			return addError(fmt.Errorf("Unknown logical operator: %d", l.Operator))
		}
	}
	return addError(fmt.Errorf("Expected string got %T", left))
}

func (logical *LogicalExpression) EvalInsert(transpiler *Transpiler) interface{} {
	panic("implement me")
}

// =============================================================================
// Key
// =============================================================================
type Key struct {
	Key  string
	Expr Expression
	NodePosition
}

func (key *Key) Expression() {}
func (key *Key) Type() NodeType {
	return NodeTypeBinaryExpression
}

func (key *Key) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return key.Expr.EvalSelect(builder, addError)
}

func (key *Key) EvalInsert(transpiler *Transpiler) interface{} {
	var lowerK = strings.ToLower(key.Key)
	switch lowerK {
	case "priority":
		transpiler.AddCol("priority")
		return key.Expr.EvalInsert(transpiler)
	default:
		return transpiler.AddError(fmt.Errorf("Unknown key: %s", key.Key), key)
	}
}
