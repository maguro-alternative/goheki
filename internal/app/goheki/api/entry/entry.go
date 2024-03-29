package entry

import (
	"fmt"
	"log"
	"io"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type CreateHandler struct {
	svc *service.IndexService
}

func NewCreateHandler(svc *service.IndexService) *CreateHandler {
	return &CreateHandler{
		svc: svc,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は受け付けない
	if r.Method != http.MethodPost {
		return
	}
	var entriesJson EntriesJson
	query := `
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
	// リクエストボディを読み込む
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(entriesJson.Entries) == 0 {
		log.Println("json unexpected error: empty body")
		http.Error(w, "json unexpected error: empty body", http.StatusBadRequest)
	}
	// jsonバリデーション
	err = entriesJson.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, entry := range entriesJson.Entries {
		// jsonバリデーション
		err = entry.Validate()
		if err != nil {
			log.Println(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Println(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// レスポンスボディに書き込む
	err = json.NewEncoder(w).Encode(&entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ReadHandler struct {
	svc *service.IndexService
}

func NewReadHandler(svc *service.IndexService) *ReadHandler {
	return &ReadHandler{
		svc: svc,
	}
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	var entriesJson EntriesJson
	query := `
		SELECT
			source_id,
			name,
			image,
			content,
			created_at
		FROM
			entry
		WHERE
			id IN (?)
	`
	// クエリパラメータからidを取得
	queryIDs, ok := r.URL.Query()["id"]
	// idが指定されていない場合は全件取得
	if !ok {
		query = `
			SELECT
				source_id,
				name,
				image,
				content,
				created_at
			FROM
				entry
		`
		// 全件取得
		err := h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query)
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// レスポンスボディに書き込む
		err = json.NewEncoder(w).Encode(&entriesJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				source_id,
				name,
				image,
				content,
				created_at
			FROM
				entry
			WHERE
				id = $1
		`
		// 1件取得
		err := h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query, queryIDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// レスポンスボディに書き込む
		err = json.NewEncoder(w).Encode(&entriesJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Println(fmt.Sprintf("db.In error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	// 複数件取得
	err = h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// レスポンスボディに書き込む
	err = json.NewEncoder(w).Encode(&entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type UpdateHandler struct {
	svc *service.IndexService
}

func NewUpdateHandler(svc *service.IndexService) *UpdateHandler {
	return &UpdateHandler{
		svc: svc,
	}
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// PUTメソッド以外は受け付けない
	if r.Method != http.MethodPut {
		return
	}
	var entriesJson EntriesJson
	query := `
		UPDATE
			entry
		SET
			source_id = :source_id,
			name = :name,
			image = :image,
			content = :content,
			created_at = :created_at
		WHERE
			id = :id
	`
	// リクエストボディを読み込む
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	if len(entriesJson.Entries) == 0 {
		log.Println("json unexpected error: empty body")
		http.Error(w, "json unexpected error: empty body", http.StatusUnprocessableEntity)
	}
	// jsonバリデーション
	err = entriesJson.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, entry := range entriesJson.Entries {
		// jsonバリデーション
		err = entry.Validate()
		if err != nil {
			log.Println(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		// 更新
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Println(fmt.Sprintf("update error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// レスポンスボディに書き込む
	err = json.NewEncoder(w).Encode(&entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type DeleteHandler struct {
	svc *service.IndexService
}

func NewDeleteHandler(svc *service.IndexService) *DeleteHandler {
	return &DeleteHandler{
		svc: svc,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// DELETEメソッド以外は受け付けない
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			entry
		WHERE
			id IN (?)
	`
	// リクエストボディを読み込む
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// jsonバリデーション
	err = delIDs.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	if len(delIDs.IDs) == 0 {
		return
	// 1件の場合はIN句を使わない
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				entry
			WHERE
				id = $1
		`
		// 1件削除
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("delete error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// レスポンスボディに書き込む
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Println(fmt.Sprintf("db.In error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("delete error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// レスポンスボディに書き込む
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
