package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/azarolol/gqlen-forum/config"
	"github.com/azarolol/gqlen-forum/db"
	"github.com/azarolol/gqlen-forum/graph"
	"github.com/go-pg/pg"
)

func main() {
	configPath := flag.String("config", "./config.toml", "path to config file")
	flag.Parse()
	config := config.ParseConfig(*configPath)

	var database db.DB

	if config.IfPg {
		var opts = pg.Options{
			User:     config.PgUser,
			Password: config.PgPassword,
			Database: config.PgDatabase,
		}
		database = db.Connect(opts)
	} else {
		database = db.CreateLocalDB()
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: database}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
