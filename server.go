package main

import (
	"log"
	"net/http"

	"gihub.com/wesleywcr/cv-manager/config"
	"gihub.com/wesleywcr/cv-manager/db"
	"gihub.com/wesleywcr/cv-manager/generated"
	"gihub.com/wesleywcr/cv-manager/resolvers"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	env := config.GetEnv()
	port := env.Port
	if port == "" {
		port = config.DEFAULT_PORT
	}
	db, error := db.New(env.DBName)
	if error != nil {
		log.Fatal((error))
	}

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"*"}, //any address can access
				AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
				AllowCredentials: true,
			},
		),
	)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{DB: db}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
