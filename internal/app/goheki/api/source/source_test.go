package source

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"

	"github.com/maguro-alternative/goheki/pkg/db"

	"github.com/maguro-alternative/goheki/internal/app/goheki/model/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestCreateSourceHandler(t *testing.T) {
	t.Run("source登録", func(t *testing.T) {
		ctx := context.Background()
		env, err := envconfig.NewEnv()
		assert.NoError(t, err)
		// データベースに接続
		indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
		assert.NoError(t, err)
		defer cleanup()
		// トランザクションの開始
		tx, err := indexDB.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		// テストデータの準備
		sources := SourcesJson{[]Source{
			{
				Name: "テストソース1",
				Url:  "https://example.com/image1.png",
				Type: "anime",
			},
			{
				Name: "テストソース2",
				Url:  "https://example.com/image2.png",
				Type: "game",
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&sources)
		req, err := http.NewRequest(http.MethodPost, "/api/source/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals SourcesJson
		var actual []Source
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, actuals, sources)

		err = tx.SelectContext(ctx, &actual, "SELECT * FROM source")
		assert.NoError(t, err)

		assert.Equal(t, sources.Sources[0].Name, actual[0].Name)
		assert.Equal(t, sources.Sources[0].Url, actual[0].Url)
		assert.Equal(t, sources.Sources[0].Type, actual[0].Type)
		assert.Equal(t, sources.Sources[1].Name, actual[1].Name)
		assert.Equal(t, sources.Sources[1].Url, actual[1].Url)
		assert.Equal(t, sources.Sources[1].Type, actual[1].Type)

		tx.RollbackCtx(ctx)
	})
}

func TestReadSourceHandler(t *testing.T) {
	t.Run("source全件取得", func(t *testing.T) {
		ctx := context.Background()
		env, err := envconfig.NewEnv()
		assert.NoError(t, err)
		// データベースに接続
		indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
		assert.NoError(t, err)
		defer cleanup()
		// トランザクションの開始
		tx, err := indexDB.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		// データベースの準備
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース1"
				s.Url = "https://example.com/image1.png"
				s.Type = "anime"
			}),
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース2"
				s.Url = "https://example.com/image2.png"
				s.Type = "game"
			}),
		)

		// テストデータの準備
		sources := SourcesJson{[]Source{
			{
				Name: f.Sources[0].Name,
				Url:  f.Sources[0].Url,
				Type: f.Sources[0].Type,
			},
			{
				Name: f.Sources[1].Name,
				Url:  f.Sources[1].Url,
				Type: f.Sources[1].Type,
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/source/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals SourcesJson
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, sources.Sources[0].Name, actuals.Sources[0].Name)
		assert.Equal(t, sources.Sources[0].Url, actuals.Sources[0].Url)
		assert.Equal(t, sources.Sources[0].Type, actuals.Sources[0].Type)
		assert.Equal(t, sources.Sources[1].Name, actuals.Sources[1].Name)
		assert.Equal(t, sources.Sources[1].Url, actuals.Sources[1].Url)
		assert.Equal(t, sources.Sources[1].Type, actuals.Sources[1].Type)
	})

	t.Run("source1件取得", func(t *testing.T) {
		ctx := context.Background()
		env, err := envconfig.NewEnv()
		assert.NoError(t, err)
		// データベースに接続
		indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
		assert.NoError(t, err)
		defer cleanup()
		// トランザクションの開始
		tx, err := indexDB.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		// データベースの準備
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース1"
				s.Url = "https://example.com/image1.png"
				s.Type = "anime"
			}),
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース2"
				s.Url = "https://example.com/image2.png"
				s.Type = "game"
			}),
		)

		// テストデータの準備
		sources := SourcesJson{[]Source{
			{
				Name: f.Sources[0].Name,
				Url:  f.Sources[0].Url,
				Type: f.Sources[0].Type,
			},
			{
				Name: f.Sources[1].Name,
				Url:  f.Sources[1].Url,
				Type: f.Sources[1].Type,
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/source/read?id=%d", f.Sources[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals SourcesJson
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, sources.Sources[0].Name, actuals.Sources[0].Name)
		assert.Equal(t, sources.Sources[0].Url, actuals.Sources[0].Url)
		assert.Equal(t, sources.Sources[0].Type, actuals.Sources[0].Type)
	})

	t.Run("source2件取得", func(t *testing.T) {
		ctx := context.Background()
		env, err := envconfig.NewEnv()
		assert.NoError(t, err)
		// データベースに接続
		indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
		assert.NoError(t, err)
		defer cleanup()
		// トランザクションの開始
		tx, err := indexDB.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		// データベースの準備
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース1"
				s.Url = "https://example.com/image1.png"
				s.Type = "anime"
			}),
			fixtures.NewSource(ctx, func(s *fixtures.Source) {
				s.Name = "テストソース2"
				s.Url = "https://example.com/image2.png"
				s.Type = "game"
			}),
		)

		// テストデータの準備
		sources := SourcesJson{[]Source{
			{
				Name: f.Sources[0].Name,
				Url:  f.Sources[0].Url,
				Type: f.Sources[0].Type,
			},
			{
				Name: f.Sources[1].Name,
				Url:  f.Sources[1].Url,
				Type: f.Sources[1].Type,
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/source/read?id=%d&id=%d", f.Sources[0].ID, f.Sources[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals SourcesJson
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, sources.Sources[0].Name, actuals.Sources[0].Name)
		assert.Equal(t, sources.Sources[0].Url, actuals.Sources[0].Url)
		assert.Equal(t, sources.Sources[0].Type, actuals.Sources[0].Type)
		assert.Equal(t, sources.Sources[1].Name, actuals.Sources[1].Name)
		assert.Equal(t, sources.Sources[1].Url, actuals.Sources[1].Url)
		assert.Equal(t, sources.Sources[1].Type, actuals.Sources[1].Type)
	})
}

func TestUpdateSourceHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	// データベースに接続
	indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
	assert.NoError(t, err)
	defer cleanup()
	// トランザクションの開始
	tx, err := indexDB.BeginTxx(ctx, nil)
	assert.NoError(t, err)
	// データベースの準備
	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}),
	)
	t.Run("source更新", func(t *testing.T) {
		// テストデータの準備
		updateSource := SourcesJson{[]Source{
			{
				Name: "テストソース1更新",
				Url:  "https://example.com/image1.png",
				Type: "anime",
			},
			{
				Name: "テストソース2更新",
				Url:  "https://example.com/image2.png",
				Type: "game",
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateSource)
		req, err := http.NewRequest(http.MethodPut, "/api/source/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals SourcesJson
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, updateSource, actuals)
	})
}

func TestDeleteSourceHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	// データベースに接続
	indexDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
	assert.NoError(t, err)
	defer cleanup()
	// トランザクションの開始
	tx, err := indexDB.BeginTxx(ctx, nil)
	assert.NoError(t, err)
	// データベースの準備
	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}),
	)
	t.Run("source削除", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewDeleteHandler(indexService)
		eJson, err := json.Marshal(&IDs{
			IDs: []int64{f.Sources[0].ID},
		})
		req, err := http.NewRequest(http.MethodDelete, "/api/source/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual IDs
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err)

		assert.Equal(t, f.Sources[0].ID, actual.IDs[0])
	})
}
