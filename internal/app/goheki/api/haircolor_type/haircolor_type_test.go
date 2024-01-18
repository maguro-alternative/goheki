package haircolortype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"

	"github.com/maguro-alternative/goheki/pkg/db"

	"github.com/maguro-alternative/goheki/internal/app/goheki/model/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestCreateHairColorTypeHandler(t *testing.T) {
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

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor_type登録", func(t *testing.T) {
		// リクエストの作成
		hairColorType := []HairColorType{
			{
				Color: "black",
			},
			{
				Color: "blue",
			},
		}
		b, err := json.Marshal(hairColorType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/haircolor_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewCreateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestReadHairColorTypeHandler(t *testing.T) {
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

	// データベースの準備
	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewHairColorType(ctx, func(h *fixtures.HairColorType) {
			h.Color = "black"
		}),
		fixtures.NewHairColorType(ctx, func(h *fixtures.HairColorType) {
			h.Color = "blue"
		}),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("haircolor_type全件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/haircolor_type/read", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)

		var res []HairColorType
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)

		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "black", res[0].Color)
		assert.Equal(t, "blue", res[1].Color)
	})

	t.Run("haircolor_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/haircolor_type/read?id=%d", *f.HairColorTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)

		var res []HairColorType
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)

		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "black", res[0].Color)
	})
}


