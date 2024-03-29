package hairlength

import (
	"fmt"
	"log"

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
	// POST以外は受け付けない
	if r.Method != http.MethodPost {
		return
	}
	var hairLengthsJson HairLengthsJson
	query := `
		INSERT INTO hairlength (
			entry_id,
			hairlength_type_id
		) VALUES (
			:entry_id,
			:hairlength_type_id
		)
	`
	if err := json.NewDecoder(r.Body).Decode(&hairLengthsJson); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, "json decode error", http.StatusBadRequest)
	}
	// jsonバリデーション
	err := hairLengthsJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, "validation error", http.StatusUnprocessableEntity)
	}
	for _, hl := range hairLengthsJson.HairLengths {
		// jsonのバリデーションを通過したデータをDBに登録
		err = hl.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, "validation error", http.StatusUnprocessableEntity)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hl); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, "db error", http.StatusInternalServerError)
		}
	}
	// jsonを返す
	err = json.NewEncoder(w).Encode(&hairLengthsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, "json encode error", http.StatusInternalServerError)
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
	var hairLengthsJson HairLengthsJson
	query := `
		SELECT
			entry_id,
			hairlength_type_id
		FROM
			hairlength
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				hairlength_type_id
			FROM
				hairlength
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairLengthsJson.HairLengths, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, "db error", http.StatusInternalServerError)
		}
		// jsonを返す
		err = json.NewEncoder(w).Encode(&hairLengthsJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, "json encode error", http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				hairlength_type_id
			FROM
				hairlength
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairLengthsJson.HairLengths, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, "db error", http.StatusInternalServerError)
		}
		// jsonを返す
		err = json.NewEncoder(w).Encode(&hairLengthsJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, "json encode error", http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, "db error", http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairLengthsJson.HairLengths, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, "db error", http.StatusInternalServerError)
	}
	// jsonを返す
	err = json.NewEncoder(w).Encode(&hairLengthsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, "json encode error", http.StatusInternalServerError)
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
	if r.Method != http.MethodPut {
		return
	}
	var hairLengthsJson HairLengthsJson
	query := `
		UPDATE
			hairlength
		SET
			hairlength_type_id = :hairlength_type_id
		WHERE
			entry_id = :entry_id
	`
	// jsonをデコード
	if err := json.NewDecoder(r.Body).Decode(&hairLengthsJson); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, "json decode error", http.StatusBadRequest)
	}
	// jsonバリデーション
	err := hairLengthsJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, "validation error", http.StatusUnprocessableEntity)
	}
	for _, hl := range hairLengthsJson.HairLengths {
		// jsonのバリデーションを通過したデータをDBに登録
		err = hl.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, "validation error", http.StatusUnprocessableEntity)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hl); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, "db error", http.StatusInternalServerError)
		}
	}
	// jsonを返す
	err = json.NewEncoder(w).Encode(&hairLengthsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, "json encode error", http.StatusInternalServerError)
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
	// DELETE以外は受け付けない
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			hairlength
		WHERE
			entry_id IN (?)
	`
	// jsonをデコード
	if err := json.NewDecoder(r.Body).Decode(&delIDs); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, "json decode error", http.StatusBadRequest)
	}
	// jsonバリデーション
	err := delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, "validation error", http.StatusUnprocessableEntity)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				hairlength
			WHERE
				entry_id = $1
		`
		_, err := h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
