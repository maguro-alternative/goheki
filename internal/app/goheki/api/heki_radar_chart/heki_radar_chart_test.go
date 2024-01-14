package hekiradarchart

import (
	"bytes"
	"context"
	"encoding/json"

	//"fmt"
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

func TestCreateHekiRadarChartHandler(t *testing.T) {
	// setup
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
	)

	// テストデータの作成
	charts := []HekiRadarChart{
		{
			EntryID: f.Entrys[0].ID,
			AI:      1,
			NU:      2,
		},
		{
			EntryID: f.Entrys[1].ID,
			AI:      3,
			NU:      4,
		},
	}

	var indexService = service.NewIndexService(
		tx,
		cookie.Store,
		env,
	)
	t.Run("heki_rader_chart登録", func(t *testing.T) {
		// リクエストの作成
		b, err := json.Marshal(charts)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/api/heki_radar_chart/create", bytes.NewBuffer(b))
		assert.NoError(t, err)
		// レスポンスの作成
		w := httptest.NewRecorder()
		handler := NewCreateHandler(indexService)
		handler.ServeHTTP(w, req)
		// レスポンスの検証
		assert.Equal(t, http.StatusOK, w.Code)
		// レスポンスの検証
		var res []HekiRadarChart
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, charts, res)
	})
}
