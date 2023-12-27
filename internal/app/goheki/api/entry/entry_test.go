package entry

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

func TestEntryHandler(t *testing.T) {
	t.Run("entry登録", func(t *testing.T) {
		ctx := context.Background()
		env, err := envconfig.NewEnv()
		assert.NoError(t, err)
		// データベースに接続
		xDB, cleanup, err := db.NewDBV1(ctx, "postgres", env.DatabaseURL)
		assert.NoError(t, err)
		defer cleanup()
		// トランザクションの開始
		tx, err := xDB.BeginTxx(ctx, nil)
		assert.NoError(t, err)
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		entry := []Entry{
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

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entry)
		req, err := http.NewRequest(http.MethodPost, "/api/entry/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entry, actual)
	})

	t.Run("entry取得", func(t *testing.T) {
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
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
		entry := []Entry{
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
		for _, entry := range entry {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		eJson, err := json.Marshal(&entry)
		req, err := http.NewRequest(http.MethodPost, "/api/entry/read", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		err = json.NewDecoder(req.Body).Decode(&entry)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entry, actual)
	})
}
