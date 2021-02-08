package main

import (
	"fmt"
	"github.com/brharrelldev/BlockListAPI/database"
	"github.com/dgraph-io/badger/v2"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/brharrelldev/BlockListAPI/auth"
	"github.com/urfave/cli/v2"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/brharrelldev/BlockListAPI/graph"
	"github.com/brharrelldev/BlockListAPI/graph/generated"
	"github.com/gorilla/mux"
)


var Version string

var cache *badger.DB
func main() {

	app := cli.NewApp()
	app.Name = "blocklist-api"
	app.Version = Version
	app.Flags =  []cli.Flag{
		&cli.StringFlag{
			Name: "port",
			EnvVars: []string{"PORT"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "cache-path",
			EnvVars: []string{"CACHE_PATH"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "user",
			EnvVars: []string{"BLOCKLIST_USER"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "password",
			EnvVars: []string{"BLOCKLIST_PASS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:        "db-path",
			EnvVars:     []string{"DB_PATH"},
			Required:    true,
		},
	}
	app.Before = func(c *cli.Context) error {

		var err error

		cache, err = badger.Open(badger.DefaultOptions(c.String("cache-path")))
		if err != nil{
			return fmt.Errorf("error opening new cache %v", err)
		}

		if err := cache.Update(func(txn *badger.Txn) error {
			if err := txn.Set([]byte("credentials"), []byte(fmt.Sprintf("%s:%s", c.String("user"),
				c.String("password")))); err != nil{
				return fmt.Errorf("error")
			}
			return nil
		}); err != nil{
			return fmt.Errorf("error adding new cache %v", err)
		}

		return nil
	}
	app.Action = func(c *cli.Context) error {
		port := c.String("port")

		router := mux.NewRouter()
		router.Use(auth.Authorize(cache))

		db, err := database.NewBlocklistDB(c.String("db-path"))
		if err != nil{
			return fmt.Errorf("error opening new blocklistDB %v", err)
		}

		srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
			DB: db,
		}}))
		srv.AddTransport(transport.POST{})

		router.Handle("/", playground.Handler("Graphql Blocklist", "/query"))
		router.Handle("/graphql", srv)

		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

		errChan := make(chan error, 1)
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt)

		go func() {
			if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil{
				errChan <- err
				return
			}
		}()

		select {
		case <-sigChan:
			if err := cache.Close(); err != nil{
				return fmt.Errorf("error shutting down cache %v", err)
			}

			if err := db.DB.Close(); err != nil{
				return fmt.Errorf("error shutting down db %v", err)
			}

		case err := <- errChan:
			if err != nil{
				return fmt.Errorf("error received %v", err)
			}

			return nil
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil{
		log.Fatal(err)
	}

}
