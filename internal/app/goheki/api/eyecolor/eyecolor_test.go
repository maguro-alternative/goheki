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

	t.Run("eyecolor作成失敗", func(t *testing.T) {
		// リクエストを作成
		ids := []int64{f.Entrys[0].ID, f.Entrys[1].ID}
		b, err := json.Marshal(ids)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/eyecolor/create", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

		var actual []EyeColor
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM eyecolor")
		assert.NoError(t, err)
		assert.Len(t, actual, 0)
	})

	t.Run("eyecolor作成", func(t *testing.T) {
		// リクエストを作成
		eyeColorsJson := EyeColorsJson{
			[]EyeColor{
				{
					EntryID: f.Entrys[0].ID,
					ColorID: f.EyeColorTypes[0].ID,
				},
				{
					EntryID: f.Entrys[1].ID,
					ColorID: f.EyeColorTypes[1].ID,
				},
			},
		}
		b, err := json.Marshal(eyeColorsJson)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/eyecolor/create", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, eyeColorsJson, res)

		var actual []EyeColor
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM eyecolor")
		assert.NoError(t, err)
	})
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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
			s.ColorID = f.EyeColorTypes[0].ID
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
			s.ColorID = f.EyeColorTypes[1].ID
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
		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res.EyeColors[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, res.EyeColors[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res.EyeColors[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, res.EyeColors[1].ColorID)
	})

	t.Run("eyecolor1件取得", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d", f.Entrys[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res.EyeColors[0].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res.EyeColors[0].ColorID)
	})

	t.Run("eyecolor2件取得", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d&entry_id=%d", f.Entrys[0].ID, f.Entrys[1].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, f.EyeColors[0].EntryID, res.EyeColors[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, res.EyeColors[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res.EyeColors[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, res.EyeColors[1].ColorID)
	})

	t.Run("eyecolor2件取得(内1件は存在しない)", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d&entry_id=0", f.Entrys[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)

		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)

		assert.Len(t, res.EyeColors, 1)
		assert.Equal(t, f.EyeColors[0].EntryID, res.EyeColors[0].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, res.EyeColors[0].ColorID)
	})

	t.Run("eyecolor1件取得失敗", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, "/api/eyecolor/read?entry_id=", nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("eyecolor2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor/read?entry_id=%d&entry_id=aaa", f.Entrys[0].ID), nil)
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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
			s.ColorID = f.EyeColorTypes[0].ID
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
			s.ColorID = f.EyeColorTypes[1].ID
		}))),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("eyecolor更新失敗", func(t *testing.T) {
		// リクエストを作成
		eyeColorsJson := IDs{}
		b, err := json.Marshal(eyeColorsJson)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/api/eyecolor/update", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewUpdateHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

		var actual []EyeColor
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM eyecolor")
		assert.NoError(t, err)
		assert.Len(t, actual, 2)

		assert.Equal(t, f.EyeColors[0].EntryID, actual[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, actual[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, actual[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, actual[1].ColorID)
	})

	t.Run("eyecolor更新", func(t *testing.T) {
		updateEyeColorsJson := EyeColorsJson{
			[]EyeColor{
				{
					EntryID: f.Entrys[0].ID,
					ColorID: f.EyeColorTypes[1].ID,
				},
				{
					EntryID: f.Entrys[1].ID,
					ColorID: f.EyeColorTypes[0].ID,
				},
			},
		}
		b, err := json.Marshal(updateEyeColorsJson)
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
		var res EyeColorsJson
		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Equal(t, updateEyeColorsJson, res)

		var actual []EyeColor
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM eyecolor")
		assert.NoError(t, err)
		assert.Equal(t, updateEyeColorsJson.EyeColors[0].EntryID, actual[0].EntryID)
		assert.Equal(t, updateEyeColorsJson.EyeColors[1].EntryID, actual[1].EntryID)
		assert.Equal(t, updateEyeColorsJson.EyeColors[0].ColorID, actual[0].ColorID)
		assert.Equal(t, updateEyeColorsJson.EyeColors[1].ColorID, actual[1].ColorID)
	})
}

func TestDeleteEyeColorHandler(t *testing.T) {
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
			s.ColorID = f.EyeColorTypes[0].ID
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
			s.ColorID = f.EyeColorTypes[1].ID
		}))),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("eyecolor削除失敗", func(t *testing.T) {
		// リクエストを作成
		ids := EyeColorsJson{}
		b, err := json.Marshal(ids)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, "/api/eyecolor/delete", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewDeleteHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

		var actual []EyeColor
		err = tx.SelectContext(ctx, &actual, "SELECT * FROM eyecolor")
		assert.NoError(t, err)
		assert.Len(t, actual, 2)

		assert.Equal(t, f.EyeColors[0].EntryID, actual[0].EntryID)
		assert.Equal(t, f.EyeColors[1].EntryID, actual[1].EntryID)
		assert.Equal(t, f.EyeColors[0].ColorID, actual[0].ColorID)
		assert.Equal(t, f.EyeColors[1].ColorID, actual[1].ColorID)
	})

	t.Run("eyecolor削除", func(t *testing.T) {
		delIDs := IDs{IDs: []int64{f.Entrys[0].ID, f.Entrys[1].ID}}
		eJson, err := json.Marshal(&delIDs)
		assert.NoError(t, err)
		// リクエストを作成
		req, err := http.NewRequest(http.MethodDelete, "/api/eyecolor/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)
		// レスポンスを作成
		rr := httptest.NewRecorder()
		handler := NewDeleteHandler(indexService)
		handler.ServeHTTP(rr, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rr.Code)
		var actual IDs
		err = json.NewDecoder(rr.Body).Decode(&actual)
		assert.NoError(t, err)
		assert.Equal(t, delIDs, actual)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM eyecolor")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
