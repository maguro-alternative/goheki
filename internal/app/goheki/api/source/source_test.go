package source

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"

	"github.com/maguro-alternative/goheki/pkg/db"

	"github.com/stretchr/testify/assert"
)

type Entry struct {
	ID        *int64    `db:"id" json:"id"`
	SourceID  int64     `db:"source_id" json:"source_id"`
	Name      string    `db:"name" json:"name"`
	Image     string    `db:"image" json:"image"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func TestSourceHandler(t *testing.T) {
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
		sources := []Source{
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
		}

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

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals []Source
		err = json.Unmarshal(w.Body.Bytes(), &sources)
		assert.NoError(t, err)

		assert.Equal(t, actuals[0].Name, sources[0].Name)

		assert.Equal(t, actuals[1].Name, sources[1].Name)
	})

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

		// テストデータの準備
		sources := []Source{
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
		}

		query := `
			INSERT INTO source (
				name,
				url,
				type
			) VALUES (
				:name,
				:url,
				:type
			)
		`
		for _, source := range sources {
			_, err = tx.NamedExecContext(ctx, query, source)
			assert.NoError(t, err)
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewAllReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/source/all-read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals []Source
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, actuals[0].Name, sources[0].Name)

		assert.Equal(t, actuals[1].Name, sources[1].Name)
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		sources := []Source{
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
		}

		query := `
			INSERT INTO source (
				name,
				url,
				type
			) VALUES (
				:name,
				:url,
				:type
			)
		`
		for _, source := range sources {
			_, err = tx.NamedExecContext(ctx, query, source)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				source
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		entrys := []Entry{
			{
				SourceID:  ids[0],
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				SourceID:  ids[1],
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entrys {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewGetReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/source/get-read?id=%d", ids[0]), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual Source
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err)

		assert.Equal(t, sources[0].Name, actual.Name)
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		sources := []Source{
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
		}

		query := `
			INSERT INTO source (
				name,
				url,
				type
			) VALUES (
				:name,
				:url,
				:type
			)
		`
		for _, source := range sources {
			_, err = tx.NamedExecContext(ctx, query, source)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				source
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)

		entrys := []Entry{
			{
				SourceID:  ids[0],
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				SourceID:  ids[1],
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entrys {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewMultipleReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/source/multiple-read?id=%d&id=%d", ids[0], ids[1]), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals []Source
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, sources[0].Name, actuals[0].Name)

		assert.Equal(t, sources[1].Name, actuals[1].Name)
	})

	t.Run("source更新", func(t *testing.T) {
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		sources := []Source{
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
		}
		updateSource := []Source{
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
		}

		query := `
			INSERT INTO source (
				name,
				url,
				type
			) VALUES (
				:name,
				:url,
				:type
			)
		`
		for _, source := range sources {
			_, err = tx.NamedExecContext(ctx, query, source)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				source
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		entrys := []Entry{
			{
				SourceID:  ids[0],
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				SourceID:  ids[1],
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entrys {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}

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

		var actuals []Source
		err = json.Unmarshal(w.Body.Bytes(), &actuals)
		assert.NoError(t, err)

		assert.Equal(t, updateSource[0].Name, actuals[0].Name)

		assert.Equal(t, updateSource[1].Name, actuals[1].Name)
	})
}
