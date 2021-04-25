package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	xerrors "github.com/pkg/errors"
)

/*
Q: 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

use golang;

CREATE TABLE lesson
(
	id INT(11) primary key,
	name VARCHAR(25),
	score FLOAT
);

INSERT INTO lesson
(id, name, score)
VALUE(1, "lihua", 90);

INSERT INTO lesson
(id, name, score)
VALUE(2, "zhangli", 60);
*/

var ErrDbNil = errors.New("db is nil, can't query")

func queryNameById(db *sql.DB, id int) (name string, err error) {
	if db == nil {
		return name, ErrDbNil
	}
	err = db.QueryRow("select name from lesson where id = ?", id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// there were no rows, but otherwise no error occurred
			return name, xerrors.Wrapf(err, "[queryNameById] not find row from db where id = %d", id)
		}
	}
	return name, err
}

func main() {
	fmt.Println("starting lesson 2")
	// Open up our database connection.
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/golang")

	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	// Execute the query
	name, err := queryNameById(db, 3)
	if err != nil {
		fmt.Println(xerrors.WithStack(err))
		return
	}
	fmt.Println("result: query name is ", name)
	fmt.Println("finished lesson 2")
	return
}
