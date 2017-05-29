// app.go

package main

import (
    "fmt"
    "log"
    "database/sql"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
    connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) { }
