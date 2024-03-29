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
	bwhsJson := BWHsJson{
		BWHs: bwhs,
	}
	ids := IDs{IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID}}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("bwh登録失敗", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		bJson, err := json.Marshal(ids)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/bwh/create", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("bwh登録失敗 形式が正しくないjson", func(t *testing.T) {
		type dummey struct {
			EntryID int64 `json:"entry_id"`
			ID      int64 `json:"id"`
		}
		type dummeyJson struct {
			BWHs []dummey `json:"bwhs"`
		}
		h := NewCreateHandler(indexService)
		bJson, err := json.Marshal(dummeyJson{
			BWHs: []dummey{
				{
					EntryID: 0,
					ID:      0,
				},
			},
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/bwh/create", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("bwh登録", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		bJson, err := json.Marshal(&bwhsJson)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/bwh/create", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHsJson
		var dbResult []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res.BWHs)
		err = tx.SelectContext(ctx, &dbResult, "SELECT * FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, bwhsJson.BWHs, dbResult)
	})
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

		var res BWHsJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res.BWHs)
	})

	t.Run("bwh1件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=%d", f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHsJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs[:1], res.BWHs)
	})

	t.Run("bwh2件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=%d&entry_id=%d", f.Entrys[0].ID, f.Entrys[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHsJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs, res.BWHs)
	})

	t.Run("bwh1件取得(存在しない)", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/bwh/read?entry_id=0", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)
		var res BWHsJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Len(t, res.BWHs, 0)
	})

	t.Run("bwh2件取得(内1件は存在しない)", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=%d&entry_id=0", f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHsJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, bwhs[:1], res.BWHs)
	})

	t.Run("bwh1件取得(形式が正しくない)", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/bwh/read?entry_id=aaa", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("bwh2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/bwh/read?entry_id=aaa&entry_id=%d", f.BWHs[0].EntryID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
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
	updateBWHsJson := BWHsJson{
		[]BWH{
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
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("bwh更新(形式が正しくない)", func(t *testing.T) {
		type dummey struct {
			EntryID int64 `json:"entry_id"`
			ID      int64 `json:"id"`
		}
		type dummeyJson struct {
			BWHs []dummey `json:"bwhs"`
		}
		h := NewUpdateHandler(indexService)
		bJson, err := json.Marshal(dummeyJson{
			BWHs: []dummey{
				{
					EntryID: 0,
					ID:      0,
				},
			},
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/bwh/update", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actual []BWH
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, f.BWHs[0].Bust, actual[0].Bust)
		assert.Equal(t, f.BWHs[0].Waist, actual[0].Waist)
		assert.Equal(t, f.BWHs[0].Hip, actual[0].Hip)
		assert.Equal(t, f.BWHs[0].Height, actual[0].Height)
		assert.Equal(t, f.BWHs[0].Weight, actual[0].Weight)
		assert.Equal(t, f.BWHs[1].Bust, actual[1].Bust)
		assert.Equal(t, f.BWHs[1].Waist, actual[1].Waist)
		assert.Equal(t, f.BWHs[1].Hip, actual[1].Hip)
		assert.Equal(t, f.BWHs[1].Height, actual[1].Height)
		assert.Equal(t, f.BWHs[1].Weight, actual[1].Weight)
	})

	t.Run("bwh更新", func(t *testing.T) {
		h := NewUpdateHandler(indexService)
		bJson, err := json.Marshal(updateBWHsJson)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/bwh/update", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res BWHsJson
		var actual []BWH
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, updateBWHsJson, res)

		err = tx.SelectContext(ctx, &actual, "SELECT * FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, updateBWHsJson.BWHs, actual)
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

	t.Run("bwh削除(形式が正しくない)", func(t *testing.T) {
		type dummey struct {
			EntryID int64 `json:"entry_id"`
			ID      int64 `json:"id"`
		}
		type dummeyJson struct {
			BWHs []dummey `json:"bwhs"`
		}
		h := NewDeleteHandler(indexService)
		bJson, err := json.Marshal(dummeyJson{
			BWHs: []dummey{
				{
					EntryID: 0,
					ID:      0,
				},
			},
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/bwh/delete", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("bwh削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{f.BWHs[0].EntryID}}
		h := NewDeleteHandler(indexService)
		bJson, err := json.Marshal(delIDs)
		req := httptest.NewRequest(http.MethodDelete, "/api/bwh/delete", bytes.NewBuffer(bJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, delIDs, res)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM bwh")
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
