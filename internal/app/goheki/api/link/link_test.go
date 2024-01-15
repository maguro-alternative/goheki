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
	)

	links := []Link{
		{
			Type:     "funart",
			URL:      "https://pixiv.com",
			Nsfw:     true,
			Darkness: false,
		},
		{
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

	t.Run("link作成", func(t *testing.T) {
		// テストの実行
		h := NewCreateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(links)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/link/create", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res []Link
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res, 2)
		assert.Equal(t, "funart", res[0].Type)
		assert.Equal(t, "https://pixiv.com", res[0].URL)
		assert.Equal(t, true, res[0].Nsfw)
		assert.Equal(t, false, res[0].Darkness)
		assert.Equal(t, "original", res[1].Type)
		assert.Equal(t, "https://pixiv.com", res[1].URL)
		assert.Equal(t, true, res[1].Nsfw)
		assert.Equal(t, true, res[1].Darkness)
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
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
		var res []Link
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res, 2)
		assert.Equal(t, links, res)
	})

	t.Run("link1件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d", *f.Links[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res []Link
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Equal(t, links[0], res[0])
	})

	t.Run("link2件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/link/read?id=%d&id=%d", *f.Links[0].ID, *f.Links[1].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res []Link
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Equal(t, links, res)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
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

	t.Run("link更新", func(t *testing.T) {
		// テストの実行
		h := NewUpdateHandler(indexService)
		// リクエストを作成
		lJson, err := json.Marshal(updateLinks)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/api/link/update", bytes.NewBuffer(lJson))
		assert.NoError(t, err)
		// レスポンスを作成
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res []Link
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res, 2)
		assert.Equal(t, updateLinks, res)
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
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "funart"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = false
		}))),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewLink(ctx, func(l *fixtures.Link) {
			l.Type = "original"
			l.URL = "https://pixiv.com"
			l.Nsfw = true
			l.Darkness = true
		}))),
	)

	delIDs := IDs{IDs:[]int64{*f.Links[0].ID, *f.Links[1].ID,}}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

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

		tx.RollbackCtx(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの検証
		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)

		assert.Len(t, res.IDs, 2)
		assert.Equal(t, delIDs, res)
	})
}
