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

	// ロールバック
	defer tx.RollbackCtx(ctx)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("hairstyle_type登録失敗", func(t *testing.T) {
		// リクエストの作成
		hairStyleType := HairStyleTypesJson{[]HairStyleType{}}
		b, err := json.Marshal(hairStyleType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/hairstyle_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewCreateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM hairstyle_type")
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("hairstyle_type登録", func(t *testing.T) {
		// リクエストの作成
		hairStyleType := HairStyleTypesJson{[]HairStyleType{
			{
				Style: "short",
			},
			{
				Style: "long",
			},
		}}
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

		var hairStyleTypesJson HairStyleTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "short", hairStyleTypesJson.HairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypesJson.HairStyleTypes[1].Style)
	})
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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

		var hairStyleTypesJson HairStyleTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "short", hairStyleTypesJson.HairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypesJson.HairStyleTypes[1].Style)
	})

	t.Run("hairstyle_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d", f.HairStyleTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypesJson HairStyleTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "short", hairStyleTypesJson.HairStyleTypes[0].Style)
	})

	t.Run("hairstyle_type2件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d&id=%d", f.HairStyleTypes[0].ID, f.HairStyleTypes[1].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypesJson HairStyleTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "short", hairStyleTypesJson.HairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypesJson.HairStyleTypes[1].Style)
	})

	t.Run("hairstyle_type_type1件取得(存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d", 0), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypesJson HairStyleTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(hairStyleTypesJson.HairStyleTypes))
	})

	t.Run("hairstyle_type2件取得(内1件は存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d&id=%d", f.HairStyleTypes[0].ID, 0), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var hairStyleTypesJson HairStyleTypesJson
		err := json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "short", hairStyleTypesJson.HairStyleTypes[0].Style)
	})

	t.Run("hairstyle_type2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%d&id=%s", f.HairStyleTypes[0].ID, "a"), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("hairstyle_type1件取得(形式が正しくない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hairstyle_type/read?id=%s", "a"), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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

	t.Run("hairstyle_type更新失敗", func(t *testing.T) {
		// リクエストの作成
		updateHairStyleType := HairStyleTypesJson{[]HairStyleType{}}
		b, err := json.Marshal(updateHairStyleType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/hairstyle_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewUpdateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var hairStyleType []HairStyleType
		err = tx.SelectContext(ctx, &hairStyleType, "SELECT * FROM hairstyle_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleType))
		assert.Equal(t, "short", hairStyleType[0].Style)
		assert.Equal(t, "long", hairStyleType[1].Style)
	})

	t.Run("hairstyle_type更新", func(t *testing.T) {
		// リクエストの作成
		updateHairStyleType := HairStyleTypesJson{[]HairStyleType{
			{
				ID:    f.HairStyleTypes[0].ID,
				Style: "medium",
			},
			{
				ID:    f.HairStyleTypes[1].ID,
				Style: "long",
			},
		}}
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

		var hairStyleTypesJson HairStyleTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &hairStyleTypesJson)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hairStyleTypesJson.HairStyleTypes))
		assert.Equal(t, "medium", hairStyleTypesJson.HairStyleTypes[0].Style)
		assert.Equal(t, "long", hairStyleTypesJson.HairStyleTypes[1].Style)
	})
}

func TestDeleteHairStyleTypeHandler(t *testing.T) {
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

	t.Run("hairstyle_type削除失敗", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{
			IDs: []int64{},
		}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/hairstyle_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewDeleteHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var count int
		err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM hairstyle_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("hairstyle_type削除", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{
			IDs: []int64{f.HairStyleTypes[0].ID, f.HairStyleTypes[1].ID},
		}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/hairstyle_type/delete", bytes.NewBuffer(b))
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
		assert.Equal(t, f.HairStyleTypes[0].ID, res.IDs[0])
		assert.Equal(t, f.HairStyleTypes[1].ID, res.IDs[1])
	})
}
