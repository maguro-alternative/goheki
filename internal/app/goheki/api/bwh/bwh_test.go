package bwh

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

func TestCreateBEHHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	yumiHeight := int64(167)
	takaneHeight := int64(169)
	takaneWeight := int64(49)
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
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)

	// テストデータの準備
	bwhs := []BWH{
		{
			EntryID: f.Entrys[0].ID,
			Bust:    92,
			Waist:   56,
			Hip:     84,
			Height:  &yumiHeight,
		},
		{
			EntryID: f.Entrys[1].ID,
			Bust:    90,
			Waist:   60,
			Hip:     92,
			Height:  &takaneHeight,
			Weight:  &takaneWeight,
		},
	}
	bwhJsons := BWHJsons{
		BWHs: bwhs,
	}
	ids := []int64{*f.Entrys[0].ID, *f.Entrys[1].ID}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("bwh登録", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		bJson, err := json.Marshal(&bwhJsons)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/bwh/create", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHJsons
		var dbResult []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res.BWHs)
		err = tx.SelectContext(ctx, &dbResult, "SELECT * FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, bwhJsons.BWHs, dbResult)
	})

	t.Run("bwh登録失敗", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		bJson, err := json.Marshal(ids)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/bwh/create", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestReadBEHHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	yumiHeight := int64(167)
	takaneHeight := int64(169)
	takaneWeight := int64(49)
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 92
			s.Waist = 56
			s.Hip = 84
			s.Height = &yumiHeight
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 90
			s.Waist = 60
			s.Hip = 92
			s.Height = &takaneHeight
			s.Weight = &takaneWeight
		}))),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)

	// テストデータの準備
	bwhs := []BWH{
		{
			EntryID: f.Entrys[0].ID,
			Bust:    92,
			Waist:   56,
			Hip:     84,
			Height:  &yumiHeight,
		},
		{
			EntryID: f.Entrys[1].ID,
			Bust:    90,
			Waist:   60,
			Hip:     92,
			Height:  &takaneHeight,
			Weight:  &takaneWeight,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("bwh全件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/bwh/read", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res)
	})

	t.Run("bwh1件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=%d", *f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs[:1], res)
	})

	t.Run("bwh2件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=%d&entry_id=%d", *f.Entrys[0].ID, *f.Entrys[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateBEHHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	yumiHeight := int64(167)
	takaneHeight := int64(169)
	takaneWeight := int64(49)
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 92
			s.Waist = 56
			s.Hip = 84
			s.Height = &yumiHeight
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 90
			s.Waist = 60
			s.Hip = 92
			s.Height = &takaneHeight
			s.Weight = &takaneWeight
		}))),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)

	// テストデータの準備
	updateBWHs := []BWH{
		{
			EntryID: f.Entrys[0].ID,
			Bust:    98,
			Waist:   59,
			Hip:     87,
			Height:  &yumiHeight,
			Weight:  nil,
		},
		{
			EntryID: f.Entrys[1].ID,
			Bust:    93,
			Waist:   62,
			Hip:     96,
			Height:  &takaneHeight,
			Weight:  &takaneWeight,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("bwh更新", func(t *testing.T) {
		h := NewUpdateHandler(indexService)
		bJson, err := json.Marshal(updateBWHs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/bwh/update", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, updateBWHs, res)
	})
}

func TestDeleteBEHHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	yumiHeight := int64(167)
	takaneHeight := int64(169)
	takaneWeight := int64(49)
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 92
			s.Waist = 56
			s.Hip = 84
			s.Height = &yumiHeight
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
		}).Connect(fixtures.NewBWH(ctx, func(s *fixtures.BWH) {
			s.Bust = 90
			s.Waist = 60
			s.Hip = 92
			s.Height = &takaneHeight
			s.Weight = &takaneWeight
		}))),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("bwh削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{*f.BWHs[0].EntryID}}
		h := NewDeleteHandler(indexService)
		bJson, err := json.Marshal(delIDs)
		req := httptest.NewRequest(http.MethodDelete, "/api/bwh/delete", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

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
