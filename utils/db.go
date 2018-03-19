package utils

import (
	"fmt"
	"reflect"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Sequence model for mongo _id auto_increment
type Sequence struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

// SQLExpr SQL expression
type SQLExpr struct {
	Expr string
	Args []interface{}
}

// InsertSQL returns insert sql and binds
// data expect struct, []struct, yiigo.X, []yiigo.X
func InsertSQL(table string, data interface{}) (string, []interface{}) {
	v := reflect.Indirect(reflect.ValueOf(data))

	sql := ""
	binds := []interface{}{}

	switch v.Kind() {
	case reflect.Map:
		if x, ok := data.(X); ok {
			sql, binds = singleInsertWithMap(sql, x)
		}
	case reflect.Struct:
		sql, binds = singleInsertWithStruct(table, v)
	case reflect.Slice:
		if count := v.Len(); count > 0 {
			elemKind := v.Type().Elem().Kind()

			if elemKind == reflect.Map {
				if x, ok := data.([]X); ok {
					sql, binds = batchInsertWithMap(table, x, count)
				}

				break
			}

			if elemKind == reflect.Struct {
				sql, binds = batchInsertWithStruct(table, v, count)

				break
			}
		}
	}

	return sql, binds
}

// UpdateSQL returns update sql and binds
// data expect struct, yiigo.X
func UpdateSQL(query string, data interface{}, args ...interface{}) (string, []interface{}) {
	v := reflect.Indirect(reflect.ValueOf(data))

	sql := ""
	binds := []interface{}{}

	switch v.Kind() {
	case reflect.Map:
		if x, ok := data.(X); ok {
			sql, binds = updateWithMap(query, x, args...)
		}
	case reflect.Struct:
		sql, binds = updateWithStruct(query, v, args...)
	}

	return sql, binds
}

// Expr returns expression, eg: yiigo.Expr("price * ? + ?", 2, 100)
func Expr(expr string, args ...interface{}) *SQLExpr {
	return &SQLExpr{Expr: expr, Args: args}
}

// SeqID get mongo auto_increment _id
func SeqID(session *mgo.Session, db string, collection string, seqs ...int) (int, error) {
	if len(seqs) == 0 {
		seqs = append(seqs, 1)
	}

	condition := bson.M{"_id": collection}

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": seqs[0]}},
		Upsert:    true,
		ReturnNew: true,
	}

	sequence := Sequence{}

	_, err := session.DB(db).C("sequence").Find(condition).Apply(change, &sequence)

	if err != nil {
		return 0, err
	}

	return sequence.Seq, nil
}

func singleInsertWithMap(table string, data X) (string, []interface{}) {
	fieldNum := len(data)

	columns := make([]string, 0, fieldNum)
	placeholders := make([]string, 0, fieldNum)
	binds := make([]interface{}, 0, fieldNum)

	for k, v := range data {
		columns = append(columns, fmt.Sprintf("`%s`", k))
		placeholders = append(placeholders, "?")
		binds = append(binds, v)
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return sql, binds
}

func singleInsertWithStruct(table string, v reflect.Value) (string, []interface{}) {
	fieldNum := v.NumField()

	columns := make([]string, 0, fieldNum)
	placeholders := make([]string, 0, fieldNum)
	binds := make([]interface{}, 0, fieldNum)

	t := v.Type()

	for i := 0; i < fieldNum; i++ {
		column := t.Field(i).Tag.Get("db")

		if column == "" {
			column = t.Field(i).Name
		}

		columns = append(columns, fmt.Sprintf("`%s`", column))
		placeholders = append(placeholders, "?")
		binds = append(binds, v.Field(i).Interface())
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return sql, binds
}

func batchInsertWithMap(table string, data []X, count int) (string, []interface{}) {
	fieldNum := len(data[0])

	fields := make([]string, 0, fieldNum)
	columns := make([]string, 0, fieldNum)
	placeholders := make([]string, 0, fieldNum)
	binds := make([]interface{}, 0, fieldNum*count)

	for k := range data[0] {
		fields = append(fields, k)
		columns = append(columns, fmt.Sprintf("`%s`", k))
	}

	fmt.Println(columns)

	for _, x := range data {
		phrs := make([]string, 0, fieldNum)

		for _, v := range fields {
			phrs = append(phrs, "?")
			binds = append(binds, x[v])
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(phrs, ", ")))
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s", table, strings.Join(columns, ", "), strings.Join(placeholders, ","))

	return sql, binds
}

func batchInsertWithStruct(table string, v reflect.Value, count int) (string, []interface{}) {
	first := reflect.Indirect(v.Index(0))

	if first.Kind() != reflect.Struct {
		panic("the data must be a slice to struct")
	}

	fieldNum := first.NumField()

	columns := make([]string, 0, fieldNum)
	placeholders := make([]string, 0, fieldNum)
	binds := make([]interface{}, 0, fieldNum*count)

	t := first.Type()

	for i := 0; i < fieldNum; i++ {
		column := t.Field(i).Tag.Get("db")

		if column == "" {
			column = t.Field(i).Name
		}

		columns = append(columns, fmt.Sprintf("`%s`", column))
	}

	for i := 0; i < count; i++ {
		phrs := make([]string, 0, fieldNum)

		for j := 0; j < fieldNum; j++ {
			phrs = append(phrs, "?")
			binds = append(binds, reflect.Indirect(v.Index(i)).Field(j).Interface())
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(phrs, ", ")))
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s", table, strings.Join(columns, ", "), strings.Join(placeholders, ","))

	return sql, binds
}

func updateWithMap(query string, data X, args ...interface{}) (string, []interface{}) {
	sets := []string{}
	binds := []interface{}{}

	for k, v := range data {
		if e, ok := v.(*SQLExpr); ok {
			sets = append(sets, fmt.Sprintf("`%s` = %s", k, e.Expr))
			binds = append(binds, e.Args...)
		} else {
			sets = append(sets, fmt.Sprintf("`%s` = ?", k))
			binds = append(binds, v)
		}
	}

	sql := strings.Replace(query, "?", strings.Join(sets, ", "), 1)
	binds = append(binds, args...)

	return sql, binds
}

func updateWithStruct(query string, v reflect.Value, args ...interface{}) (string, []interface{}) {
	fieldNum := v.NumField()

	sets := make([]string, 0, fieldNum)
	binds := make([]interface{}, 0, fieldNum+len(args))

	t := v.Type()

	for i := 0; i < fieldNum; i++ {
		column := t.Field(i).Tag.Get("db")

		if column == "" {
			column = t.Field(i).Name
		}

		field := v.Field(i).Interface()

		if e, ok := field.(*SQLExpr); ok {
			sets = append(sets, fmt.Sprintf("`%s` = %s", column, e.Expr))
			binds = append(binds, e.Args...)
		} else {
			sets = append(sets, fmt.Sprintf("`%s` = ?", column))
			binds = append(binds, field)
		}
	}

	sql := strings.Replace(query, "?", strings.Join(sets, ", "), 1)
	binds = append(binds, args...)

	return sql, binds
}
