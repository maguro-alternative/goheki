package entry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/maguro-alternative/goheki/configs/envconfig"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/service/cookie"

	//"github.com/maguro-alternative/goheki/pkg/db"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	ctx := context.Background()
	env, err := envconfig.NewEnv()
	assert.NoError(t, err)
	// データベースに接続
	//assert.Equal(t, "postgres", env.DatabaseURL)
	indexDB, err := sqlx.Open("postgres", env.DatabaseURL)
	assert.NoError(t, err)
	var indexService = service.NewIndexService(
		indexDB,
		cookie.Store,
		env,
	)
	// テストデータの準備
	entry := []Entry{
		{
			Name:     "テストエントリ1",
			Image:    "https://example.com/image1.png",
			Content:  "テスト内容1",
			CreateAt: time.Now(),
		},
		{
			Name:     "テストエントリ2",
			Image:    "https://example.com/image2.png",
			Content:  "テスト内容2",
			CreateAt: time.Now(),
		},
	}
	var entry2 []Entry

	// テストの実行
	h := NewCreateHandler(indexService)
	eJson, err := json.Marshal(&entry)
	req, err := http.NewRequest(http.MethodPost, "/api/entry/create", bytes.NewBuffer(eJson))
	assert.NoError(t, err)

	err = json.NewDecoder(req.Body).Decode(&entry2)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	// トランザクションの開始
	tx, err := h.svc.DB.BeginTx(ctx, nil)
	assert.NoError(t, err)

	// テストの実行
	h.ServeHTTP(w, req)

	// ロールバック
	tx.Rollback()

	// 応答の検証
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var actual []Entry
	err = json.NewDecoder(res.Body).Decode(&actual)
	assert.NoError(t, err)

	assert.Equal(t, entry2, actual)
	assert.Equal(t, nil, actual)
}
