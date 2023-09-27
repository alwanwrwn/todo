package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/rs/cors"

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
    rows, err := db.Query("SELECT * FROM todo_list WHERE completed = FALSE")
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

func updateTodoItem(w http.ResponseWriter, r *http.Request) {
    fmt.Println("update todo")

    var item TodoItem

    err := json.NewDecoder(r.Body).Decode(&item)

    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    // Update the to-do item in the database
    _, err = db.Exec("UPDATE todo_list SET task = $1, completed = $2 WHERE id = $3", item.Task, item.Completed, item.ID)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    
    w.WriteHeader(http.StatusOK)
}

func toDoRouter() http.Handler {
    r := mux.NewRouter()

    // Define HTTP routes and handlers
    r.HandleFunc("/get-todo", getTodoItems).Methods("GET")
    
    r.HandleFunc("/add-todo", addTodoItem).Methods("POST")

    r.HandleFunc("/update-todo", updateTodoItem).Methods("POST")

    r.Methods(http.MethodOptions).HandlerFunc(handleOptions)

    return r
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers to allow the specific origin, methods, and headers
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Respond with a 200 OK status for the preflight request
    w.WriteHeader(http.StatusOK)
}

func main() { 
    // Create a new CORS middleware instance
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
        Debug:          true,     
      })

    http.Handle("/", c.Handler(toDoRouter())) 

    // Start the HTTP server
    port := ":8080"
    fmt.Printf("Server is listening on %s...\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
