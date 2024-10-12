package transpiler

import "github.com/huandu/go-sqlbuilder"

type SqlBuilder struct {
  Select *sqlbuilder.SelectBuilder
  Insert *sqlbuilder.InsertBuilder
  Delete *sqlbuilder.DeleteBuilder
}

func NewSqlBuilder() *SqlBuilder {
  return &SqlBuilder{
    Select: sqlbuilder.NewSelectBuilder(),
    Insert: sqlbuilder.NewInsertBuilder(),
    Delete: sqlbuilder.NewDeleteBuilder(),
  }
}
