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

		// ロールバック
		defer tx.RollbackCtx(ctx)
		// テストデータの準備
		tag := TagsJson{[]Tag{
			{
				Name: "テストタグ1",
			},
			{
				Name: "テストタグ2",
			},
		}}

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

		assert.Equal(t, http.StatusOK, w.Code)

		var tags TagsJson
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, tag, tags)

		var actuals []Tag
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM tag")
		assert.NoError(t, err)

		assert.Equal(t, tag.Tags[0].Name, actuals[0].Name)
		assert.Equal(t, tag.Tags[1].Name, actuals[1].Name)
	})

	t.Run("tag登録失敗", func(t *testing.T) {
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
		// ロールバック
		defer tx.RollbackCtx(ctx)
		// テストデータの準備
		tag := TagsJson{}

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

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actuals []Tag
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM tag")
		assert.NoError(t, err)

		assert.Equal(t, 0, len(actuals))
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags TagsJson
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, f.Tags[0].Name, tags.Tags[0].Name)
		assert.Equal(t, f.Tags[1].Name, tags.Tags[1].Name)
	})

	t.Run("tag1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/tag/read?id=%d",f.Tags[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals TagsJson
		err = json.NewDecoder(w.Body).Decode(&actuals)
		assert.NoError(t, err)

		assert.Equal(t, f.Tags[0].Name, actuals.Tags[0].Name)
	})

	t.Run("tag2件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/tag/read?id=%d&id=%d", f.Tags[0].ID, f.Tags[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags TagsJson
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, f.Tags[0].Name, tags.Tags[0].Name)
		assert.Equal(t, f.Tags[1].Name, tags.Tags[1].Name)
	})

	t.Run("tag1件取得(存在しない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/tag/read?id=0", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var actuals TagsJson
		err = json.NewDecoder(w.Body).Decode(&actuals)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(actuals.Tags))
	})

	t.Run("tag2件取得(内1件存在しない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/tag/read?id=%d&id=0", f.Tags[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var tags TagsJson
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, f.Tags[0].Name, tags.Tags[0].Name)
		assert.Equal(t, 1, len(tags.Tags))
	})

	t.Run("tag1件取得(形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/tag/read?id=aaa", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("tag2件取得(内1件形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/tag/read?id=aaa&id=1", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateTagHandler(t *testing.T) {
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

	// ロールバック
	defer tx.RollbackCtx(ctx)
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

	t.Run("tag更新失敗", func(t *testing.T) {
		updateTag := TagsJson{[]Tag{}}

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

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var tags []Tag
		err = tx.SelectContext(ctx, &tags, "SELECT * FROM tag")
		assert.NoError(t, err)

		assert.Equal(t, 2, len(tags))

		assert.Equal(t, f.Tags[0].Name, tags[0].Name)
		assert.Equal(t, f.Tags[1].Name, tags[1].Name)
	})

	t.Run("tag更新", func(t *testing.T) {
		updateTag := TagsJson{[]Tag{
			{
				Name: "テストタグ3",
			},
			{
				Name: "テストタグ4",
			},
		}}

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

		assert.Equal(t, http.StatusOK, w.Code)

		var tags TagsJson
		err = json.NewDecoder(w.Body).Decode(&tags)
		assert.NoError(t, err)

		assert.Equal(t, updateTag, tags)

		var actuals []Tag
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM tag")
		assert.NoError(t, err)

		assert.Equal(t, updateTag.Tags[0].Name, actuals[0].Name)
		assert.Equal(t, updateTag.Tags[1].Name, actuals[1].Name)
	})
}

func TestDeleteTagHandler(t *testing.T) {
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

	// ロールバック
	defer tx.RollbackCtx(ctx)
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

	t.Run("tag削除失敗", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{}}

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

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var tags []Tag
		err = tx.SelectContext(ctx, &tags, "SELECT * FROM tag")
		assert.NoError(t, err)

		assert.Equal(t, 2, len(tags))

		assert.Equal(t, f.Tags[0].Name, tags[0].Name)
		assert.Equal(t, f.Tags[1].Name, tags[1].Name)
	})

	t.Run("tag削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{f.Tags[0].ID, f.Tags[1].ID}}

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

		assert.Equal(t, http.StatusOK, w.Code)

		var actual IDs
		err = json.NewDecoder(w.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
