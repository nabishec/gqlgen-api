package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nabishec/graphapi/db"
	"github.com/nabishec/graphapi/graph"
	"github.com/nabishec/graphapi/inmemory"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	//logging
	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Logging file couldn`t be opened:", err)
	}
	log.SetOutput(logFile)
	// defer func() { never happens because server doesn't shut down
	// 	log.Println("Program's ending")
	// 	if err := logFile.Close(); err != nil {
	// 		log.Println("Error closing log file:", err)
	// 	}
	// }()
	log.Println("Program started")
	//running program with selected storage method
	var storage *graph.Resolver

	configOfStorage := choosingStorageMethod()
	if configOfStorage == 1 {
		database, err := createDatabase()
		if err != nil {
			log.Fatal("Failed connection to database:", err)
		}
		defer func() {
			if err := database.CloseDatabase(); err != nil {
				log.Fatal("Failed closing of database:", err)
			} else {
				log.Print("Database conection closed successfully")
			}
		}()
		if err := database.MigrationsUp(); err != nil {
			log.Fatal("Failed running migration", err)
		}

		//checking connection
		if err := database.PingDatabase(); err != nil {
			log.Fatal("!Failed connection to database:", err)
		} else {
			log.Println("Connection to database successful")
		}

		storage = graph.NewResolver(db.NewDatabaseResolver(database.DB))
	} else {
		storage = graph.NewResolver(inmemory.NewMemoryResolver())
	}
	server(storage)
}

func server(storage *graph.Resolver) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: storage}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createDatabase() (*db.Database, error) {
	database := db.NewDatabase()
	dbConnfig := &db.DataSourceName{
		Protocol: "postgres",
		Username: "postgres",
		Password: "secret",
		Host:     "localhost",
		Port:     "5432",
		DBname:   "postgres",
		Options:  "sslmode=disable",
	}
	err := database.ConnectDatabase(dbConnfig)
	if err == nil {
		log.Println("Connectiion  to database is successfull")
	}
	return database, err
}

func choosingStorageMethod() int {

	var in = bufio.NewReader(os.Stdin)

	var configuration int

	fmt.Println("Choose a storage method: 1-database 2-localmemory")
	fmt.Println("Answer 1 or 2:")

	fmt.Fscan(in, &configuration)

	if configuration != 1 && configuration != 2 {
		log.Fatal("Answer is incorrect")
	}
	return configuration
}
