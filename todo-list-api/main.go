package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", "postgres://username:password@localhost/todoapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")
}

func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(todo.Title, todo.Status).Scan(&todo.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, todo)
}

func getTodos(c echo.Context) error {
	rows, err := db.Query("SELECT id, title, status FROM todos")
	if err != nil {
		return err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(http.StatusOK, todos)
}

func getTodoByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	row := db.QueryRow("SELECT id, title, status FROM todos WHERE id = $1", id)
	todo := &Todo{}
	err = row.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Todo not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, todo)
}

func updateTodoByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}

	result, err := db.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3", todo.Title, todo.Status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return c.String(http.StatusNotFound, "Todo not found")
	}

	return c.String(http.StatusOK, "Todo updated successfully")
}

func deleteTodoByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	result, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return c.String(http.StatusNotFound, "Todo not found")
	}

	return c.String(http.StatusOK, "Todo deleted successfully")
}

func main() {
	initDB()

	e := echo.New()

	// Routes
	e.POST("/todos", createTodo)
	e.GET("/todos", getTodos)
	e.GET("/todos/:id", getTodoByID)
	e.PUT("/todos/:id", updateTodoByID)
	e.DELETE("/todos/:id", deleteTodoByID)

	log.Println("Server is running on port 8080...")
	log.Fatal(e.Start(":8080"))
}