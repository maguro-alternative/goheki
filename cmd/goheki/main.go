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

	_ "embed"
	"context"
	"log"
	"net/http"

	"github.com/justinas/alice"
)

//go:embed schema.sql
var schema string // schema.sqlの内容をschemaに代入

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
	// テーブルの作成
	if _, err := indexDB.ExecContext(ctx, schema); err != nil {
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
	mux.Handle("/api/entry/update", middleChain.Then(entry.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry/delete", middleChain.Then(entry.NewDeleteHandler(indexService)))
	mux.Handle("/api/tag/create", middleChain.Then(tag.NewCreateHandler(indexService)))
	mux.Handle("/api/tag/read", middleChain.Then(tag.NewReadHandler(indexService)))
	mux.Handle("/api/tag/update", middleChain.Then(tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/tag/delete", middleChain.Then(tag.NewDeleteHandler(indexService)))

	log.Print("Server listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}