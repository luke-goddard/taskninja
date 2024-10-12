package transpiler

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/luke-goddard/taskninja/interpreter/ast"
)

type table string
type joinExpression string

// Joinable is an interface that is used to evaluate
// the queries joins and convert them into SQL joint statements.
type Joinable interface {
	Join(table string, onExpr ...string) Joinable
}

type JoinTranspiler struct {
	sql            *sqlbuilder.SelectBuilder
	previouslySeen map[table]map[joinExpression]bool
}

// NewJoinTranspiler creates a new JoinTranspiler instance.
// This instance make sure that joins are unique and not repeated.
func NewJoinTranspiler(builder *sqlbuilder.SelectBuilder) *JoinTranspiler {
	return &JoinTranspiler{
		sql:            builder,
		previouslySeen: make(map[table]map[joinExpression]bool),
	}
}

// Join adds a join statement to the SQL query.
func (j *JoinTranspiler) Join(table string, onExpr ...string) Joinable {
	for _, expr := range onExpr {
		j.sql.Join(table, expr)
	}
	return j
}

func (j *JoinTranspiler) join(table table, onExpr joinExpression) bool {
	if !j.hasSeen(table, joinExpression(onExpr)) {
		j.sql.Join(string(table), string(onExpr))
		return true
	}
	return false
}

func (j *JoinTranspiler) hasSeen(table table, expr joinExpression) bool {
	var ret = false
	if _, ok := j.previouslySeen[table]; !ok {
		j.previouslySeen[table] = make(map[joinExpression]bool)
	}
	if _, ok := j.previouslySeen[table][expr]; !ok {
		j.previouslySeen[table][expr] = true
	} else {
		ret = true
	}
	return ret
}

func (j *JoinTranspiler) reset() {
	j.sql = sqlbuilder.NewSelectBuilder()
	j.previouslySeen = make(map[table]map[joinExpression]bool)
}

func (j *JoinTranspiler) Visit(node ast.Node) *JoinTranspiler {
	// switch node := node.(type) {
	// case *ast.Tag:
	// 	j.Join("tags", "tags.id = taskTags.tagId")
	// }
	return j
}
