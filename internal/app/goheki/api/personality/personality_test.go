package personality

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

func TestCreatePersonalityHandler(t *testing.T) {
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
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("personality追加", func(t *testing.T) {
		personalitys := []Personality{
			{
				EntryID: f.Entrys[0].ID,
				Type:    "jun",
			},
			{
				EntryID: f.Entrys[1].ID,
				Type:    "ten",
			},
		}
		// テストの実行
		h := NewCreateHandler(indexService)
		pJson, err := json.Marshal(personalitys)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/personality/create", bytes.NewBuffer(pJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの確認
		var res []Personality
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, personalitys, res)
	})
}

func TestReadPersonalityHandler(t *testing.T) {
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
		}).Connect(fixtures.NewPersonality(ctx, func(s *fixtures.Personality) {
			s.Type = "jun"
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
		}).Connect(fixtures.NewPersonality(ctx, func(s *fixtures.Personality) {
			s.Type = "ten"
		}))),
	)
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("personality全件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, "/api/personality/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの確認
		var res []Personality
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, f.Personalities, res)
	})

	t.Run("personality1件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/personality/read?id=%d", *f.Personalities[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの確認
		var res []Personality
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, f.Personalities[0].Type, res[0].Type)
	})

	t.Run("personality2件取得", func(t *testing.T) {
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/personality/read?id=%d&id=%d", *f.Personalities[0].ID, *f.Personalities[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの確認
		var res []Personality
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, f.Personalities[0].Type, res[0].Type)
		assert.Equal(t, f.Personalities[1].Type, res[1].Type)
	})
}

func TestUpdatePersonalityHandler(t *testing.T) {
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
		}).Connect(fixtures.NewPersonality(ctx, func(s *fixtures.Personality) {
			s.Type = "jun"
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
		}).Connect(fixtures.NewPersonality(ctx, func(s *fixtures.Personality) {
			s.Type = "ten"
		}))),
	)
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("personality更新", func(t *testing.T) {
		personalitys := []Personality{
			{
				ID:      f.Personalities[0].ID,
				EntryID: f.Entrys[0].ID,
				Type:    "tun",
			},
			{
				ID:      f.Personalities[1].ID,
				EntryID: f.Entrys[1].ID,
				Type:    "dere",
			},
		}
		// テストの実行
		h := NewUpdateHandler(indexService)
		pJson, err := json.Marshal(personalitys)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, "/api/personality/update", bytes.NewBuffer(pJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスの確認
		var res []Personality
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, personalitys[0].Type, res[0].Type)
	})
}
