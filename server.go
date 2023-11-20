package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kilianp07/graphql/graph"
)

const defaultPort = "8080"

var db *sql.DB

func initDB() {
	dataSourceName := "username:password@tcp(localhost:3306)/dbname"
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = createTables()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")
}
func createTables() error {
	// Créer la table Game
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Game (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			publicationDate INT,
			platform VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Créer la table Editor
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Editor (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Créer la table Studio
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Studio (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err

	}

	// Créer la table GameEditor
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS GameEditor (
			gameID INT,
			editorID INT,
			PRIMARY KEY (gameID, editorID),
			FOREIGN KEY (gameID) REFERENCES Game(id),
			FOREIGN KEY (editorID) REFERENCES Editor(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err

	}

	// Créer la table GameStudio
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS GameStudio (
			gameID INT,
			studioID INT,
			PRIMARY KEY (gameID, studioID),
			FOREIGN KEY (gameID) REFERENCES Game(id),
			FOREIGN KEY (studioID) REFERENCES Studio(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
		return err

	}

	log.Println("Tables created successfully")
	return err
}

func main() {
	initDB()
	defer db.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	schema := graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB: db,
		},
	})

	srv := handler.NewDefaultServer(schema)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
