package generator

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var ErrDriver = errors.New("未匹配到数据库类型")
var ErrConnect = errors.New("连接数据库错误")
var ErrQuery = errors.New("查询数据库错误")

type ColumnSchema struct {
	TableName              string
	ColumnName             string
	DataType               string
	UdtName                string
	IsNullable             string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
}

// PgType 数据库字段类型转go类型
func PgType(c *ColumnSchema) (string, string) {
	requiredImport := ""
	var gt string
	switch c.DataType {
	case "char", "character varying", "varchar", "text", "character":
		gt = "string"
	case "date", "time", "timestamp", "timestampz", "timestamp without time zone", "timestamp with time zone":
		gt, requiredImport = "*time.Time", "time"
	case "smallint", "integer", "bigint", "serial", "big serial", "interval":
		gt = "int64"
	case "float", "double precision", "decimal", "numeric", "money", "real":
		gt = "float64"
	case "bytea", "tsvector", "USER-DEFINED", "uuid", "inet":
		gt = "string"
	case "ARRAY":
		cCopy := &ColumnSchema{DataType: c.UdtName}
		gt, _ = PgType(cCopy)
	case "json", "jsonb":
		gt = "yutil.JsonB"
	case "boolean":
		gt = "bool"
	case "_int2", "_int4", "_int8":
		gt = "pq.Int64Array"
	case "_text", "_uuid", "_varchar", "_char", "_bpchar":
		gt = "pq.StringArray"
	case "_bool":
		gt = "pq.BoolArray"
	case "_jsonb":
		gt = "pq.StringArray"
	}
	if gt == "" {
		return "", ""
	}
	return gt, requiredImport
}

func MySQLType(col *ColumnSchema) (string, string) {
	requiredImport := ""

	var gt string = ""
	switch col.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":

		gt = "string"

	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":

		gt, requiredImport = "*time.Time", "time"
	case "smallint", "int", "mediumint", "bigint":

		gt = "int64"

	case "float", "decimal", "double":

		gt = "float64"

	case "tinyint":
		gt = "int"
	}
	if gt == "" {
		return "", ""
	}
	return gt, requiredImport
}

// getPgSchema 查询pg数据库结构
func getPgSchema(dbAddr, dbName string) map[string][]*ColumnSchema {
	tables := make(map[string][]*ColumnSchema)
	datasource := fmt.Sprintf("%s dbname=%s", dbAddr, dbName)
	conn, err := sql.Open("postgres", datasource)
	if err != nil {
		PrintErr(ErrConnect)
		panic(err)
	}
	defer conn.Close()
	q := `SELECT TABLE_NAME,COLUMN_NAME,data_type,udt_name FROM information_schema.COLUMNS WHERE table_schema='public' ORDER BY TABLE_NAME,ordinal_position;`
	rows, err := conn.Query(q)
	if err != nil {
		PrintErr(ErrQuery)
		panic(err)
	}

	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.DataType, &cs.UdtName)
		if err != nil {
			PrintErr(err)
		}

		if _, ok := tables[cs.TableName]; !ok {
			tables[cs.TableName] = make([]*ColumnSchema, 0)
		}
		tables[cs.TableName] = append(tables[cs.TableName], &cs)
	}
	if err := rows.Err(); err != nil {
		PrintErr(err)
	}
	return tables
}

func getMySQLSchema(dbAddr, dbName string) map[string][]*ColumnSchema {
	tables := make(map[string][]*ColumnSchema)

	conn, err := sql.Open("mysql",fmt.Sprintf("%s/information_schema",dbAddr))
	if err != nil {
		PrintErr(ErrConnect)
		panic(err)
	}
	defer conn.Close()
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
		"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, " +
		"COLUMN_KEY FROM COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION"
	rows, err := conn.Query(q, dbName)
	if err != nil {
		PrintErr(err)
		log.Fatal(err)
	}

	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale,
			&cs.ColumnType, &cs.ColumnKey)
		if err != nil {
			PrintErr(err)

		}

		if _, ok := tables[cs.TableName]; !ok {
			tables[cs.TableName] = make([]*ColumnSchema, 0)
		}
		tables[cs.TableName] = append(tables[cs.TableName], &cs)
	}
	if err := rows.Err(); err != nil {
		PrintErr(err)
	}
	return tables
}

// getSchema 反射表结构
func getSchema(driverName string, dbaddr, dbName string) map[string][]*ColumnSchema{

	switch driverName {

	case "mysql":
		return getMySQLSchema(dbaddr, dbName)
	default:

		return getPgSchema(dbaddr, dbName)
	}
}
