package source

import (
	"bytes"
	"context"
	"encoding/json"
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		entrys := []Entry{
			{
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		query := `
			INSERT INTO entry (
				name,
				image,
				content,
				created_at
			) VALUES (
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
		query = `
			SELECT
				id
			FROM
				entry
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)

		source := []Source{
			{
				EntryID: ids[0],
				Name:    "テストソース1",
				Url:     "https://example.com/image1.png",
				Type:    "anime",
			},
			{
				EntryID: ids[1],
				Name:    "テストソース2",
				Url:     "https://example.com/image2.png",
				Type:    "game",
			},
		}


		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&source)
		req, err := http.NewRequest(http.MethodPost, "/api/source/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var sources []Source
		err = json.Unmarshal(w.Body.Bytes(), &sources)
		assert.NoError(t, err)

		assert.Equal(t, source[0].Name, sources[0].Name)

		assert.Equal(t, source[1].Name, sources[1].Name)
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		entrys := []Entry{
			{
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		query := `
			INSERT INTO entry (
				name,
				image,
				content,
				created_at
			) VALUES (
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
		query = `
			SELECT
				id
			FROM
				entry
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)

		source := []Source{
			{
				EntryID: ids[0],
				Name:    "テストソース1",
				Url:     "https://example.com/image1.png",
				Type:    "anime",
			},
			{
				EntryID: ids[1],
				Name:    "テストソース2",
				Url:     "https://example.com/image2.png",
				Type:    "game",
			},
		}

		query = `
			INSERT INTO source (
				entry_id,
				name,
				url,
				type
			) VALUES (
				:entry_id,
				:name,
				:url,
				:type
			)
		`
		for _, s := range source {
			_, err = tx.NamedExecContext(ctx, query, s)
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

		var sources []Source
		err = json.Unmarshal(w.Body.Bytes(), &sources)
		assert.NoError(t, err)

		assert.Equal(t, source[0].Name, sources[0].Name)

		assert.Equal(t, source[1].Name, sources[1].Name)
	})
}
