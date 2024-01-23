package hairlengthtype

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

func TestCreateHairLengthTypeHandler(t *testing.T) {
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
	t.Run("hairlength_type登録", func(t *testing.T) {
		// リクエストの作成
		hairLengthType := HairLengthTypesJson{[]HairLengthType{
			{
				Length: "short",
			},
			{
				Length: "long",
			},
		}}
		b, err := json.Marshal(hairLengthType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/hairlength_type/create", bytes.NewBuffer(b))
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

func TestReadHairLengthTypeHandler(t *testing.T) {
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

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "long"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength_type全件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/hairlength_type/read", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)

		var hairLengthTypesJson HairLengthTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairLengthTypesJson)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 2, len(hairLengthTypesJson.HairLengthTypes))
		assert.Equal(t, "short", hairLengthTypesJson.HairLengthTypes[0].Length)
		assert.Equal(t, "long", hairLengthTypesJson.HairLengthTypes[1].Length)
	})

	t.Run("hairlength_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairlength_type/read?id=%d", f.HairLengthTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)

		var hairLengthTypesJson HairLengthTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairLengthTypesJson)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "short", hairLengthTypesJson.HairLengthTypes[0].Length)
	})

	t.Run("hairlength_type2件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairlength_type/read?id=%d&id=%d", f.HairLengthTypes[0].ID, f.HairLengthTypes[1].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)

		var hairLengthTypesJson HairLengthTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairLengthTypesJson)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, 2, len(hairLengthTypesJson.HairLengthTypes))
		assert.Equal(t, "short", hairLengthTypesJson.HairLengthTypes[0].Length)
		assert.Equal(t, "long", hairLengthTypesJson.HairLengthTypes[1].Length)
	})

	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestUpdateHairLengthTypeHandler(t *testing.T) {
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

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "long"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength_type更新", func(t *testing.T) {
		// リクエストの作成
		updateHairLengthType := HairLengthTypesJson{[]HairLengthType{
			{
				ID:     f.HairLengthTypes[0].ID,
				Length: "short",
			},
			{
				ID:     f.HairLengthTypes[1].ID,
				Length: "long",
			},
		}}
		b, err := json.Marshal(updateHairLengthType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/hairlength_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewUpdateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairLengthTypesJson HairLengthTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &hairLengthTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairLengthTypesJson.HairLengthTypes))
		assert.Equal(t, "short", hairLengthTypesJson.HairLengthTypes[0].Length)
		assert.Equal(t, "long", hairLengthTypesJson.HairLengthTypes[1].Length)
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}

func TestDeleteHairLengthTypeHandler(t *testing.T) {
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

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "short"
		}),
		fixtures.NewHairLengthType(ctx, func(h *fixtures.HairLengthType) {
			h.Length = "long"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("hairlength_type削除", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{
			IDs: []int64{f.HairLengthTypes[0].ID, f.HairLengthTypes[1].ID},
		}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/hairlength_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewDeleteHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res.IDs))
		assert.Equal(t, f.HairLengthTypes[0].ID, res.IDs[0])
		assert.Equal(t, f.HairLengthTypes[1].ID, res.IDs[1])
	})
	// ロールバック
	tx.RollbackCtx(ctx)
}
