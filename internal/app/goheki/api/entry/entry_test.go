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
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read", bytes.NewBuffer(eJson))
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

	t.Run("entry更新", func(t *testing.T) {
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
		updateEntrys := []Entry{
			{
				Name:      "テストエントリ3",
				Image:     "https://example.com/image3.png",
				Content:   "テスト内容3",
				CreatedAt: fixedTime,
			},
			{
				Name:      "テストエントリ4",
				Image:     "https://example.com/image4.png",
				Content:   "テスト内容4",
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
		var ids []int64
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		for i, id := range ids {
			updateEntrys[i].ID = &id
		}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateEntrys)
		req, err := http.NewRequest(http.MethodPut, "/api/entry/update", bytes.NewBuffer(eJson))
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

		assert.Equal(t, updateEntrys, actual)
	})

	t.Run("entry削除", func(t *testing.T) {
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
		query = `
			SELECT
				id
			FROM
				entry
		`
		var ids []int64
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		delIDs := DeleteIDs{IDs: ids}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewDeleteHandler(indexService)
		eJson, err := json.Marshal(&delIDs)
		req, err := http.NewRequest(http.MethodDelete, "/api/entry/delete", bytes.NewBuffer(eJson))
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

		var actual DeleteIDs
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
