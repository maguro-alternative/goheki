package hairstyletype

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

func TestCreateHairStyleTypeHandler(t *testing.T) {
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
	t.Run("hairstyle_type登録", func(t *testing.T) {
		// リクエストの作成
		hairStyleType := []HairStyleType{
			{
				Style: "short",
			},
			{
				Style: "long",
			},
		}
		b, err := json.Marshal(hairStyleType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/hairstyle_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewCreateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypes []HairStyleType
		err = json.Unmarshal(w.Body.Bytes(), &hairStyleTypes)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypes))
		assert.Equal(t, "short", hairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypes[1].Style)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestReadHairStyleTypeHandler(t *testing.T) {
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
		fixtures.NewHairStyleType(ctx, func(s *fixtures.HairStyleType) {
			s.Style = "short"
		}),
		fixtures.NewHairStyleType(ctx, func(s *fixtures.HairStyleType) {
			s.Style = "long"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairstyle_type全件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/hairstyle_type/read", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypes []HairStyleType
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypes)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypes))
		assert.Equal(t, "short", hairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypes[1].Style)
	})

	t.Run("hairstyle_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d", *f.HairStyleTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleType []HairStyleType
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleType)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(hairStyleType))
		assert.Equal(t, "short", hairStyleType[0].Style)
	})

	t.Run("hairstyle_type2件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d&id=%d", *f.HairStyleTypes[0].ID, *f.HairStyleTypes[1].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypes []HairStyleType
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypes)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypes))
		assert.Equal(t, "short", hairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypes[1].Style)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateHairStyleTypeHandler(t *testing.T) {
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
		fixtures.NewHairStyleType(ctx, func(s *fixtures.HairStyleType) {
			s.Style = "short"
		}),
		fixtures.NewHairStyleType(ctx, func(s *fixtures.HairStyleType) {
			s.Style = "long"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairstyle_type更新", func(t *testing.T) {
		// リクエストの作成
		updateHairStyleType := []HairStyleType{
			{
				ID:    f.HairStyleTypes[0].ID,
				Style: "short",
			},
			{
				ID:    f.HairStyleTypes[1].ID,
				Style: "long",
			},
		}
		b, err := json.Marshal(updateHairStyleType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/hairstyle_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewUpdateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypes []HairStyleType
		err = json.Unmarshal(w.Body.Bytes(), &hairStyleTypes)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypes))
		assert.Equal(t, "short", hairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypes[1].Style)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}
