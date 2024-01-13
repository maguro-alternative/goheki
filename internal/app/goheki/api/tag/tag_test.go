package tag

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

	"github.com/maguro-alternative/goheki/internal/app/goheki/model/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestCreateTagHandler(t *testing.T) {
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
}

func TestReadTagHandler(t *testing.T) {
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
	// データベースの準備
	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ1"
			s.Image = "https://example.com/image1.png"
			s.Content = "テスト内容1"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	t.Run("tag全件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/tag/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags []Tag
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, f.Tags[0].Name, tags[0].Name)
		assert.Equal(t, f.Tags[1].Name, tags[1].Name)
	})

	t.Run("tag1件取得", func(t *testing.T) {
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
		var ids []int64
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
		query = `
			SELECT
				id
			FROM
				tag
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/tag/read?id=%d",ids[0]), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals []Tag
		err = json.NewDecoder(w.Body).Decode(&actuals)
		assert.NoError(t, err)

		assert.Equal(t, tag[0].Name, actuals[0].Name)
	})

	t.Run("tag2件取得", func(t *testing.T) {
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
		query = `
			SELECT
				id
			FROM
				tag
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/tag/read?id=%d&id=%d", ids[0], ids[1]), nil)
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

}

func TestUpdateTagHandler(t *testing.T) {
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
		updateTag := []Tag{
			{
				Name: "テストタグ3",
			},
			{
				Name: "テストタグ4",
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
		eJson, err := json.Marshal(&updateTag)
		req, err := http.NewRequest(http.MethodPut, "/api/tag/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags []Tag
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, updateTag[0].Name, tags[0].Name)
		assert.Equal(t, updateTag[1].Name, tags[1].Name)
	})
}

func TestDeleteTagHandler(t *testing.T) {
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
		delIDs := IDs{IDs: []int64{ids[0]}}

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

		var actual IDs
		err = json.NewDecoder(w.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
