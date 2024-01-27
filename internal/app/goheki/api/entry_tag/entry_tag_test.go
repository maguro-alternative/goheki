package entry_tag

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

func TestCreateEntryTagHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)

	t.Run("entry_tag登録失敗", func(t *testing.T) {
		entryTags := IDs{
			IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entryTags)
		req, err := http.NewRequest(http.MethodPost, "/api/entry_tag/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var actual []EntryTag
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM entry_tag")
		assert.NoError(t, err)

		assert.Equal(t, 0, len(actual))
	})

	t.Run("entry_tag登録", func(t *testing.T) {
		entryTags := []EntryTag{
			{
				EntryID: f.Entrys[0].ID,
				TagID:   f.Tags[0].ID,
			},
			{
				EntryID: f.Entrys[1].ID,
				TagID:   f.Tags[1].ID,
			},
		}
		entryTagsJson := EntryTagsJson{EntryTags: entryTags}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entryTagsJson)
		req, err := http.NewRequest(http.MethodPost, "/api/entry_tag/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		var actual []EntryTag
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		err = tx.SelectContext(ctx, &actual, "SELECT * FROM entry_tag")
		assert.NoError(t, err)

		assert.Equal(t, entryTags[0].EntryID, actual[0].EntryID)
		assert.Equal(t, entryTags[0].TagID, actual[0].TagID)
		assert.Equal(t, entryTags[1].EntryID, actual[1].EntryID)
		assert.Equal(t, entryTags[1].TagID, actual[1].TagID)
	})
}

func TestReadEntryTagHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	f.Build(t,
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[0].ID
			s.TagID = f.Tags[0].ID
		}),
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[1].ID
			s.TagID = f.Tags[1].ID
		}),
	)
	t.Run("entry_tag全件取得", func(t *testing.T) {
		entryTags := []EntryTag{
			{
				ID:      f.EntryTags[0].ID,
				EntryID: f.Entrys[0].ID,
				TagID:   f.Tags[0].ID,
			},
			{
				ID:      f.EntryTags[1].ID,
				EntryID: f.Entrys[1].ID,
				TagID:   f.Tags[1].ID,
			},
		}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/entry_tag/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, entryTags, res.EntryTags)
	})

	t.Run("entry_tag1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d", f.EntryTags[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, f.EntryTags[0].ID, res.EntryTags[0].ID)
	})

	t.Run("entry_tag2件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=%d", f.EntryTags[0].ID, f.EntryTags[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, f.EntryTags[0].ID, res.EntryTags[0].ID)
		assert.Equal(t, f.EntryTags[1].ID, res.EntryTags[1].ID)
	})

	t.Run("entry_tag1件取得(存在しない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read?id=0", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(res.EntryTags))
	})

	t.Run("entry_tag2件取得(内1件は存在しない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=1", f.EntryTags[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntryTagsJson
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, f.EntryTags[0].ID, res.EntryTags[0].ID)
	})

	t.Run("entry_tag1件取得(形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read?id=aaa", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusInternalServerError, r.StatusCode)
	})

	t.Run("entry_tad2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=aaa", f.EntryTags[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusInternalServerError, r.StatusCode)
	})
}

func TestUpdateEntryHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	f.Build(t,
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[0].ID
			s.TagID = f.Tags[0].ID
		}),
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[1].ID
			s.TagID = f.Tags[1].ID
		}),
	)

	t.Run("entry_tag更新失敗", func(t *testing.T) {
		updateEntryTags := EntryTag{}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateEntryTags)
		req, err := http.NewRequest(http.MethodPut, "/api/entry_tag/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var actuals []EntryTag
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry_tag")
		assert.NoError(t, err)

		assert.Equal(t, actuals[0].EntryID, f.EntryTags[0].EntryID)
		assert.Equal(t, actuals[0].TagID, f.EntryTags[0].TagID)
		assert.Equal(t, actuals[1].EntryID, f.EntryTags[1].EntryID)
		assert.Equal(t, actuals[1].TagID, f.EntryTags[1].TagID)
	})

	t.Run("entry_tag更新", func(t *testing.T) {
		updateEntryTags := EntryTagsJson{[]EntryTag{
			{
				ID:      f.EntryTags[0].ID,
				EntryID: f.Entrys[0].ID,
				TagID:   f.Tags[1].ID,
			},
			{
				ID:      f.EntryTags[1].ID,
				EntryID: f.Entrys[1].ID,
				TagID:   f.Tags[0].ID,
			},
		}}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateEntryTags)
		req, err := http.NewRequest(http.MethodPut, "/api/entry_tag/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actuals EntryTagsJson
		var actual []EntryTag
		err = json.NewDecoder(res.Body).Decode(&actuals)
		assert.NoError(t, err)
		assert.Equal(t, updateEntryTags, actuals)

		err = tx.SelectContext(ctx, &actual, "SELECT * FROM entry_tag")
		assert.NoError(t, err)

		assert.Equal(t, actual[0].EntryID, f.Entrys[0].ID)
		assert.Equal(t, actual[0].TagID, f.Tags[1].ID)
		assert.Equal(t, actual[1].EntryID, f.Entrys[1].ID)
		assert.Equal(t, actual[1].TagID, f.Tags[0].ID)
	})
}

func TestDeleteEntryTagHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	f.Build(t,
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[0].ID
			s.TagID = f.Tags[0].ID
		}),
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[1].ID
			s.TagID = f.Tags[1].ID
		}),
	)

	t.Run("entry_tag削除失敗", func(t *testing.T) {
		delIDs := EntryTagsJson{}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewDeleteHandler(indexService)
		eJson, err := json.Marshal(&delIDs)
		req, err := http.NewRequest(http.MethodDelete, "/api/entry_tag/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var actuals []EntryTag
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry_tag")
		assert.NoError(t, err)

		assert.Equal(t, actuals[0].EntryID, f.EntryTags[0].EntryID)
		assert.Equal(t, actuals[0].TagID, f.EntryTags[0].TagID)
		assert.Equal(t, actuals[1].EntryID, f.EntryTags[1].EntryID)
		assert.Equal(t, actuals[1].TagID, f.EntryTags[1].TagID)
	})

	t.Run("entry_tag削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{f.EntryTags[0].ID, f.EntryTags[1].ID}}
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

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual IDs
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
