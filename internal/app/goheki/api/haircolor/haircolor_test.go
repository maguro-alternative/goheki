package haircolor

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

func TestCreateHairColorHandler(t *testing.T) {
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
	hairColors := []HairColor{
		{
			EntryID: *f.Entrys[0].ID,
			Color:   "銀",
		},
		{
			EntryID: *f.Entrys[1].ID,
			Color:   "銀",
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor登録", func(t *testing.T) {
		// リクエストの準備
		handler := NewCreateHandler(indexService)
		body, err := json.Marshal(hairColors)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/haircolor/create", bytes.NewBuffer(body))

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res []HairColor
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, hairColors, res)
	})
}

func TestReadHairColorHandler(t *testing.T) {
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
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.Color = "銀"
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
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.Color = "銀"
		}))),
	)

	// テストデータの準備
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor全件取得", func(t *testing.T) {
		// リクエストの準備
		handler := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/haircolor/read", nil)

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res []HairColor
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})

	t.Run("haircolor1件取得", func(t *testing.T) {
		// リクエストの準備
		handler := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/haircolor/read?entry_id=%d", *f.Entrys[0].ID), nil)

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res []HairColor
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "銀", res[0].Color)
	})

	t.Run("haircolor2件取得", func(t *testing.T) {
		// リクエストの準備
		handler := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/haircolor/read?entry_id=%d&entry_id=%d", *f.Entrys[0].ID, *f.Entrys[1].ID), nil)

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res []HairColor
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "銀", res[0].Color)
		assert.Equal(t, "銀", res[1].Color)
	})
}

func TestUpdateHairColorHandler(t *testing.T) {
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
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.Color = "銀"
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
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.Color = "銀"
		}))),
	)

	// テストデータの準備
	updateHairColors := []HairColor{
		{
			EntryID: *f.Entrys[0].ID,
			Color:   "金",
		},
		{
			EntryID: *f.Entrys[1].ID,
			Color:   "金",
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor更新", func(t *testing.T) {
		// リクエストの準備
		handler := NewUpdateHandler(indexService)
		body, err := json.Marshal(updateHairColors)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/haircolor/update", bytes.NewBuffer(body))

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res []HairColor
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, updateHairColors, res)
	})
}
