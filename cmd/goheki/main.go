package main

import (
	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/pkg/db"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"
	"github.com/maguro-alternative/goheki/internal/app/goheki/article"
	"github.com/maguro-alternative/goheki/internal/app/goheki/middleware"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"

	"github.com/maguro-alternative/goheki/internal/app/goheki/api/bwh"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/entry"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/entry_tag"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/haircolor"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/hairlength"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/hairstyle"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/heki_radar_chart"
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

	mux.Handle("/api/bwh/create", middleChain.Then(bwh.NewCreateHandler(indexService)))
	mux.Handle("/api/bwh/read", middleChain.Then(bwh.NewReadHandler(indexService)))
	mux.Handle("/api/bwh/update", middleChain.Then(bwh.NewUpdateHandler(indexService)))
	mux.Handle("/api/bwh/delete", middleChain.Then(bwh.NewDeleteHandler(indexService)))
	mux.Handle("/api/entry/create", middleChain.Then(entry.NewCreateHandler(indexService)))
	mux.Handle("/api/entry/read", middleChain.Then(entry.NewReadHandler(indexService)))
	mux.Handle("/api/entry/update", middleChain.Then(entry.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry/delete", middleChain.Then(entry.NewDeleteHandler(indexService)))
	mux.Handle("/api/tag/create", middleChain.Then(tag.NewCreateHandler(indexService)))
	mux.Handle("/api/tag/read", middleChain.Then(tag.NewReadHandler(indexService)))
	mux.Handle("/api/tag/update", middleChain.Then(tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/tag/delete", middleChain.Then(tag.NewDeleteHandler(indexService)))
	mux.Handle("/api/entry_tag/create", middleChain.Then(entry_tag.NewCreateHandler(indexService)))
	mux.Handle("/api/entry_tag/read", middleChain.Then(entry_tag.NewReadHandler(indexService)))
	mux.Handle("/api/entry_tag/update", middleChain.Then(entry_tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry_tag/delete", middleChain.Then(entry_tag.NewDeleteHandler(indexService)))
	mux.Handle("/api/haircolor/create", middleChain.Then(haircolor.NewCreateHandler(indexService)))
	mux.Handle("/api/haircolor/read", middleChain.Then(haircolor.NewReadHandler(indexService)))
	mux.Handle("/api/haircolor/update", middleChain.Then(haircolor.NewUpdateHandler(indexService)))
	mux.Handle("/api/haircolor/delete", middleChain.Then(haircolor.NewDeleteHandler(indexService)))
	mux.Handle("/api/hairlength/create", middleChain.Then(hairlength.NewCreateHandler(indexService)))
	mux.Handle("/api/hairlength/read", middleChain.Then(hairlength.NewReadHandler(indexService)))
	mux.Handle("/api/hairlength/update", middleChain.Then(hairlength.NewUpdateHandler(indexService)))
	mux.Handle("/api/hairlength/delete", middleChain.Then(hairlength.NewDeleteHandler(indexService)))
	mux.Handle("/api/hairstyle/create", middleChain.Then(hairstyle.NewCreateHandler(indexService)))
	mux.Handle("/api/hairstyle/read", middleChain.Then(hairstyle.NewReadHandler(indexService)))
	mux.Handle("/api/hairstyle/update", middleChain.Then(hairstyle.NewUpdateHandler(indexService)))
	mux.Handle("/api/hairstyle/delete", middleChain.Then(hairstyle.NewDeleteHandler(indexService)))
	mux.Handle("/api/heki_radar_chart/create", middleChain.Then(hekiradarchart.NewCreateHandler(indexService)))

	log.Print("Server listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}