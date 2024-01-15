package hairstyle

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

func TestCreateHairStyleHandler(t *testing.T) {
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

	// テストデータの準備
	hairStyles := []HairStyle{
		{
			EntryID: f.Entrys[0].ID,
			Style:   "ショートカット",
		},
		{
			EntryID: f.Entrys[1].ID,
			Style:   "ウェービーロングヘアー",
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairstyle登録", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		body, err := json.Marshal(hairStyles)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/hairstyles/create", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairStyle
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, hairStyles, res)
	})
}

func TestReadHairStyleHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ショートカット"
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
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ウェービーロングヘアー"
		}))),
	)

	// テストデータの準備
	hairStyles := []HairStyle{
		{
			EntryID: f.Entrys[0].ID,
			Style:   f.HairStyles[0].Style,
		},
		{
			EntryID: f.Entrys[1].ID,
			Style:   f.HairStyles[1].Style,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairstyle全件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/hairstyles/read", nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairStyle
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, hairStyles, res)
	})

	t.Run("hairstyle1件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyles/read?entry_id=%d", *f.Entrys[0].ID), nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairStyle
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, hairStyles[0], res[0])
	})

	t.Run("hairstyle2件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyles/read?entry_id=%d&entry_id=%d", *f.Entrys[0].ID, *f.Entrys[1].ID), nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairStyle
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, hairStyles, res)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateHairStyleHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ショートカット"
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
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ウェービーロングヘアー"
		}))),
	)

	// テストデータの準備
	updateHairStyles := []HairStyle{
		{
			EntryID: f.Entrys[0].ID,
			Style:   "ロングヘアー",
		},
		{
			EntryID: f.Entrys[1].ID,
			Style:   "ショートヘアー",
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairstyle更新", func(t *testing.T) {
		h := NewUpdateHandler(indexService)
		body, err := json.Marshal(updateHairStyles)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/hairstyles/update", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairStyle
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, updateHairStyles, res)
	})
}

func TestDeleteHairStyleHandler(t *testing.T) {
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
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ショートカット"
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
		}).Connect(fixtures.NewHairStyle(ctx, func(s *fixtures.HairStyle) {
			s.Style = "ウェービーロングヘアー"
		}))),
	)
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("hairstyle削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{*f.Entrys[0].ID, *f.Entrys[1].ID}}
		h := NewDeleteHandler(indexService)
		body, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/hairstyles/delete", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, delIDs, res)
	})
}
