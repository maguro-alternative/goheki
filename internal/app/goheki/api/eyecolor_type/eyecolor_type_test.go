package eyecolortype

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

func TestCreateEyeColorTypeHandler(t *testing.T) {
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
	t.Run("eyecolor_type登録", func(t *testing.T) {
		// リクエストの作成
		eyeColorType := []EyeColorType{
			{
				Color: "black",
			},
			{
				Color: "blue",
			},
		}
		b, err := json.Marshal(eyeColorType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/eyecolor_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewCreateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res []EyeColorType
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, eyeColorType, res)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestReadEyeColorTypeHandler(t *testing.T) {
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
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "black"
		}),
		fixtures.NewEyeColorType(ctx, func(s *fixtures.EyeColorType) {
			s.Color = "blue"
		}),
	)
	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("eyecolor_type全件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/eyecolor_type/read", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res []EyeColorType
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, f.EyeColorTypes[0].ID, res[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res[0].Color)
		assert.Equal(t, f.EyeColorTypes[1].ID, res[1].ID)
		assert.Equal(t, f.EyeColorTypes[1].Color, res[1].Color)
	})

	t.Run("eyecolor_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor_type/read?id=%d", *f.EyeColorTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res []EyeColorType
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, f.EyeColorTypes[0].ID, res[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res[0].Color)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}
