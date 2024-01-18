package eyecolor

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

func TestCreateEyeColorHandler(t *testing.T) {
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
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "青"
		}),
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "小豆色"
		}),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("eyecolor作成", func(t *testing.T) {
		// リクエストを作成
		eyeColors := []EyeColor{
			{
				EntryID: *f.Entrys[0].ID,
				ColorID: *f.EyeColorTypes[0].ID,
			},
			{
				EntryID: *f.Entrys[1].ID,
				ColorID: *f.EyeColorTypes[1].ID,
			},
		}
		b, err := json.Marshal(eyeColors)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/eyecolor/create", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res []EyeColor
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, eyeColors, res)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestReadEyeColorHandler(t *testing.T) {
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
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "青"
		}),
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "小豆色"
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
		}).Connect(fixtures.NewEyeColor(ctx, func(s *fixtures.EyeColor) {
			s.ColorID = *f.EyeColorTypes[0].ID
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
		}).Connect(fixtures.NewEyeColor(ctx, func(s *fixtures.EyeColor) {
			s.ColorID = *f.EyeColorTypes[1].ID
		}))),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("eyecolor全件取得", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, "/api/eyecolor/read", nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res []EyeColor
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, res[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, res[1].ColorID)
	})

	t.Run("eyecolor1件取得", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d", *f.Entrys[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res []EyeColor
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res[0].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res[0].ColorID)
	})

	t.Run("eyecolor2件取得", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d&entry_id=%d", *f.Entrys[0].ID, *f.Entrys[1].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res []EyeColor
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, res[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, res[1].ColorID)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateEyeColorHandler(t *testing.T) {
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
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "青"
		}),
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "小豆色"
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
		}).Connect(fixtures.NewEyeColor(ctx, func(s *fixtures.EyeColor) {
			s.ColorID = *f.EyeColorTypes[0].ID
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
		}).Connect(fixtures.NewEyeColor(ctx, func(s *fixtures.EyeColor) {
			s.ColorID = *f.EyeColorTypes[1].ID
		}))),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("eyecolor更新", func(t *testing.T) {
		updateEyeColors := []EyeColor{
			{
				EntryID: *f.Entrys[0].ID,
				ColorID: *f.EyeColorTypes[1].ID,
			},
			{
				EntryID: *f.Entrys[1].ID,
				ColorID: *f.EyeColorTypes[0].ID,
			},
		}
		b, err := json.Marshal(updateEyeColors)
		assert.NoError(t, err)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodPut, "/api/eyecolor/update", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewUpdateHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res []EyeColor
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, updateEyeColors, res)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}
