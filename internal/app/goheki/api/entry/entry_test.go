package entry

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

func TestCreateEntryHandler(t *testing.T) {
	t.Run("entry登録", func(t *testing.T) {
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
		var ids []int64
		// テストデータの準備
		query := `
			INSERT INTO source (
				name,
				url,
				type
			) VALUES (
				:name,
				:url,
				:type
			)
		`
		sources := []Source{
			{
				Name: "テストソース1",
				Url:  "https://example.com/image1.png",
				Type: "anime",
			},
			{
				Name: "テストソース2",
				Url:  "https://example.com/image2.png",
				Type: "game",
			},
		}
		for _, source := range sources {
			_, err = tx.NamedExecContext(ctx, query, source)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				source
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		entrys := []Entry{
			{
				SourceID:  ids[0],
				Name:      "テストエントリ1",
				Image:     "https://example.com/image1.png",
				Content:   "テスト内容1",
				CreatedAt: fixedTime,
			},
			{
				SourceID:  ids[1],
				Name:      "テストエントリ2",
				Image:     "https://example.com/image2.png",
				Content:   "テスト内容2",
				CreatedAt: fixedTime,
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entrys)
		req, err := http.NewRequest(http.MethodPost, "/api/entry/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys, actual)
	})
}

func TestReadEntryHandler(t *testing.T) {
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
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ1"
			e.Image = "https://example.com/image1.png"
			e.Content = "テスト内容1"
			e.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ2"
			e.Image = "https://example.com/image2.png"
			e.Content = "テスト内容2"
			e.CreatedAt = fixedTime
		})),
	)
	entrys := []Entry{
		{
			SourceID:  *f.Sources[0].ID,
			Name:      f.Entrys[0].Name,
			Image:     f.Entrys[0].Image,
			Content:   f.Entrys[0].Content,
			CreatedAt: f.Entrys[0].CreatedAt,
		},
		{
			SourceID:  *f.Sources[1].ID,
			Name:      f.Entrys[1].Name,
			Image:     f.Entrys[1].Image,
			Content:   f.Entrys[1].Content,
			CreatedAt: f.Entrys[1].CreatedAt,
		},
	}
	t.Run("entry全件取得", func(t *testing.T) {

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		eJson, err := json.Marshal(&entrys)
		req, err := http.NewRequest(http.MethodGet, "/api/entry/read", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys, actual)
	})

	t.Run("entry1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d", *f.Entrys[0].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys[0], actual[0])
	})

	t.Run("entry2件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=%d", *f.Entrys[0].ID, *f.Entrys[1].ID), nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entrys, actual)
	})
}

func TestUpdateEntryHandler(t *testing.T) {
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
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ1"
			e.Image = "https://example.com/image1.png"
			e.Content = "テスト内容1"
			e.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ2"
			e.Image = "https://example.com/image2.png"
			e.Content = "テスト内容2"
			e.CreatedAt = fixedTime
		})),
	)

	updateEntrys := []Entry{
		{
			SourceID:  *f.Sources[0].ID,
			Name:      "テストエントリ3",
			Image:     "https://example.com/image3.png",
			Content:   "テスト内容3",
			CreatedAt: fixedTime,
		},
		{
			SourceID:  *f.Sources[1].ID,
			Name:      "テストエントリ4",
			Image:     "https://example.com/image4.png",
			Content:   "テスト内容4",
			CreatedAt: fixedTime,
		},
	}
	t.Run("entry更新", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewUpdateHandler(indexService)
		eJson, err := json.Marshal(&updateEntrys)
		req, err := http.NewRequest(http.MethodPut, "/api/entry/update", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []Entry
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, updateEntrys, actual)
	})
}

func TestDeleteEntryHandler(t *testing.T) {
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
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ1"
			e.Image = "https://example.com/image1.png"
			e.Content = "テスト内容1"
			e.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func (s *fixtures.Source)  {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func (e *fixtures.Entry)  {
			e.Name = "テストエントリ2"
			e.Image = "https://example.com/image2.png"
			e.Content = "テスト内容2"
			e.CreatedAt = fixedTime
		})),
	)
	t.Run("entry削除", func(t *testing.T) {
		ids := []int64{*f.Entrys[0].ID, *f.Entrys[1].ID}
		delIDs := IDs{IDs: ids}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewDeleteHandler(indexService)
		eJson, err := json.Marshal(&delIDs)
		req, err := http.NewRequest(http.MethodDelete, "/api/entry/delete", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		assert.NoError(t, err)

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual IDs
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, delIDs, actual)
	})
}
