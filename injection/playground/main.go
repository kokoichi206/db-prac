package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5"
	// file not found
	// #include <libpq-fe.h>
	// _ "github.com/jbarham/gopgsqldriver"
)

type database struct {
	db *sql.DB
	// db *pgx.Conn
}

type user struct {
	Name string
	Age  int
}

// playground-postgres-1  | 2023-10-08 17:23:16.709 UTC [61] LOG:  statement:
// playground-postgres-1  |        SELECT name, age FROM users WHERE name = 'John'
// playground-postgres-1  |
const getAllUsersStmt = `
SELECT name, age FROM users WHERE name = '%s'
`

func (d *database) getAllUsers(ctx context.Context, name string) ([]user, error) {

	rows, err := d.db.QueryContext(ctx, fmt.Sprintf(getAllUsersStmt, name))

	if err != nil {
		return nil, fmt.Errorf("failed to select messages: %w", err)
	}

	return parseRows(rows)
}

// playground-postgres-1  | 2023-10-08 17:21:59.919 UTC [58] LOG:  execute <unnamed>:
// playground-postgres-1  |        SELECT name, age FROM users WHERE name = $1
// playground-postgres-1  |
// playground-postgres-1  | 2023-10-08 17:21:59.919 UTC [58] DETAIL:  parameters: $1 = 'John'
const getAllUsersPlaceholderStmt = `
SELECT name, age FROM users WHERE name = $1
`

func (d *database) getAllUsersPlaceholder(ctx context.Context, name string) ([]user, error) {
	rows, err := d.db.QueryContext(ctx, getAllUsersPlaceholderStmt, name)
	if err != nil {
		return nil, fmt.Errorf("failed to select messages: %w", err)
	}

	return parseRows(rows)
}

const orderByParam = `
SELECT name, age FROM users ORDER BY %s
`

func (d *database) getAllUsersOrderBy(ctx context.Context, param string) ([]user, error) {
	query := fmt.Sprintf(orderByParam, pq.QuoteIdentifier(param))
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select messages: %w", err)
	}

	return parseRows(rows)
}

const orderByParamPlaceHolder = `
SELECT name, age FROM users ORDER BY $1
`

func (d *database) getAllUsersOrderByPH(ctx context.Context, param string) ([]user, error) {
	rows, err := d.db.QueryContext(ctx, orderByParamPlaceHolder, param)
	if err != nil {
		return nil, fmt.Errorf("failed to select messages: %w", err)
	}

	return parseRows(rows)
}


func pgxTest(conn *pgx.Conn, param string) ([]user, error) {
	rows, err := conn.Query(context.Background(), orderByParam, param)
	if err != nil {
		return nil, fmt.Errorf("failed to select messages: %w", err)
	}

	return parsePgxRows(rows)
}

func parsePgxRows(rows pgx.Rows) ([]user, error) {
	resp := []user{}

	for rows.Next() {
		var n sql.NullString
		var a sql.NullInt64

		if err := rows.Scan(&n, &a); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		resp = append(resp, user{
			Name: n.String,
			Age:  int(a.Int64),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return resp, nil
}

func parseRows(rows *sql.Rows) ([]user, error) {
	resp := []user{}

	for rows.Next() {
		var n sql.NullString
		var a sql.NullInt64

		if err := rows.Scan(&n, &a); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		resp = append(resp, user{
			Name: n.String,
			Age:  int(a.Int64),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return resp, nil
}

func connect(driver, host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	source := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	sqlDB, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql: %w", err)
	}

	return sqlDB, nil
}

func main() {
	// // ===================== pgx ======================
	// url := "postgres://ubuntu:ubuntu@localhost:4949/postgres"
	// db, err := pgx.Connect(context.Background(), url)
	// if err != nil {
	// 	panic(err)
	// }
	// // u, err := pgxTest(db, "EXTRACTVALUE(0, (SELECT CONCAT(0x7e, (SELECT (ELT(0x7e, 1) LIMIT 1)), 0x7e)))) AND '1'='1")
	// // u, err := pgxTest(db, "EXTRACTVALUE(0, (SELECT CONCAT('$', name, ':', password) FROM users LIMIT 1))")
	// u, err := pgxTest(db, "(SELECT CONCAT('$', name, ':', password) FROM users LIMIT 1)")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("u: %v\n", u)

	// ===================== lib/pq ======================
	db, err := connect("postgres", "localhost", "4949", "ubuntu", "ubuntu", "postgres", "disable")
	if err != nil {
		panic(err)
	}

	d := &database{db: db}

	// users, err := d.getAllUsersOrderBy(context.Background(), "name; DROP TABLE users; --")
	users, err := d.getAllUsersOrderBy(context.Background(), "age")
	if err != nil {
		panic(err)
	}

	fmt.Println(users)

	users, err = d.getAllUsersOrderByPH(context.Background(), "age")
	if err != nil {
		panic(err)
	}

	fmt.Println(users)

	// uss, err := d.getAllUsersPlaceholder(context.Background(), "John")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(uss)
	// return

	// users, err := d.getAllUsers(context.Background(), "John")
	// // users, err := d.getAllUsers(context.Background(), "John'; DROP TABLE users; --")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(users)
}
