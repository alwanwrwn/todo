package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"

    _ "github.com/lib/pq"
)

var db *sql.DB

func init() {
    // Replace with your PostgreSQL connection string
    connStr := "user=alwanwirawan dbname=todo_list_db sslmode=disable"
    
    // Open a connection to the database
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

type TodoItem struct {
    ID          int     `json:"id"`
    Task        string  `json:"task"`
    Completed   bool    `json:"completed"`
}

func getTodoItems(w http.ResponseWriter, r *http.Request) {
    fmt.Println("get todo")

    // Query the database to retrieve all to-do items
    rows, err := db.Query("SELECT * FROM todo_list")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Create a slice to store the retrieved to-do items
    var todoItems []TodoItem

    // Iterate through the rows and scan them into the struct
    for rows.Next() {
        var item TodoItem
        if err := rows.Scan(&item.ID, &item.Task, &item.Completed); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        todoItems = append(todoItems, item)
    }

    // Convert the slice of to-do items to JSON
    jsonItems, err := json.Marshal(todoItems)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the response content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Write the JSON response
    w.WriteHeader(http.StatusOK)
    w.Write(jsonItems)
}


func addTodoItem(w http.ResponseWriter, r *http.Request) {
    fmt.Println("add todo")

    var item TodoItem

    err := json.NewDecoder(r.Body).Decode(&item)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Insert the new to-do item into the database
    _, err = db.Exec("INSERT INTO todo_list (task, completed) VALUES ($1, $2)", item.Task, item.Completed)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func main() { 
    r := mux.NewRouter()

    // Define HTTP routes and handlers
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        // Implement logic to retrieve and display to-do list items from the database
        
        fmt.Println("Daftar todo x")
        // db.Exec("INSERT INTO todo_list (task, completed) VALUES ('test input', false)")
        //testPrintge()
    })

    r.HandleFunc("/get-todo", getTodoItems).Methods("GET")

    r.HandleFunc("/add-todo", addTodoItem).Methods("POST")

    // Start the HTTP server
    port := ":8080"
    fmt.Printf("Server is listening on %s...\n", port)
    if err := http.ListenAndServe(port, r); err != nil {
        log.Fatal(err)
    }
}

