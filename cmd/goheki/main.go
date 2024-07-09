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
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/link"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/personality"
	"github.com/maguro-alternative/goheki/internal/app/goheki/api/source"
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

	mux.Handle("/api/bwh/create", middleChain.Append(middleware.BasicAuth).Then(bwh.NewCreateHandler(indexService)))
	mux.Handle("/api/bwh/read", middleChain.Append(middleware.BasicAuth).Then(bwh.NewReadHandler(indexService)))
	mux.Handle("/api/bwh/update", middleChain.Append(middleware.BasicAuth).Then(bwh.NewUpdateHandler(indexService)))
	mux.Handle("/api/bwh/delete", middleChain.Append(middleware.BasicAuth).Then(bwh.NewDeleteHandler(indexService)))

	mux.Handle("/api/entry/create", middleChain.Append(middleware.BasicAuth).Then(entry.NewCreateHandler(indexService)))
	mux.Handle("/api/entry/read", middleChain.Append(middleware.BasicAuth).Then(entry.NewReadHandler(indexService)))
	mux.Handle("/api/entry/update", middleChain.Append(middleware.BasicAuth).Then(entry.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry/delete", middleChain.Append(middleware.BasicAuth).Then(entry.NewDeleteHandler(indexService)))

	mux.Handle("/api/tag/create", middleChain.Append(middleware.BasicAuth).Then(tag.NewCreateHandler(indexService)))
	mux.Handle("/api/tag/read", middleChain.Append(middleware.BasicAuth).Then(tag.NewReadHandler(indexService)))
	mux.Handle("/api/tag/update", middleChain.Append(middleware.BasicAuth).Then(tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/tag/delete", middleChain.Append(middleware.BasicAuth).Then(tag.NewDeleteHandler(indexService)))

	mux.Handle("/api/entry_tag/create", middleChain.Append(middleware.BasicAuth).Then(entry_tag.NewCreateHandler(indexService)))
	mux.Handle("/api/entry_tag/read", middleChain.Append(middleware.BasicAuth).Then(entry_tag.NewReadHandler(indexService)))
	mux.Handle("/api/entry_tag/update", middleChain.Append(middleware.BasicAuth).Then(entry_tag.NewUpdateHandler(indexService)))
	mux.Handle("/api/entry_tag/delete", middleChain.Append(middleware.BasicAuth).Then(entry_tag.NewDeleteHandler(indexService)))

	mux.Handle("/api/haircolor/create", middleChain.Append(middleware.BasicAuth).Then(haircolor.NewCreateHandler(indexService)))
	mux.Handle("/api/haircolor/read", middleChain.Append(middleware.BasicAuth).Then(haircolor.NewReadHandler(indexService)))
	mux.Handle("/api/haircolor/update", middleChain.Append(middleware.BasicAuth).Then(haircolor.NewUpdateHandler(indexService)))
	mux.Handle("/api/haircolor/delete", middleChain.Append(middleware.BasicAuth).Then(haircolor.NewDeleteHandler(indexService)))

	mux.Handle("/api/hairlength/create", middleChain.Append(middleware.BasicAuth).Then(hairlength.NewCreateHandler(indexService)))
	mux.Handle("/api/hairlength/read", middleChain.Append(middleware.BasicAuth).Then(hairlength.NewReadHandler(indexService)))
	mux.Handle("/api/hairlength/update", middleChain.Append(middleware.BasicAuth).Then(hairlength.NewUpdateHandler(indexService)))
	mux.Handle("/api/hairlength/delete", middleChain.Append(middleware.BasicAuth).Then(hairlength.NewDeleteHandler(indexService)))

	mux.Handle("/api/hairstyle/create", middleChain.Append(middleware.BasicAuth).Then(hairstyle.NewCreateHandler(indexService)))
	mux.Handle("/api/hairstyle/read", middleChain.Append(middleware.BasicAuth).Then(hairstyle.NewReadHandler(indexService)))
	mux.Handle("/api/hairstyle/update", middleChain.Append(middleware.BasicAuth).Then(hairstyle.NewUpdateHandler(indexService)))
	mux.Handle("/api/hairstyle/delete", middleChain.Append(middleware.BasicAuth).Then(hairstyle.NewDeleteHandler(indexService)))

	mux.Handle("/api/heki_radar_chart/create", middleChain.Append(middleware.BasicAuth).Then(hekiradarchart.NewCreateHandler(indexService)))
	mux.Handle("/api/heki_radar_chart/read", middleChain.Append(middleware.BasicAuth).Then(hekiradarchart.NewReadHandler(indexService)))
	mux.Handle("/api/heki_radar_chart/update", middleChain.Append(middleware.BasicAuth).Then(hekiradarchart.NewUpdateHandler(indexService)))
	mux.Handle("/api/heki_radar_chart/delete", middleChain.Append(middleware.BasicAuth).Then(hekiradarchart.NewDeleteHandler(indexService)))

	mux.Handle("/api/link/create", middleChain.Append(middleware.BasicAuth).Then(link.NewCreateHandler(indexService)))
	mux.Handle("/api/link/read", middleChain.Append(middleware.BasicAuth).Then(link.NewReadHandler(indexService)))
	mux.Handle("/api/link/update", middleChain.Append(middleware.BasicAuth).Then(link.NewUpdateHandler(indexService)))
	mux.Handle("/api/link/delete", middleChain.Append(middleware.BasicAuth).Then(link.NewDeleteHandler(indexService)))

	mux.Handle("/api/personality/create", middleChain.Append(middleware.BasicAuth).Then(personality.NewCreateHandler(indexService)))
	mux.Handle("/api/personality/read", middleChain.Append(middleware.BasicAuth).Then(personality.NewReadHandler(indexService)))
	mux.Handle("/api/personality/update", middleChain.Append(middleware.BasicAuth).Then(personality.NewUpdateHandler(indexService)))
	mux.Handle("/api/personality/delete", middleChain.Append(middleware.BasicAuth).Then(personality.NewDeleteHandler(indexService)))

	mux.Handle("/api/source/create", middleChain.Append(middleware.BasicAuth).Then(source.NewCreateHandler(indexService)))
	mux.Handle("/api/source/read", middleChain.Append(middleware.BasicAuth).Then(source.NewReadHandler(indexService)))
	mux.Handle("/api/source/update", middleChain.Append(middleware.BasicAuth).Then(source.NewUpdateHandler(indexService)))
	mux.Handle("/api/source/delete", middleChain.Append(middleware.BasicAuth).Then(source.NewDeleteHandler(indexService)))

	log.Print("Server listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}