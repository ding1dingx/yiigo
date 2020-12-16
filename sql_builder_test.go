package yiigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToQuery(t *testing.T) {
	query, binds := builder.Wrap(
		Table("user"),
		Where("id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user WHERE id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		Where("name = ? AND age > ?", "shenghui0779", 20),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user WHERE name = ? AND age > ?", query)
	assert.Equal(t, []interface{}{"shenghui0779", 20}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		WhereIn("age IN (?)", []int{20, 30}),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user WHERE age IN (?, ?)", query)
	assert.Equal(t, []interface{}{20, 30}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		Select("id", "name", "age"),
		Where("id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT id, name, age FROM user WHERE id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		Distinct("name"),
		Where("id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT DISTINCT name FROM user WHERE id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		Join("address", "user.id = address.user_id"),
		Where("user.id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user INNER JOIN address ON user.id = address.user_id WHERE user.id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		LeftJoin("address", "user.id = address.user_id"),
		Where("user.id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user LEFT JOIN address ON user.id = address.user_id WHERE user.id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		RightJoin("address", "user.id = address.user_id"),
		Where("user.id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user RIGHT JOIN address ON user.id = address.user_id WHERE user.id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		FullJoin("address", "user.id = address.user_id"),
		Where("user.id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user FULL JOIN address ON user.id = address.user_id WHERE user.id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("address"),
		Select("user_id", "COUNT(*) AS total"),
		GroupBy("user_id"),
		Having("user_id = ?", 1),
	).ToQuery()

	assert.Equal(t, "SELECT user_id, COUNT(*) AS total FROM address GROUP BY user_id HAVING user_id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		Where("age > ?", 20),
		OrderBy("id DESC"),
		Offset(5),
		Limit(10),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user WHERE age > ? ORDER BY id DESC OFFSET ? LIMIT ?", query)
	assert.Equal(t, []interface{}{20, 5, 10}, binds)

	query, binds = builder.Wrap(
		Table("user00"),
		Where("id = ?", 1),
		Union(builder.Wrap(Table("user01"), Where("id = ?", 2))),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user00 WHERE id = ? UNION SELECT * FROM user01 WHERE id = ?", query)
	assert.Equal(t, []interface{}{1, 2}, binds)

	query, binds = builder.Wrap(
		Table("user_0"),
		Where("id = ?", 1),
		UnionAll(builder.Wrap(Table("user_1"), Where("id = ?", 2))),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user_0 WHERE id = ? UNION ALL SELECT * FROM user_1 WHERE id = ?", query)
	assert.Equal(t, []interface{}{1, 2}, binds)

	query, binds = builder.Wrap(
		Table("user_0"),
		WhereIn("id IN (?)", []int{1, 2}),
		Union(builder.Wrap(Table("user_1"), Where("id = ?", 3))),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user_0 WHERE id IN (?, ?) UNION SELECT * FROM user_1 WHERE id = ?", query)
	assert.Equal(t, []interface{}{1, 2, 3}, binds)

	query, binds = builder.Wrap(
		Table("user_0"),
		Where("id = ?", 1),
		Union(builder.Wrap(Table("user_1"), WhereIn("id IN (?)", []int{2, 3}))),
	).ToQuery()

	assert.Equal(t, "SELECT * FROM user_0 WHERE id = ? UNION SELECT * FROM user_1 WHERE id IN (?, ?)", query)
	assert.Equal(t, []interface{}{1, 2, 3}, binds)
}

func TestToInsert(t *testing.T) {
	type User struct {
		ID     int    `db:"-"`
		Name   string `db:"name"`
		Gender string `db:"gender"`
		Age    int    `db:"age"`
	}

	query, binds := builder.Wrap(Table("user")).ToInsert(&User{
		Name:   "shenghui0779",
		Gender: "M",
		Age:    29,
	})

	assert.Equal(t, "INSERT INTO user (name, gender, age) VALUES (?, ?, ?)", query)
	assert.Equal(t, []interface{}{"shenghui0779", "M", 29}, binds)
}

func TestToBatchInsert(t *testing.T) {
	type User struct {
		ID     int    `db:"-"`
		Name   string `db:"name"`
		Gender string `db:"gender"`
		Age    int    `db:"age"`
	}

	query, binds := builder.Wrap(Table("user")).ToBatchInsert([]*User{
		{
			Name:   "shenghui0779",
			Gender: "M",
			Age:    29,
		},
		{
			Name:   "test",
			Gender: "W",
			Age:    20,
		},
	})

	assert.Equal(t, "INSERT INTO user (name, gender, age) VALUES (?, ?, ?), (?, ?, ?)", query)
	assert.Equal(t, []interface{}{"shenghui0779", "M", 29, "test", "W", 20}, binds)
}

func TestToUpdate(t *testing.T) {
	type User struct {
		Name   string `db:"name"`
		Gender string `db:"gender"`
		Age    int    `db:"age"`
	}

	query, binds := builder.Wrap(
		Table("user"),
		Where("id = ?", 1),
	).ToUpdate(&User{
		Name:   "shenghui0779",
		Gender: "M",
		Age:    29,
	})

	assert.Equal(t, "UPDATE user SET name = ?, gender = ?, age = ? WHERE id = ?", query)
	assert.Equal(t, []interface{}{"shenghui0779", "M", 29, 1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		WhereIn("id IN (?)", []int{1, 2}),
	).ToUpdate(&User{
		Name:   "shenghui0779",
		Gender: "M",
		Age:    29,
	})

	assert.Equal(t, "UPDATE user SET name = ?, gender = ?, age = ? WHERE id IN (?, ?)", query)
	assert.Equal(t, []interface{}{"shenghui0779", "M", 29, 1, 2}, binds)
}

func TestToDelete(t *testing.T) {
	query, binds := builder.Wrap(
		Table("user"),
		Where("id = ?", 1),
	).ToDelete()

	assert.Equal(t, "DELETE FROM user WHERE id = ?", query)
	assert.Equal(t, []interface{}{1}, binds)

	query, binds = builder.Wrap(
		Table("user"),
		WhereIn("id IN (?)", []int{1, 2}),
	).ToDelete()

	assert.Equal(t, "DELETE FROM user WHERE id IN (?, ?)", query)
	assert.Equal(t, []interface{}{1, 2}, binds)
}
