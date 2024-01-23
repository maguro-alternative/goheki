package hairlength

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

func TestCreateHairLengthHandler(t *testing.T) {
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
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "long"
		}),
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
	heirLengths := []HairLength{
		{
			EntryID:          f.Entrys[0].ID,
			HairLengthTypeID: f.HairLengthTypes[0].ID,
		},
		{
			EntryID:          f.Entrys[1].ID,
			HairLengthTypeID: f.HairLengthTypes[1].ID,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength登録", func(t *testing.T) {
		h := NewCreateHandler(indexService)
		body, err := json.Marshal(heirLengths)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/hairlength", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairLength
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, heirLengths, res)
	})
}

func TestReadHairLengthHandler(t *testing.T) {
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
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "long"
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[0].ID
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
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[1].ID
		}))),
	)

	// テストデータの準備
	heirLengths := []HairLength{
		{
			EntryID:          f.Entrys[0].ID,
			HairLengthTypeID: f.HairLengthTypes[0].ID,
		},
		{
			EntryID:          f.Entrys[1].ID,
			HairLengthTypeID: f.HairLengthTypes[1].ID,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength全件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, "/api/hairlength/read", nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairLength
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, heirLengths, res)
	})

	t.Run("hairlength1件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairlength/read?entry_id=%d", f.Entrys[0].ID), nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairLength
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, heirLengths[:1], res)
	})

	t.Run("hairlength2件取得", func(t *testing.T) {
		h := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairlength/read?entry_id=%d&entry_id=%d", f.Entrys[0].ID, f.Entrys[1].ID), nil)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		// tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairLength
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, heirLengths, res)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateHairLengthHandler(t *testing.T) {
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
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "long"
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[0].ID
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
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[1].ID
		}))),
	)

	// テストデータの準備
	updateHeirLengths := []HairLength{
		{
			EntryID:          f.Entrys[0].ID,
			HairLengthTypeID: f.HairLengthTypes[1].ID,
		},
		{
			EntryID:          f.Entrys[1].ID,
			HairLengthTypeID: f.HairLengthTypes[0].ID,
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength更新", func(t *testing.T) {
		h := NewUpdateHandler(indexService)
		body, err := json.Marshal(updateHeirLengths)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/hairlength/update", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res []HairLength
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, updateHeirLengths, res)
	})
}

func TestDeleteHairLengthHandler(t *testing.T) {
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
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(s *fixtures.HairLengthType) {
			s.Length = "long"
		}),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "閃乱カグラ"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "雪泉"
			s.Image = "https://example.com/image1.png"
			s.Content = "かわいい"
			s.CreatedAt = fixedTime
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[0].ID
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
		}).Connect(fixtures.NewHairLength(ctx, func(s *fixtures.HairLength) {
			s.HairLengthTypeID = f.HairLengthTypes[1].ID
		}))),
	)

	// テストデータの準備
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength削除", func(t *testing.T) {
		h := NewDeleteHandler(indexService)
		body, err := json.Marshal(IDs{IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID}})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/hairlength/delete", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		tx.RollbackCtx(ctx)
		assert.Equal(t, http.StatusOK, w.Code)

		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, []int64{f.Entrys[0].ID, f.Entrys[1].ID}, res.IDs)
	})
}
