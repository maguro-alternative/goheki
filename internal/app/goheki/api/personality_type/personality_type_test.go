package personalitytype

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

func TestCreatePersonalityTypeHandler(t *testing.T) {
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

	t.Run("personality_type登録失敗", func(t *testing.T) {
		// リクエストの作成
		personalityType := PersonalityTypesJson{[]PersonalityType{}}
		b, err := json.Marshal(personalityType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/personality_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(actuals))
	})

	t.Run("personality_type登録", func(t *testing.T) {
		// リクエストの作成
		personalityType := PersonalityTypesJson{[]PersonalityType{
			{
				Type: "大和撫子",
			},
			{
				Type: "天然",
			},
		}}
		b, err := json.Marshal(personalityType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/personality_type/create", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, personalityType, res)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
		assert.Equal(t, "大和撫子", actuals[0].Type)
		assert.Equal(t, "天然", actuals[1].Type)
	})
}

func TestReadPersonalityTypeHandler(t *testing.T) {
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
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "大和撫子"
		}),
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "天然"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("personality_type全件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/personality_type/read", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, 2, len(res.PersonalityTypes))
		assert.Equal(t, "大和撫子", res.PersonalityTypes[0].Type)
		assert.Equal(t, "天然", res.PersonalityTypes[1].Type)
	})

	t.Run("personality_type1件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/personality_type/read?id=%d", f.PersonalityTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, 1, len(res.PersonalityTypes))
		assert.Equal(t, "大和撫子", res.PersonalityTypes[0].Type)
	})

	t.Run("personality_type2件取得", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/personality_type/read?id=%d&id=%d", f.PersonalityTypes[0].ID, f.PersonalityTypes[1].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, 2, len(res.PersonalityTypes))
		assert.Equal(t, "大和撫子", res.PersonalityTypes[0].Type)
		assert.Equal(t, "天然", res.PersonalityTypes[1].Type)
	})

	t.Run("personality_type1件取得(存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/personality_type/read?id=0", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, 0, len(res.PersonalityTypes))
	})

	t.Run("personality_type2件取得(存在しない)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/personality_type/read?id=0&id=0", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, 0, len(res.PersonalityTypes))
	})

	t.Run("personality_type1件取得(不正なID)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, "/api/personality_type/read?id=a", nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("personality_type2件取得(内1件不正なID)", func(t *testing.T) {
		// リクエストの作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/personality_type/read?id=a&id=%d", f.PersonalityTypes[0].ID), nil)
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewReadHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdatePersonalityTypeHandler(t *testing.T) {
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
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "大和撫子"
		}),
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "天然"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("personality_type更新失敗", func(t *testing.T) {
		// リクエストの作成
		updatePersonalityType := IDs{IDs: []int64{}}
		b, err := json.Marshal(updatePersonalityType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut,"/api/personality_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewUpdateHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
		assert.Equal(t, "大和撫子", actuals[0].Type)
		assert.Equal(t, "天然", actuals[1].Type)
	})

	t.Run("personality_type更新", func(t *testing.T) {
		// リクエストの作成
		updatePersonalityType := PersonalityTypesJson{[]PersonalityType{
			{
				ID:   f.PersonalityTypes[0].ID,
				Type: "クール",
			},
			{
				ID:   f.PersonalityTypes[1].ID,
				Type: "ミステリアス",
			},
		}}
		b, err := json.Marshal(updatePersonalityType)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut,"/api/personality_type/update", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewUpdateHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスのデコード
		var res PersonalityTypesJson
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		// レスポンスの検証
		assert.Equal(t, updatePersonalityType, res)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
		assert.Equal(t, "クール", actuals[0].Type)
		assert.Equal(t, "ミステリアス", actuals[1].Type)
	})
}

func TestDeletePersonalityTypeHandler(t *testing.T) {
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
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "大和撫子"
		}),
		fixtures.NewPersonalityType(ctx, func(s *fixtures.PersonalityType) {
			s.Type = "天然"
		}),
	)

	// テスト対象のハンドラを作成
	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)

	t.Run("personality_type削除失敗", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{IDs: []int64{}}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/personality_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewDeleteHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(actuals))
		assert.Equal(t, "大和撫子", actuals[0].Type)
		assert.Equal(t, "天然", actuals[1].Type)
	})

	t.Run("personality_type削除", func(t *testing.T) {
		// リクエストの作成
		delIDs := IDs{
			IDs: []int64{f.PersonalityTypes[0].ID, f.PersonalityTypes[1].ID},
		}
		b, err := json.Marshal(delIDs)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, "/api/personality_type/delete", bytes.NewBuffer(b))
		// レスポンスの作成
		w := httptest.NewRecorder()
		// テスト対象のハンドラを実行
		handler := NewDeleteHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)

		var res IDs
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, delIDs, res)

		var actuals []PersonalityType
		err = tx.SelectContext(ctx, &actuals, "SELECT * FROM personality_type")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(actuals))
	})
}
