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

	// ロールバック
	defer tx.RollbackCtx(ctx)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("eyecolor_type登録失敗", func(t *testing.T) {
		// リクエストの作成
		ids := IDs{
			IDs: []int64{1, 2},
		}
		b, err := json.Marshal(ids)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/eyecolor_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewCreateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actuals []EyeColorType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM eyecolor_type")
		assert.NoError(t, err)
		assert.Len(t, actuals, 0)
	})

	t.Run("eyecolor_type登録", func(t *testing.T) {
		// リクエストの作成
		eyeColorTypesJson := EyeColorTypesJson{
			[]EyeColorType{
				{
					Color: "black",
				},
				{
					Color: "blue",
				},
			},
		}
		b, err := json.Marshal(eyeColorTypesJson)
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
		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, eyeColorTypesJson, res)

		var actualEyeColorTypes []EyeColorType
		err = tx.SelectContext(ctx, &actualEyeColorTypes, "SELECT * FROM eyecolor_type")
		assert.NoError(t, err)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[0].Color, actualEyeColorTypes[0].Color)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[1].Color, actualEyeColorTypes[1].Color)
	})
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
	// ロールバック
	defer tx.RollbackCtx(ctx)
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
		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, f.EyeColorTypes[0].ID, res.EyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res.EyeColorTypes[0].Color)
		assert.Equal(t, f.EyeColorTypes[1].ID, res.EyeColorTypes[1].ID)
		assert.Equal(t, f.EyeColorTypes[1].Color, res.EyeColorTypes[1].Color)
	})

	t.Run("eyecolor_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor_type/read?id=%d", f.EyeColorTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, f.EyeColorTypes[0].ID, res.EyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res.EyeColorTypes[0].Color)
	})

	t.Run("eyecolor_type2件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor_type/read?id=%d&id=%d", f.EyeColorTypes[0].ID, f.EyeColorTypes[1].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, f.EyeColorTypes[0].ID, res.EyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res.EyeColorTypes[0].Color)
		assert.Equal(t, f.EyeColorTypes[1].ID, res.EyeColorTypes[1].ID)
		assert.Equal(t, f.EyeColorTypes[1].Color, res.EyeColorTypes[1].Color)
	})

	t.Run("eyecolor_type1件取得(存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/eyecolor_type/read?id=0", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Len(t, res.EyeColorTypes, 0)
	})

	t.Run("eyecolor_type2件取得(内1件は存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor_type/read?id=%d&id=0", f.EyeColorTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		assert.Len(t, res.EyeColorTypes, 1)
		assert.Equal(t, f.EyeColorTypes[0].ID, res.EyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, res.EyeColorTypes[0].Color)
	})

	t.Run("eyecolor_type1件取得(バリデーションエラー)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/eyecolor_type/read?id=invalid", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("eyecolor_type2件取得(内1件は形式が正しくない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/eyecolor_type/read?id=%d&id=invalid", f.EyeColorTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行

		h := NewReadHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateEyeColorTypeHandler(t *testing.T) {
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

	t.Run("eyecolor_type更新失敗", func(t *testing.T) {
		// リクエストの作成
		eyeColorTypesJson := EyeColorTypesJson{}
		b, err := json.Marshal(eyeColorTypesJson)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/eyecolor_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()

		// テスト対象のハンドラを実行
		h := NewUpdateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actualEyeColorTypes []EyeColorType
		err = tx.SelectContext(ctx, &actualEyeColorTypes, "SELECT * FROM eyecolor_type")
		assert.NoError(t, err)

		assert.Equal(t, f.EyeColorTypes[0].ID, actualEyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, actualEyeColorTypes[0].Color)
		assert.Equal(t, f.EyeColorTypes[1].ID, actualEyeColorTypes[1].ID)
		assert.Equal(t, f.EyeColorTypes[1].Color, actualEyeColorTypes[1].Color)
	})

	t.Run("eyecolor_type更新", func(t *testing.T) {
		// リクエストの作成
		eyeColorTypesJson := EyeColorTypesJson{
			[]EyeColorType{
				{
					ID:    f.EyeColorTypes[0].ID,
					Color: "red",
				},
				{
					ID:    f.EyeColorTypes[1].ID,
					Color: "green",
				},
			},
		}
		b, err := json.Marshal(eyeColorTypesJson)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/api/eyecolor_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewUpdateHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res EyeColorTypesJson
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, eyeColorTypesJson, res)

		var actualEyeColorTypes []EyeColorType
		err = tx.SelectContext(ctx, &actualEyeColorTypes, "SELECT * FROM eyecolor_type")
		assert.NoError(t, err)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[0].ID, actualEyeColorTypes[0].ID)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[0].Color, actualEyeColorTypes[0].Color)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[1].ID, actualEyeColorTypes[1].ID)
		assert.Equal(t, eyeColorTypesJson.EyeColorTypes[1].Color, actualEyeColorTypes[1].Color)
	})
}

func TestDeleteEyeColorTypeHandler(t *testing.T) {
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

	t.Run("eyecolor_type削除失敗", func(t *testing.T) {
		// リクエストの作成
		ids := IDs{}
		b, err := json.Marshal(ids)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/eyecolor_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()

		// テスト対象のハンドラを実行
		h := NewDeleteHandler(indexService)
		h.ServeHTTP(w, req)

		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actualEyeColorTypes []EyeColorType
		err = tx.SelectContext(ctx, &actualEyeColorTypes, "SELECT * FROM eyecolor_type")
		assert.NoError(t, err)

		assert.Equal(t, f.EyeColorTypes[0].ID, actualEyeColorTypes[0].ID)
		assert.Equal(t, f.EyeColorTypes[0].Color, actualEyeColorTypes[0].Color)
		assert.Equal(t, f.EyeColorTypes[1].ID, actualEyeColorTypes[1].ID)
		assert.Equal(t, f.EyeColorTypes[1].Color, actualEyeColorTypes[1].Color)
	})

	t.Run("eyecolor_type削除", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{
			IDs: []int64{f.EyeColorTypes[0].ID, f.EyeColorTypes[1].ID},
		}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/eyecolor_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		h := NewDeleteHandler(indexService)
		h.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res IDs
		err = json.NewDecoder(w.Body).Decode(&res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, delIDs, res)
	})
}
