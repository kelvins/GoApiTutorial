
# Building and Testing a REST API in GoLang using Gorilla Mux and MySQL

[![Build Status](https://travis-ci.org/kelvins/GoApiTutorial.svg?branch=master)](https://travis-ci.org/kelvins/GoApiTutorial)
[![Coverage Status](https://coveralls.io/repos/github/kelvins/GoApiTutorial/badge.svg?branch=master)](https://coveralls.io/github/kelvins/GoApiTutorial?branch=master)

Link to the Medium story: https://goo.gl/8qe9Du


In this tutorial we will learn how to build and test a simple REST API in Go using [Gorilla Mux](https://github.com/gorilla/mux) router and the [MySQL](https://www.mysql.com/) database. We will also create the application following the [test-driven development](https://en.wikipedia.org/wiki/Test-driven_development) (TDD) methodology.

### Goals

* Become familiar with the TDD methodology.

* Become familiar with the Gorilla Mux package.

* Learn how to use MySQL in Go.

### Prerequisites

* You must have a working Go and MySQL environments.

* Basic familiarity with Go and MySQL.

### About the Application

The application is a simple REST API server that will provide endpoints to allow accessing and manipulating ‘users’.

### API Specification

* Create a new user in response to a valid POST request at /user,

* Update a user in response to a valid PUT request at /user/{id},

* Delete a user in response to a valid DELETE request at /user/{id},

* Fetch a user in response to a valid GET request at /user/{id}, and

* Fetch a list of users in response to a valid GET request at /users.

The {id} will determine which user the request will work with.

### Creating the Database

As our application is simple, we will create only one table called users with the following fields:

* id : is the primary key.

* name : is the name of the user.

* age : is the age of the user.

Let’s use the following statement to create the database and the table.

```sql
    CREATE DATABASE rest_api_example;
    USE rest_api_example;
    CREATE TABLE users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        age INT NOT NULL
    );
```

It is a very simple table but it’s ok for this example.

### Getting Dependencies

Before we start writing our application, we need to get some dependencies that we will use. We need to get the two following packages:

* mux : The Gorilla Mux router.

* mysql : The MySQL driver.

You can easily use go get to get it:

    go get github.com/gorilla/mux
    go get github.com/go-sql-driver/mysql

### Getting Started

First of all, let’s create a file called app.go and add an App structure to hold our application. This structure provides references to the router and the database that we will use on our application. To make it testable let’s also create two methods to initialize and run the application:

```go
    // app.go
    
    package main
    
    import (
        "database/sql"

        "github.com/gorilla/mux"
        _ "github.com/go-sql-driver/mysql"
    )
    
    type App struct {
        Router *mux.Router
        DB     *sql.DB
    }
    
    func (a *App) Initialize(user, password, dbname string) { }
    
    func (a *App) Run(addr string) { }
```

The Initialize method is responsible for create a database connection and wire up the routes, and the Run method will simply start the application.

Note that we have to import both mux and mysql packages here.

Now, let’s create the main.go file which will contain the entry point for the application:

```go
    // main.go
    
    package main
    
    func main() {
        a := App{} 
        // You need to set your Username and Password here
        a.Initialize("DB_USERNAME", "DB_PASSWORD", "rest_api_example")
    
        a.Run(":8080")
    }
```

Note that on this step you need to set the username and password.

Now, let’s create a file called model.go which is used to define our user structure and provide some useful functions to deal with database operations.

```go
    // model.go

    package main

    import (
        "database/sql"
        "errors"
    )

    type user struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Age   int    `json:"age"`
    }

    func (u *user) getUser(db *sql.DB) error {
        return errors.New("Not implemented")
    }

    func (u *user) updateUser(db *sql.DB) error {
        return errors.New("Not implemented")
    }

    func (u *user) deleteUser(db *sql.DB) error {
        return errors.New("Not implemented")
    }

    func (u *user) createUser(db *sql.DB) error {
        return errors.New("Not implemented")
    }

    func getUsers(db *sql.DB, start, count int) ([]user, error) {
        return nil, errors.New("Not implemented")
    }
```

At this point we should have a file structure like that:

    ┌── app.go
    ├── main.go
    └── model.go

Now it’s time to write some tests for our API.

### Writing Tests

As we are following the test-driven development (TDD) methodology, we need to write the test even before we write the functions itself.

As we will run the tests using a database, we need to make sure the database is set up before running the tests and cleaned up after the tests. So let’s create the main_test.go file. In the main_test.go file let’s create the TestMain function which is executed before all tests and will do these stuff for us.

```go
    // main_test.go

    package main

    import (
        "os"
        "log"
        "testing"
    )

    var a App

    func TestMain(m *testing.M) {
        a = App{}
        a.Initialize("DB_USERNAME", "DB_PASSWORD", "rest_api_example")

        ensureTableExists()

        code := m.Run()

        clearTable()

        os.Exit(code)
    }

    func ensureTableExists() {
        if _, err := a.DB.Exec(tableCreationQuery); err != nil {
            log.Fatal(err)
        }
    }

    func clearTable() {
        a.DB.Exec("DELETE FROM users")
        a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
    }

    const tableCreationQuery = `
    CREATE TABLE IF NOT EXISTS users
    (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        age INT NOT NULL
    )`
```

Note that the global variable a represents the application that we want to test.

We use the ensureTableExists function that the table we need for testing is available. The tableCreationQuery is a constant which is a query used to create the database table.

After run the tests we need to call the clearTable function to clean the database up.

In order to run the tests we need to implement the Initialize function in the app.go file, to create a database connection and initialize the router. Now the Initialize function should look like this:

```go
    // app.go
    
    func (a *App) Initialize(user, password, dbname string) {
        connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

        var err error
        a.DB, err = sql.Open("mysql", connectionString)
        if err != nil {
            log.Fatal(err)
        }

        a.Router = mux.NewRouter()
    }
```

At this point even if we don’t have any tests we should be able to run go test without finding any runtime errors. Let’s try it out:

    go test -v

Executing this command should result something like this:

    testing: warning: no tests to run
    PASS
    ok      _/home/user/app 0.051s

### Writing API Tests

Let’s start testing the response of the /users endpoint with an empty table.

```go
    // main_test.go

    func TestEmptyTable(t *testing.T) {
        clearTable()

        req, _ := http.NewRequest("GET", "/users", nil)
        response := executeRequest(req)

        checkResponseCode(t, http.StatusOK, response.Code)

        if body := response.Body.String(); body != "[]" {
            t.Errorf("Expected an empty array. Got %s", body)
        }
    }
```

This test will delete all records in the users table and send a GET request to the /users endpoint.

We use the executeRequest function to execute the request, and checkResponseCode function to test that the HTTP response code is what we expect, and finally, we check the body of the response and check if it is what we expect.

So, let’s implement the executeRequest and checkResponseCode functions.

```go
    // main_test.go
    
    func executeRequest(req *http.Request) *httptest.ResponseRecorder {
        rr := httptest.NewRecorder()
        a.Router.ServeHTTP(rr, req)
    
        return rr
    }

    func checkResponseCode(t *testing.T, expected, actual int) {
        if expected != actual {
            t.Errorf("Expected response code %d. Got %d\n", expected, actual)
        }
    }
```

Make sure you have imported the "net/http" and "net/http/httptest" packages and run the tests again. If everything goes well you should get something like this:

    === RUN   TestEmptyTable
    --- FAIL: TestEmptyTable (0.02s)
        main_test.go:71: Expected response code 200. Got 404
        main_test.go:58: Expected an empty array. Got 404 page not found
    FAIL
    exit status 1
    FAIL    _/home/user/app 0.055s

As expected, the tests will fail because we don’t have implemented anything yet, so let’s continue implementing other tests before we implement the functions for the application itself.

Let’s implement a test that tries to fetch a nonexistent user.

```go
    // main_test.go
    
    func TestGetNonExistentUser(t *testing.T) {
        clearTable()

        req, _ := http.NewRequest("GET", "/user/45", nil)
        response := executeRequest(req)

        checkResponseCode(t, http.StatusNotFound, response.Code)

        var m map[string]string
        json.Unmarshal(response.Body.Bytes(), &m)
        if m["error"] != "User not found" {
            t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
        }
    }
```

This test basically tests two things: the status code which should be 404 and if the response contains the expected error message.

Note that in this step we need to import the "encoding/json" package to use the json.Unmarshal function.

Now, let’s implement a test to create a user.

```go
    // main_test.go

    func TestCreateUser(t *testing.T) {
        clearTable()

        payload := []byte(`{"name":"test user","age":30}`)

        req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
        response := executeRequest(req)

        checkResponseCode(t, http.StatusCreated, response.Code)

        var m map[string]interface{}
        json.Unmarshal(response.Body.Bytes(), &m)

        if m["name"] != "test user" {
            t.Errorf("Expected user name to be 'test user'. Got '%v'", m["name"])
        }

        if m["age"] != 30.0 {
            t.Errorf("Expected user age to be '30'. Got '%v'", m["age"])
        }

        // the id is compared to 1.0 because JSON unmarshaling converts numbers to
        // floats, when the target is a map[string]interface{}
        if m["id"] != 1.0 {
            t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
        }
    }
```

In this test, we manually add a new user to the database and, by accessing the correspondent endpoint, we check if the status code is 201 (the resource was created) and if the JSON response contains the correct information that was added.

Note that in this step we need to import the "bytes" package to use the bytes.NewBuffer function.

Now, let’s implement a test to fetch an existing user.

```go
    // main_test.go
    
    func TestGetUser(t *testing.T) {
        clearTable()
        addUsers(1)

        req, _ := http.NewRequest("GET", "/user/1", nil)
        response := executeRequest(req)

        checkResponseCode(t, http.StatusOK, response.Code)
    }
```

This test basically add a new user to the database and check if the correct endpoint results in an HTTP response with status code 200 (success).

In this test above we use the addUsers function which is used to add a new user to the database for the tests. So, let’s implement this function:

```go
    // main_test.go
    
    func addUsers(count int) {
        if count < 1 {
            count = 1
        }

        for i := 0; i < count; i++ {
            statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)", ("User " + strconv.Itoa(i+1)), ((i+1) * 10))
            a.DB.Exec(statement)
        }
    }
```

Note that in this step we need to import the "strconv" package to use the strconv.Itoa function to convert an integer to a string.

Now, let’s test the update option:

```go
    // main_test.go

    func TestUpdateUser(t *testing.T) {
        clearTable()
        addUsers(1)

        req, _ := http.NewRequest("GET", "/user/1", nil)
        response := executeRequest(req)
        var originalUser map[string]interface{}
        json.Unmarshal(response.Body.Bytes(), &originalUser)

        payload := []byte(`{"name":"test user - updated name","age":21}`)

        req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
        response = executeRequest(req)

        checkResponseCode(t, http.StatusOK, response.Code)

        var m map[string]interface{}
        json.Unmarshal(response.Body.Bytes(), &m)

        if m["id"] != originalUser["id"] {
            t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
        }

        if m["name"] == originalUser["name"] {
            t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser["name"], m["name"], m["name"])
        }

        if m["age"] == originalUser["age"] {
            t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", originalUser["age"], m["age"], m["age"])
        }
    }
```

In the above test, we basically add a new user to the database and then we use the correct endpoint to update it.

It tests if the status code is 200 indicating success and if the JSON response contains the updated details about the user.

And the last test, for now, will try to delete a user.

```go
    // main_test.go
    
    func TestDeleteUser(t *testing.T) {
        clearTable()
        addUsers(1)

        req, _ := http.NewRequest("GET", "/user/1", nil)
        response := executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code)

        req, _ = http.NewRequest("DELETE", "/user/1", nil)
        response = executeRequest(req)

        checkResponseCode(t, http.StatusOK, response.Code)

        req, _ = http.NewRequest("GET", "/user/1", nil)
        response = executeRequest(req)
        checkResponseCode(t, http.StatusNotFound, response.Code)
    }
```

In this test we basically create a new user and test if it exists in the database, then we user the correct endpoint to delete the user and checks if it was properly deleted.

At this point we should be able to run go test -v in your project directory.

All tests should fail but it’s ok because we did not implement the application functions yet. So let’s implement it to make these tests pass.

### Creating the Application Functionalities

Let’s begin implementing the methods in the model.go file. These methods are responsible for executing the database statements and it can be implemented as follows:

```go
    // model.go
    
    func (u *user) getUser(db *sql.DB) error {
        statement := fmt.Sprintf("SELECT name, age FROM users WHERE id=%d", u.ID)
        return db.QueryRow(statement).Scan(&u.Name, &u.Age)
    }

    func (u *user) updateUser(db *sql.DB) error {
        statement := fmt.Sprintf("UPDATE users SET name='%s', age=%d WHERE id=%d", u.Name, u.Age, u.ID)
        _, err := db.Exec(statement)
        return err
    }

    func (u *user) deleteUser(db *sql.DB) error {
        statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
        _, err := db.Exec(statement)
        return err
    }

    func (u *user) createUser(db *sql.DB) error {
        statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)", u.Name, u.Age)
        _, err := db.Exec(statement)

        if err != nil {
            return err
        }

        err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

        if err != nil {
            return err
        }

        return nil
    }

    func getUsers(db *sql.DB, start, count int) ([]user, error) {
        statement := fmt.Sprintf("SELECT id, name, age FROM users LIMIT %d OFFSET %d", count, start)
        rows, err := db.Query(statement)

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
```

The getUsers function fetches records from the users table and limits the number of records based on the count value passed by parameter. The start parameter determines how many records are skipped at the beginning.

At this point, we need to remove the errors package and import the fmt package.

The model is done, now we need to implement the App functions, including the **routes** and **route handlers**.

Let’s start creating the getUser function to fetch a single user.

```go
    // app.go
    
    func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid user ID")
            return
        }

        u := user{ID: id}
        if err := u.getUser(a.DB); err != nil {
            switch err {
            case sql.ErrNoRows:
                respondWithError(w, http.StatusNotFound, "User not found")
            default:
                respondWithError(w, http.StatusInternalServerError, err.Error())
            }
            return
        }

        respondWithJSON(w, http.StatusOK, u)
    }
```

This handler basically retrieves the id of the user from the requested URL and uses the getUser function, from the model, to fetch the user details.

If the user is not found it will respond with the status code 404. This function uses the respondWithError and respondWithJSON functions to process errors and normal responses. These functions are implemented as follows:

```go
    // app.go
    
    func respondWithError(w http.ResponseWriter, code int, message string) {
        respondWithJSON(w, code, map[string]string{"error": message})
    }

    func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
        response, _ := json.Marshal(payload)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        w.Write(response)
    }
```

The rest of the handlers can be implemented in a similar manner:

```go
    // app.go
    
    func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
        count, _ := strconv.Atoi(r.FormValue("count"))
        start, _ := strconv.Atoi(r.FormValue("start"))

        if count > 10 || count < 1 {
            count = 10
        }
        if start < 0 {
            start = 0
        }

        users, err := getUsers(a.DB, start, count)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }

        respondWithJSON(w, http.StatusOK, users)
    }
```

This handler uses the count and start parameters from the querystring to fetch count number of users, starting at position start in the database. By default, start is set to 0 and count is set to 10. If these parameters aren’t provided, this handler will respond with the first 10 users.

Let’s implement the handler to create a user.

```go
    // app.go
    
    func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
        var u user
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&u); err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid request payload")
            return
        }
        defer r.Body.Close()

        if err := u.createUser(a.DB); err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }

        respondWithJSON(w, http.StatusCreated, u)
    }
```

This handler assumes that the request body is a JSON object containing the details of the user to be created. It extracts that object into a user and then uses the createUser function.

The handler to update a user:

```go
    // app.go
    
    func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid user ID")
            return
        }

        var u user
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&u); err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
            return
        }
        defer r.Body.Close()
        u.ID = id

        if err := u.updateUser(a.DB); err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }

        respondWithJSON(w, http.StatusOK, u)
    }
```

This handler extracts the user details from the request body and the id from the URL, and uses the id and the body to update the user.

And the last handler that we will implement is used to delete a user.

```go
    // app.go
    
    func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid User ID")
            return
        }

        u := user{ID: id}
        if err := u.deleteUser(a.DB); err != nil {
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }

        respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
    }
```

This handler extracts the id from the URL and uses it to delete the corresponding user.

Now that we have all handlers implemented we must define the routes which will use them.

```go
    // app.go
    
    func (a *App) initializeRoutes() {
        a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
        a.Router.HandleFunc("/user", a.createUser).Methods("POST")
        a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
        a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
        a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
    }
```

The routes are defined based on the API specification defined earlier. The {id:[0-9]+} part of the path indicates that Gorilla Mux should treat process a URL only if the id is a number. For all matching requests, Gorilla Mux then stores the the actual numeric value in the id variable.

Now we just need to implement the Run function and call initializeRoutes from the Initialize method.

```go
    // app.go
    
    func (a *App) Initialize(user, password, dbname string) {
        connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

        var err error
        a.DB, err = sql.Open("mysql", connectionString)
        if err != nil {
            log.Fatal(err)
        }

        a.Router = mux.NewRouter()
        a.initializeRoutes()
    }

    func (a *App) Run(addr string) {
        log.Fatal(http.ListenAndServe(addr, a.Router))
    }
```

Remember to import all packages needed.

```go
    // app.go
    
    import (
        "database/sql"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "strconv"

        _ "github.com/go-sql-driver/mysql"
        "github.com/gorilla/mux"
    )
```

The final version of the app.go file should look like this: [https://github.com/kelvins/GoApiTutorial/blob/master/app.go](https://github.com/kelvins/GoApiTutorial/blob/master/app.go)

Now if we run the tests again:

    go test -v

We should get the following results:

    === RUN   TestEmptyTable
    --- PASS: TestEmptyTable (0.02s)
    === RUN   TestGetNonExistentUser
    --- PASS: TestGetNonExistentUser (0.02s)
    === RUN   TestCreateUser
    --- PASS: TestCreateUser (0.01s)
    === RUN   TestGetUser
    --- PASS: TestGetUser (0.01s)
    === RUN   TestUpdateUser
    --- PASS: TestUpdateUser (0.01s)
    === RUN   TestDeleteUser
    --- PASS: TestDeleteUser (0.01s)
    PASS
    ok   github.com/kelvins/goapi 0.124s

The complete code can be found on **Github** at the following link: [https://github.com/kelvins/GoApiTutorial](https://github.com/kelvins/GoApiTutorial)

### Travis CI and Coveralls

If you are familiar with [Travis CI](https://travis-ci.org/) and [Coveralls](https://coveralls.io/) you can use the following settings for the build environment on the .travis.yml file:

    language: go

    go:
      - 1.6
      - 1.8
      - tip

    services:
      - mysql

    before_install:
      - mysql -e 'CREATE DATABASE IF NOT EXISTS rest_api_example;'

    install:
      - go get golang.org/x/tools/cmd/cover
      - go get github.com/mattn/goveralls
      - go get github.com/gorilla/mux
      - go get github.com/go-sql-driver/mysql

    script:
      - go test -covermode=count -coverprofile=coverage.out
      - $HOME/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci

**Note** that for the **Travis CI** run the tests properly using the **MySQL** database you need to set the **username** as root and leave the **password** empty.

If you are not familiar with it, I suggest starting reading the **Getting Started **section of both [Travis CI](https://docs.travis-ci.com/user/getting-started/) and [Coveralls](https://coveralls.zendesk.com/hc/en-us). These tools are well documented and quite easy to understand and use.

If you want to manually test the API by manually sending requests I suggest to use the **Insomnia** application. It is a cross-platform REST API client that is very easy to use. It can be found here: [https://github.com/getinsomnia/insomnia](https://github.com/getinsomnia/insomnia)

### References

Almost all this tutorial was created (and some codes copied) based on the tutorial written by Kulshekhar Kabra which can be found at the following link:
[**Building and Testing a REST API in Go with Gorilla Mux and PostgreSQL**
*Learn how to build simple and well-tested REST APIs backed by PostgreSQL in Go, using Gorilla Mux - a highly stable and…*semaphoreci.com](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)
