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
		fixtures.NewHairColorType(ctx, func(s *fixtures.HairColorType) {
			s.Color = "銀"
		}),
	)

	// テストデータの準備
	hairColorsJson := HairColorsJson{[]HairColor{
			{
				EntryID: f.Entrys[0].ID,
				ColorID: f.HairColorTypes[0].ID,
			},
			{
				EntryID: f.Entrys[1].ID,
				ColorID: f.HairColorTypes[0].ID,
			},
		},
	}
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("haircolor登録失敗", func(t *testing.T) {
		// リクエストの準備
		handler := NewCreateHandler(indexService)
		ids := IDs{IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID}}
		body, err := json.Marshal(ids)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/haircolor/create", bytes.NewBuffer(body))

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// データベースの検証
		var hairColors []HairColor
		err = tx.SelectContext(ctx, &hairColors, "SELECT * FROM haircolor")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(hairColors))
	})

	t.Run("haircolor登録", func(t *testing.T) {
		// リクエストの準備
		handler := NewCreateHandler(indexService)
		body, err := json.Marshal(hairColorsJson)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/haircolor/create", bytes.NewBuffer(body))

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)

		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res HairColorsJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, hairColorsJson, res)

		// データベースの検証
		var hairColors []HairColor
		err = tx.SelectContext(ctx, &hairColors, "SELECT * FROM haircolor")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairColors))
		assert.Equal(t, f.Entrys[0].ID, hairColors[0].EntryID)
		assert.Equal(t, f.Entrys[1].ID, hairColors[1].EntryID)
		assert.Equal(t, f.HairColorTypes[0].ID, hairColors[0].ColorID)
		assert.Equal(t, f.HairColorTypes[0].ID, hairColors[1].ColorID)
	})

	// トランザクションのロールバック
	tx.RollbackCtx(ctx)
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
		fixtures.NewHairColorType(ctx, func(s *fixtures.HairColorType) {
			s.Color = "銀"
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[0].ID
		})),
		fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[1].ID
			s.ColorID = f.HairColorTypes[0].ID
		}),
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
		// トランザクションのロールバック
		// tx.RollbackCtx(ctx)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res HairColorsJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res.HairColors))
		assert.Equal(t, f.HairColors[0].ColorID, res.HairColors[0].ColorID)
		assert.Equal(t, f.HairColors[1].ColorID, res.HairColors[1].ColorID)
	})

	t.Run("haircolor1件取得", func(t *testing.T) {
		// リクエストの準備
		handler := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/haircolor/read?entry_id=%d", f.Entrys[0].ID), nil)

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// トランザクションのロールバック
		// tx.RollbackCtx(ctx)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res HairColorsJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res.HairColors))
		assert.Equal(t, f.HairColors[0].ColorID, res.HairColors[0].ColorID)
	})

	t.Run("haircolor2件取得", func(t *testing.T) {
		// リクエストの準備
		handler := NewReadHandler(indexService)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/haircolor/read?entry_id=%d&entry_id=%d", f.Entrys[0].ID, f.Entrys[1].ID), nil)

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// トランザクションのロールバック
		// tx.RollbackCtx(ctx)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res HairColorsJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res.HairColors))
		assert.Equal(t, f.HairColors[0].ColorID, res.HairColors[0].ColorID)
		assert.Equal(t, f.HairColors[1].ColorID, res.HairColors[1].ColorID)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
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
		fixtures.NewHairColorType(ctx, func(s *fixtures.HairColorType) {
			s.Color = "銀"
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[0].ID
		})),
		fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[1].ID
			s.ColorID = f.HairColorTypes[0].ID
		}),
		fixtures.NewHairColorType(ctx, func(s *fixtures.HairColorType) {
			s.Color = "金"
		}),
	)

	// テストデータの準備
	updateHairColors := HairColorsJson{[]HairColor{
		{
			EntryID: f.Entrys[0].ID,
			ColorID: f.HairColorTypes[1].ID,
		},
		{
			EntryID: f.Entrys[1].ID,
			ColorID: f.HairColorTypes[1].ID,
		},
	}}
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
		// トランザクションのロールバック
		tx.RollbackCtx(ctx)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res HairColorsJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, updateHairColors, res)
	})
}

func TestDeleteHairColorHandler(t *testing.T) {
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
		fixtures.NewHairColorType(ctx, func(s *fixtures.HairColorType) {
			s.Color = "銀"
		}).Connect(fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[0].ID
		})),
		fixtures.NewHairColor(ctx, func(s *fixtures.HairColor) {
			s.EntryID = f.Entrys[1].ID
			s.ColorID = f.HairColorTypes[0].ID
		}),
	)

	// テストデータの準備
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor削除", func(t *testing.T) {
		// リクエストの準備
		handler := NewDeleteHandler(indexService)
		body, err := json.Marshal(IDs{IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID}})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/haircolor/delete", bytes.NewBuffer(body))

		// レスポンスの準備
		w := httptest.NewRecorder()
		// ハンドラの実行
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		var res IDs
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, []int64{f.Entrys[0].ID, f.Entrys[1].ID}, res.IDs)
	})
	// トランザクションのロールバック
	tx.RollbackCtx(ctx)
}
