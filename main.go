package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"github.com/AlkBur/VueGolang/handlers"
	"os"
	"os/signal"
	"time"
	"net/http"
	"context"
	"database/sql"
)

func main() {
	db := initDB("storage.db")
	defer db.Close()
	migrate(db)

	router := gin.Default()

	router.StaticFile("/", "public/index.html")
	router.StaticFile("/index.html", "public/index.html")
	router.Static("/css/", "./public/css")
	router.Static("/js/", "./public/js")
	router.GET("/tasks", handlers.GetTasks(db))
	router.PUT("/tasks", handlers.PutTask(db))
	router.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// Start as a web server
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}

	return db
}

func migrate(db *sql.DB) {
	sql := `
    CREATE TABLE IF NOT EXISTS tasks(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL
    );
    `

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}
