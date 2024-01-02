package tag

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"

	"github.com/maguro-alternative/goheki/pkg/db"

	"github.com/stretchr/testify/assert"
)

func TestTagHandler(t *testing.T) {
	t.Run("tag登録", func(t *testing.T) {
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
		tag := []Tag{
			{
				Name: "テストタグ1",
			},
			{
				Name: "テストタグ2",
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&tag)
		req, err := http.NewRequest(http.MethodPost, "/api/tag/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags []Tag
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, tag[0].Name, tags[0].Name)
		assert.Equal(t, tag[1].Name, tags[1].Name)
	})

	t.Run("tag取得", func(t *testing.T) {
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
		tag := []Tag{
			{
				Name: "テストタグ1",
			},
			{
				Name: "テストタグ2",
			},
		}

		query := `
			INSERT INTO tag (
				name
			) VALUES (
				:name
			)
		`
		for _, tag := range tag {
			_, err = tx.NamedExecContext(ctx, query, tag)
			assert.NoError(t, err)
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		eJson, err := json.Marshal(&tag)
		req, err := http.NewRequest(http.MethodGet, "/api/tag/read", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags []Tag
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, tag[0].Name, tags[0].Name)
		assert.Equal(t, tag[1].Name, tags[1].Name)
	})

	t.Run("tag更新", func(t *testing.T) {
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
		tag := []Tag{
			{
				Name: "テストタグ1",
			},
			{
				Name: "テストタグ2",
			},
		}

		query := `
			INSERT INTO tag (
				name
			) VALUES (
				:name
			)
		`
		for _, tag := range tag {
			_, err = tx.NamedExecContext(ctx, query, tag)
			assert.NoError(t, err)
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&tag)
		req, err := http.NewRequest(http.MethodPut, "/api/tag/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags []Tag
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, tag[0].Name, tags[0].Name)
		assert.Equal(t, tag[1].Name, tags[1].Name)
	})

	t.Run("tag削除", func(t *testing.T) {
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
		tag := []Tag{
			{
				Name: "テストタグ1",
			},
			{
				Name: "テストタグ2",
			},
		}
		var ids []int64

		query := `
			INSERT INTO tag (
				name
			) VALUES (
				:name
			)
		`
		for _, tag := range tag {
			_, err = tx.NamedExecContext(ctx, query, tag)
			assert.NoError(t, err)
		}
		selectQuery := `
			SELECT
				id
			FROM
				tag
		`
		err = tx.SelectContext(ctx, &ids, selectQuery)
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
		req, err := http.NewRequest(http.MethodDelete, "/api/tag/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual DeleteIDs
		err = json.NewDecoder(w.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
