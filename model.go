// model.go

package main

import (
	"database/sql"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *user) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT name, age FROM users WHERE id=$1", u.ID).Scan(&u.Name, &u.Age)
}

func (u *user) updateUser(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", u.Name, u.Age, u.ID)
	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)
	return err
}

func (u *user) createUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users(name, age) VALUES($1, $2) RETURNING id", u.Name, u.Age).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	rows, err := db.Query("SELECT id, name, age FROM users LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
