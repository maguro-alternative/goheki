package main

import (
	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/pkg/db"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"
	"github.com/maguro-alternative/goheki/internal/app/goheki/article"
	"github.com/maguro-alternative/goheki/internal/app/goheki/middleware"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/entry"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/tag"

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
	mux.Handle("/api/entry/all-read", middleChain.Then(entry.NewAllReadHandler(indexService)))
	mux.Handle("/api/entry/get-read", middleChain.Then(entry.NewGetReadHandler(indexService)))
	mux.Handle("/api/entry/multiple-read", middleChain.Then(entry.NewMultipleReadHandler(indexService)))
	mux.Handle("/api/entry/update", middleChain.Then(entry.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry/delete", middleChain.Then(entry.NewDeleteHandler(indexService)))
	mux.Handle("/api/tag/create", middleChain.Then(tag.NewCreateHandler(indexService)))
	mux.Handle("/api/tag/all-read", middleChain.Then(tag.NewAllReadHandler(indexService)))
	mux.Handle("/api/tag/update", middleChain.Then(tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/tag/delete", middleChain.Then(tag.NewDeleteHandler(indexService)))

	log.Print("Server listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}