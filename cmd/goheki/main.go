package main

import (
	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/pkg/db"
	"github.com/maguro-alternative/goheki/pkg/cookie"
	"github.com/maguro-alternative/goheki/internal/app/goheki/article"
	"github.com/maguro-alternative/goheki/internal/app/goheki/middleware"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"

	"log"
	"net/http"

	"github.com/justinas/alice"
)

func main() {
	// load env
	env, err := envconfig.NewEnv()
	if err != nil {
		log.Fatal(err)
	}
	indexDB, err := db.NewPostgresDB(env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	var indexService = service.NewIndexService(
		indexDB,
		cookie.Store,
		env,
	)

	// register routes
	mux := http.NewServeMux()
	middleChain := alice.New(middleware.CORS)
	mux.Handle("/", middleChain.Then(article.NewIndexHandler(indexService)))
	//mux.Handle("/show", middleChain.Then(article.Show))
	//mux.Handle("/create", middleChain.Then(article.Create))
	//mux.Handle("/edit", middleChain.Then(article.Edit))
	//mux.Handle("/delete", middleChain.Then(article.Delete))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}