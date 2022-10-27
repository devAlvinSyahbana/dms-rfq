package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/devAlvinSyahbana/golang-rfq/graph"
	"github.com/devAlvinSyahbana/golang-rfq/graph/generated"
	middlewares "github.com/devAlvinSyahbana/golang-rfq/middleware"
	"github.com/devAlvinSyahbana/golang-rfq/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const defaultPort = "4080"

type authString string

const (
	host     = "localhost"
	portDB   = "5432"
	user     = "postgres"
	password = "@231020"
	dbname   = "dms_rfq_new"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, portDB, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.Use(middlewares.AuthMiddleware)
	router.HandleFunc("/download/{id}", service.GeneratePDFMux(db))
	schema := generated.Config{Resolvers: &graph.Resolver{DB: db}}
	schema.Directives.RequireLogin = func(ctx context.Context, obj interface{}, next graphql.Resolver, hastoken bool) (interface{}, error) {
		if hastoken {
			token := middlewares.CtxValue(ctx)
			if token == nil {
				return nil, fmt.Errorf("Please provide token")
			}
			return next(ctx)
		}
		return next(ctx)
	}

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5173"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		Debug:            true,
		AllowOriginFunc:  func(origin string) bool { return true },
	}).Handler)

	gql := generated.NewExecutableSchema(schema)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", middlewares.AuthMiddleware(handler.GraphQL(gql)))

	log.Printf("connect to http://127.0.0.1:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":4080", router))
}
