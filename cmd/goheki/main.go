package main

import (
	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/pkg/db"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"
	"github.com/maguro-alternative/goheki/internal/app/goheki/article"
	"github.com/maguro-alternative/goheki/internal/app/goheki/middleware"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/entry"

	"context"
	"log"
	"net/http"

	"github.com/justinas/alice"
)

func main() {
	ctx := context.Background()
	// load env
	env, err := envconfig.NewEnv()
	if err != nil {
		log.Fatal(err)
	}
	indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	var indexService = service.NewIndexService(
		indexDB,
		cookie.Store,
		env,
	)

	// register routes
	mux := http.NewServeMux()
	middleChain := alice.New(middleware.CORS)
	mux.Handle("/", middleChain.Then(article.NewIndexHandler(indexService)))
	mux.Handle("/api/entry/create", middleChain.Then(entry.NewCreateHandler(indexService)))
	mux.Handle("/api/entry/read", middleChain.Then(entry.NewReadHandler(indexService)))
	//mux.Handle("/edit", middleChain.Then(article.Edit))
	//mux.Handle("/delete", middleChain.Then(article.Delete))

	log.Print("Server listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}