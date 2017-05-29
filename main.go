// main.go

package main

func main() {
    a := App{}
    a.Initialize("DB_USERNAME", "DB_PASSWORD", "rest_api_example")

    a.Run(":8080")
}
