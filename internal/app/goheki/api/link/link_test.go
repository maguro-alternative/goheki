package link

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

func TestCreateLinkHandler(t *testing.T) {
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

	links := []Link{
		{
			EntryID:  f.Entrys[0].ID,
			Type:     "funart",
			URL:      "https://pixiv.com",
			Nsfw:     true,
			Darkness: false,
		},
		{
			EntryID:  f.Entrys[1].ID,
			Type:     "original",
			URL:      "https://pixiv.com",
			Nsfw:     true,
			Darkness: true,
		},
	}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("link作成失敗", func(t *testing.T) {
		// テストの実行
		h := NewCreateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(LinksJson{})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/link/create", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		// レスポンスの検証
		var actual []Link
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM link")
		assert.NoError(t, err)

		assert.Len(t, actual, 0)
	})

	t.Run("link作成", func(t *testing.T) {
		// テストの実行
		h := NewCreateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(LinksJson{links})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/link/create", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 2)
		assert.Equal(t, "funart", res.Links[0].Type)
		assert.Equal(t, "https://pixiv.com", res.Links[0].URL)
		assert.Equal(t, true, res.Links[0].Nsfw)
		assert.Equal(t, false, res.Links[0].Darkness)
		assert.Equal(t, "original", res.Links[1].Type)
		assert.Equal(t, "https://pixiv.com", res.Links[1].URL)
		assert.Equal(t, true, res.Links[1].Nsfw)
		assert.Equal(t, true, res.Links[1].Darkness)

		var resLinks []Link
		err = tx.SelectContext(ctx, &resLinks, "SELECT * FROM link")
		assert.NoError(t, err)

		assert.Len(t, resLinks, 2)
		assert.Equal(t, "funart", resLinks[0].Type)
		assert.Equal(t, "https://pixiv.com", resLinks[0].URL)
		assert.Equal(t, true, resLinks[0].Nsfw)
		assert.Equal(t, false, resLinks[0].Darkness)
		assert.Equal(t, "original", resLinks[1].Type)
		assert.Equal(t, "https://pixiv.com", resLinks[1].URL)
		assert.Equal(t, true, resLinks[1].Nsfw)
		assert.Equal(t, true, resLinks[1].Darkness)
	})
}

func TestReadLinkHandler(t *testing.T) {
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "original"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = true
		}))),
	)

	links := []Link{
		{
			EntryID:  f.Links[0].EntryID,
			Type:     f.Links[0].Type,
			URL:      f.Links[0].URL,
			Nsfw:     f.Links[0].Nsfw,
			Darkness: f.Links[0].Darkness,
		},
		{
			EntryID:  f.Links[1].EntryID,
			Type:     f.Links[1].Type,
			URL:      f.Links[1].URL,
			Nsfw:     f.Links[1].Nsfw,
			Darkness: f.Links[1].Darkness,
		},
	}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("link全件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, "/api/link/read", nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 2)
		assert.Equal(t, links, res.Links)
	})

	t.Run("link1件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d", f.Links[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 1)
		assert.Equal(t, links[0], res.Links[0])
	})

	t.Run("link2件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d&id=%d", f.Links[0].ID, f.Links[1].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 2)
		assert.Equal(t, links, res.Links)
	})

	t.Run("link1件取得(存在しない)", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, "/api/link/read?id=0", nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 0)
	})

	t.Run("link2件取得(内1件存在しない)", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d&id=0", f.Links[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 1)
		assert.Equal(t, links[0], res.Links[0])
	})

	t.Run("link1件取得(形式が正しくない)", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%s", "a"), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("link2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d&id=%s", f.Links[0].ID, "a"), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateLinkHandler(t *testing.T) {
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "original"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = true
		}))),
	)

	updateLinks := []Link{
		{
			ID:       f.Links[0].ID,
			EntryID:  f.Links[0].EntryID,
			Type:     "funart",
			URL:      "https://pixiv.com",
			Nsfw:     true,
			Darkness: true,
		},
		{
			ID:       f.Links[1].ID,
			EntryID:  f.Links[1].EntryID,
			Type:     "original",
			URL:      "https://pixiv.com",
			Nsfw:     true,
			Darkness: false,
		},
	}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("link更新失敗", func(t *testing.T) {
		// テストの実行
		h := NewUpdateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(LinksJson{})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/api/link/update", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		// レスポンスの検証
		var actual []Link
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM link")
		assert.NoError(t, err)

		assert.Len(t, actual, 2)
		assert.Equal(t, f.Links[0].EntryID, actual[0].EntryID)
		assert.Equal(t, f.Links[0].Type, actual[0].Type)
		assert.Equal(t, f.Links[0].URL, actual[0].URL)
		assert.Equal(t, f.Links[0].Nsfw, actual[0].Nsfw)
		assert.Equal(t, f.Links[0].Darkness, actual[0].Darkness)
		assert.Equal(t, f.Links[1].EntryID, actual[1].EntryID)
		assert.Equal(t, f.Links[1].Type, actual[1].Type)
		assert.Equal(t, f.Links[1].URL, actual[1].URL)
		assert.Equal(t, f.Links[1].Nsfw, actual[1].Nsfw)
	})

	t.Run("link更新", func(t *testing.T) {
		// テストの実行
		h := NewUpdateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(LinksJson{updateLinks})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/api/link/update", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res LinksJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.Links, 2)
		assert.Equal(t, updateLinks, res.Links)
	})
}

func TestDeleteLinkHandler(t *testing.T) {
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "アイドルマスター"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "四条貴音"
			s.Image = "https://example.com/image2.png"
			s.Content = "お姫ちん"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "original"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = true
		}))),
	)

	delIDs := IDs{IDs:[]int64{f.Links[0].ID, f.Links[1].ID,}}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("link削除失敗", func(t *testing.T) {
		// テストの実行
		h := NewDeleteHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(LinksJson{})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, "/api/link/delete", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		// レスポンスの検証
		var actual []Link
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM link")
		assert.NoError(t, err)

		assert.Len(t, actual, 2)
		assert.Equal(t, f.Links[0].EntryID, actual[0].EntryID)
		assert.Equal(t, f.Links[0].Type, actual[0].Type)
		assert.Equal(t, f.Links[0].URL, actual[0].URL)
		assert.Equal(t, f.Links[0].Nsfw, actual[0].Nsfw)
		assert.Equal(t, f.Links[0].Darkness, actual[0].Darkness)
		assert.Equal(t, f.Links[1].EntryID, actual[1].EntryID)
		assert.Equal(t, f.Links[1].Type, actual[1].Type)
		assert.Equal(t, f.Links[1].URL, actual[1].URL)
		assert.Equal(t, f.Links[1].Nsfw, actual[1].Nsfw)
	})

	t.Run("link削除", func(t *testing.T) {
		// テストの実行
		h := NewDeleteHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, "/api/link/delete", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.IDs, 2)
		assert.Equal(t, delIDs, res)
	})
}
