package entry

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

func TestCreateEntryHandler(t *testing.T) {
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
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}),
	)

	t.Run("entry登録失敗", func(t *testing.T) {
		entrys := []IDs{
			{
				IDs: []int64{1, 2},
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entrys)
		req, err := http.NewRequest(http.MethodPost, "/api/entry/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var actuals []Entry
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(actuals))
	})

	t.Run("entry登録", func(t *testing.T) {
		entriesJson := EntriesJson{
			[]Entry{
				{
					SourceID:  f.Sources[0].ID,
					Name:      "雪泉",
					Image:     "https://example.com/image1.png",
					Content:   "かわいい",
					CreatedAt: fixedTime,
				},
				{
					SourceID:  f.Sources[1].ID,
					Name:      "四条貴音",
					Image:     "https://example.com/image2.png",
					Content:   "お姫ちん",
					CreatedAt: fixedTime,
				},
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entriesJson)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/entry/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		r := w.Result()
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var res EntriesJson
		var actuals []Entry
		var count int
		err = json.NewDecoder(r.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Equal(t, entriesJson, res)

		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, entriesJson.Entries[0].SourceID, actuals[0].SourceID)
		assert.Equal(t, entriesJson.Entries[0].Name, actuals[0].Name)
		assert.Equal(t, entriesJson.Entries[0].Image, actuals[0].Image)
		assert.Equal(t, entriesJson.Entries[0].Content, actuals[0].Content)
		assert.Equal(t, entriesJson.Entries[1].SourceID, actuals[1].SourceID)
		assert.Equal(t, entriesJson.Entries[1].Name, actuals[1].Name)
		assert.Equal(t, entriesJson.Entries[1].Image, actuals[1].Image)
		assert.Equal(t, entriesJson.Entries[1].Content, actuals[1].Content)

		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestReadEntryHandler(t *testing.T) {
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
	)
	entrys := []Entry{
		{
			SourceID:  f.Sources[0].ID,
			Name:      f.Entrys[0].Name,
			Image:     f.Entrys[0].Image,
			Content:   f.Entrys[0].Content,
			CreatedAt: f.Entrys[0].CreatedAt,
		},
		{
			SourceID:  f.Sources[1].ID,
			Name:      f.Entrys[1].Name,
			Image:     f.Entrys[1].Image,
			Content:   f.Entrys[1].Content,
			CreatedAt: f.Entrys[1].CreatedAt,
		},
	}
	t.Run("entry全件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys[0].SourceID, actual.Entries[0].SourceID)
		assert.Equal(t, entrys[0].Name, actual.Entries[0].Name)
		assert.Equal(t, entrys[0].Image, actual.Entries[0].Image)
		assert.Equal(t, entrys[0].Content, actual.Entries[0].Content)
		assert.Equal(t, entrys[1].SourceID, actual.Entries[1].SourceID)
		assert.Equal(t, entrys[1].Name, actual.Entries[1].Name)
		assert.Equal(t, entrys[1].Image, actual.Entries[1].Image)
		assert.Equal(t, entrys[1].Content, actual.Entries[1].Content)
	})

	t.Run("entry1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d", f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys[0].SourceID, actual.Entries[0].SourceID)
		assert.Equal(t, entrys[0].Name, actual.Entries[0].Name)
		assert.Equal(t, entrys[0].Image, actual.Entries[0].Image)
		assert.Equal(t, entrys[0].Content, actual.Entries[0].Content)
	})

	t.Run("entry2件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=%d", f.Entrys[0].ID, f.Entrys[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys[0].SourceID, actual.Entries[0].SourceID)
		assert.Equal(t, entrys[0].Name, actual.Entries[0].Name)
		assert.Equal(t, entrys[0].Image, actual.Entries[0].Image)
		assert.Equal(t, entrys[0].Content, actual.Entries[0].Content)
		assert.Equal(t, entrys[1].SourceID, actual.Entries[1].SourceID)
		assert.Equal(t, entrys[1].Name, actual.Entries[1].Name)
		assert.Equal(t, entrys[1].Image, actual.Entries[1].Image)
		assert.Equal(t, entrys[1].Content, actual.Entries[1].Content)
	})

	t.Run("entry0件取得", func(t *testing.T) {
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
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(actual.Entries))
	})

	t.Run("entry存在しないidを含む2件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=%d", f.Entrys[0].ID, 0), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys[0].SourceID, actual.Entries[0].SourceID)
		assert.Equal(t, entrys[0].Name, actual.Entries[0].Name)
		assert.Equal(t, entrys[0].Image, actual.Entries[0].Image)
		assert.Equal(t, entrys[0].Content, actual.Entries[0].Content)
		assert.Equal(t, 1, len(actual.Entries))
	})

	t.Run("entry存在しないidを含む1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d", 0), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		// tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(actual.Entries))
	})

	t.Run("entry1件取得(形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read?id=aaa", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		h.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("entry2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=aaa&id=%d", f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		h.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
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
	)

	updateEntriesJson := EntriesJson{
		[]Entry{
			{
				ID:        f.Entrys[0].ID,
				SourceID:  f.Sources[0].ID,
				Name:      "テストエントリ3",
				Image:     "https://example.com/image3.png",
				Content:   "テスト内容3",
				CreatedAt: fixedTime,
			},
			{
				ID:        f.Entrys[1].ID,
				SourceID:  f.Sources[1].ID,
				Name:      "テストエントリ4",
				Image:     "https://example.com/image4.png",
				Content:   "テスト内容4",
				CreatedAt: fixedTime,
			},
		},
	}

	t.Run("entry更新失敗", func(t *testing.T) {
		entrys := []IDs{
			{
				IDs: []int64{1, 2},
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&entrys)
		req, err := http.NewRequest(http.MethodPut, "/api/entry/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var actuals []Entry
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
	})

	t.Run("entry更新", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateEntriesJson)
		req, err := http.NewRequest(http.MethodPut, "/api/entry/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual EntriesJson
		var actuals []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)
		assert.Equal(t, updateEntriesJson, actual)

		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, updateEntriesJson.Entries[0].SourceID, actuals[0].SourceID)
		assert.Equal(t, updateEntriesJson.Entries[0].Name, actuals[0].Name)
		assert.Equal(t, updateEntriesJson.Entries[0].Image, actuals[0].Image)
		assert.Equal(t, updateEntriesJson.Entries[0].Content, actuals[0].Content)
		assert.Equal(t, updateEntriesJson.Entries[1].SourceID, actuals[1].SourceID)
		assert.Equal(t, updateEntriesJson.Entries[1].Name, actuals[1].Name)
		assert.Equal(t, updateEntriesJson.Entries[1].Image, actuals[1].Image)
		assert.Equal(t, updateEntriesJson.Entries[1].Content, actuals[1].Content)
	})
}

func TestDeleteEntryHandler(t *testing.T) {
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
	)

	t.Run("entry削除失敗", func(t *testing.T) {
		entrys := EntriesJson{}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewDeleteHandler(indexService)
		eJson, err := json.Marshal(&entrys)
		req, err := http.NewRequest(http.MethodDelete, "/api/entry/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var actuals []Entry
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
	})

	t.Run("entry削除", func(t *testing.T) {
		ids := []int64{f.Entrys[0].ID, f.Entrys[1].ID}
		delIDs := IDs{IDs: ids}
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

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM entry")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
