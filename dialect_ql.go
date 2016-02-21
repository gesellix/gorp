// Copyright 2012 James Cooper. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package gorp provides a simple way to marshal Go structs to and from
// SQL databases.  It uses the database/sql package, and should work with any
// compliant database/sql driver.
//
// Source code and project home:
// https://github.com/go-gorp/gorp

package gorp

import (
	"fmt"
	"reflect"
	"strings"
)

type QlDialect struct {
}

func (d QlDialect) QuerySuffix() string {
	return ";"
}

func (d QlDialect) CreateTableSuffix() string {
	return ""
}

func (d QlDialect) CreateIndexSuffix() string {
	return ""
}

func (d QlDialect) DropIndexSuffix() string {
	return ""
}

func (d QlDialect) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	switch val.Kind() {
	case reflect.Ptr:
		return d.ToSqlType(val.Elem(), maxsize, isAutoIncr)
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "int32"
	case reflect.Int64, reflect.Uint64:
		return "int64"
	case reflect.Float64:
		return "float64"
	case reflect.Float32:
		return "float32"
	}

	switch val.Name() {
	case "NullInt64":
		return "bigint"
	case "NullFloat64":
		return "float64"
	case "NullBool":
		return "bool"
	case "Time", "NullTime":
		return "time"
	}

	if maxsize > 0 {
		return fmt.Sprintf("varchar(%d)", maxsize)
	} else {
		return "string"
	}

}

func (d QlDialect) AutoIncrStr() string {
	return ""
}

func (d QlDialect) AutoIncrBindValue() string {
	return ""
}

func (d QlDialect) AutoIncrInsertSuffix(col *ColumnMap) string {
	return " returning " + d.QuoteField(col.ColumnName)
}

func (d QlDialect) TruncateClause() string {
	return "TRUNCATE"
}

// Returns "$(i+1)"
func (d QlDialect) BindVar(i int) string {
	return fmt.Sprintf("$%d", i + 1)
}

func (d QlDialect) QuoteField(f string) string {
	return f
}

func (d QlDialect) QuotedTableForQuery(schema string, table string) string {
	if strings.TrimSpace(schema) == "" {
		return d.QuoteField(table)
	}

	return schema + "." + d.QuoteField(table)
}

func (d QlDialect) IfSchemaNotExists(command, schema string) string {
	return fmt.Sprintf("%s IF NOT EXISTS", command)
}

func (d QlDialect) IfTableExists(command, schema, table string) string {
	return fmt.Sprintf("%s IF EXISTS", command)
}

func (d QlDialect) IfTableNotExists(command, schema, table string) string {
	return fmt.Sprintf("%s IF NOT EXISTS", command)
}
