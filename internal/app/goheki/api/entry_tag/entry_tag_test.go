package entry_tag

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

func TestCreateEntryTagHandler(t *testing.T) {
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
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ1"
			s.Image = "https://example.com/image1.png"
			s.Content = "テスト内容1"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	t.Run("entry_tag登録", func(t *testing.T) {
		entryTags := []EntryTag{
			{
				EntryID: f.Entrys[0].ID,
				TagID:   f.Tags[0].ID,
			},
			{
				EntryID: f.Entrys[1].ID,
				TagID:   f.Tags[1].ID,
			},
		}

		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewCreateHandler(indexService)
		eJson, err := json.Marshal(&entryTags)
		req, err := http.NewRequest(http.MethodPost, "/api/entry_tag/create", bytes.NewBuffer(eJson))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		// テストの実行
		h.ServeHTTP(w, req)

		// ロールバック
		tx.RollbackCtx(ctx)

		// 応答の検証
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var actual []EntryTag
		err = json.NewDecoder(res.Body).Decode(&actual)
		assert.NoError(t, err)

		assert.Equal(t, entryTags, actual)
	})
}

func TestReadEntryTagHandler(t *testing.T) {
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
			s.Name = "テストソース1"
			s.Url = "https://example.com/image1.png"
			s.Type = "anime"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ1"
			s.Image = "https://example.com/image1.png"
			s.Content = "テスト内容1"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewSource(ctx, func(s *fixtures.Source) {
			s.Name = "テストソース2"
			s.Url = "https://example.com/image2.png"
			s.Type = "game"
		}).Connect(fixtures.NewEntry(ctx, func(s *fixtures.Entry) {
			s.Name = "テストエントリ2"
			s.Image = "https://example.com/image2.png"
			s.Content = "テスト内容2"
			s.CreatedAt = fixedTime
		})),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ1"
		}),
		fixtures.NewTag(ctx, func(s *fixtures.Tag) {
			s.Name = "テストタグ2"
		}),
	)
	f.Build(t,
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[0].ID
			s.TagID = f.Tags[0].ID
		}),
		fixtures.NewEntryTag(ctx, func(s *fixtures.EntryTag) {
			s.EntryID = f.Entrys[1].ID
			s.TagID = f.Tags[1].ID
		}),
	)
	t.Run("entry全件取得", func(t *testing.T) {
		entryTags := []EntryTag{
			{
				ID:      f.EntryTags[0].ID,
				EntryID: f.Entrys[0].ID,
				TagID:   f.Tags[0].ID,
			},
			{
				ID:      f.EntryTags[1].ID,
				EntryID: f.Entrys[1].ID,
				TagID:   f.Tags[1].ID,
			},
		}
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		eJson, err := json.Marshal(&entryTags)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, "/api/entry_tag/read", bytes.NewBuffer(eJson))
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

		var actuals []EntryTag
		err = json.NewDecoder(res.Body).Decode(&actuals)
		assert.NoError(t, err)

		assert.Equal(t, entryTags, actuals)
	})

	t.Run("entry1件取得", func(t *testing.T) {
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d", *f.EntryTags[0].ID), nil)
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

		var actuals []EntryTag
		err = json.NewDecoder(res.Body).Decode(&actuals)
		assert.NoError(t, err)

		assert.Equal(t, f.EntryTags[0].ID, actuals[0].ID)
	})

	t.Run("entry2件取得", func(t *testing.T) {
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
		var ids []int64
		var idsJson IDs
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
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
		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entrys {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				entry
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		idsJson.IDs = ids
		var indexService = service.NewIndexService(
			tx,
			cookie.Store,
			env,
		)
		// テストの実行
		h := NewReadHandler(indexService)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/entry/read?id=%d&id=%d", ids[0], ids[1]), nil)
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
	t.Run("entry更新", func(t *testing.T) {
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
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

		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
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
		updateEntrys := []Entry{
			{
				SourceID:  ids[0],
				Name:      "テストエントリ3",
				Image:     "https://example.com/image3.png",
				Content:   "テスト内容3",
				CreatedAt: fixedTime,
			},
			{
				SourceID:  ids[1],
				Name:      "テストエントリ4",
				Image:     "https://example.com/image4.png",
				Content:   "テスト内容4",
				CreatedAt: fixedTime,
			},
		}
		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entrys {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				entry
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
		for i, id := range ids {
			updateEntrys[i].ID = &id
		}
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

func TestDeleteEntryTagHandler(t *testing.T) {
	t.Run("entry削除", func(t *testing.T) {
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
		var ids []int64
		fixedTime := time.Date(2023, time.December, 27, 10, 55, 22, 0, time.UTC)
		// テストデータの準備
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
		entry := []Entry{
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
		query = `
			INSERT INTO entry (
				source_id,
				name,
				image,
				content,
				created_at
			) VALUES (
				:source_id,
				:name,
				:image,
				:content,
				:created_at
			)
		`
		for _, entry := range entry {
			_, err = tx.NamedExecContext(ctx, query, entry)
			assert.NoError(t, err)
		}
		query = `
			SELECT
				id
			FROM
				entry
		`
		err = tx.SelectContext(ctx, &ids, query)
		assert.NoError(t, err)
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
