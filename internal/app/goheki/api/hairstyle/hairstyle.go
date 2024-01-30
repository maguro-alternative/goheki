package hairstyle

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
	var hairStylesJson HairStylesJson
	query := `
		INSERT INTO hairstyle (
			entry_id,
			style_id
		) VALUES (
			:entry_id,
			:style_id
		)
	`
	// json読み込み
	if err := json.NewDecoder(r.Body).Decode(&hairStylesJson); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// jsonバリデーション
	err := hairStylesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, hs := range hairStylesJson.HairStyles {
		// jsonバリデーション
		err = hs.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hs); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&hairStylesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	// GET以外は受け付けない
	if r.Method != http.MethodGet {
		return
	}
	var hairStylesJson HairStylesJson
	query := `
		SELECT
			entry_id,
			style_id
		FROM
			hairstyle
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				style_id
			FROM
				hairstyle
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairStylesJson.HairStyles, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// json書き込み
		err = json.NewEncoder(w).Encode(&hairStylesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				style_id
			FROM
				hairstyle
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairStylesJson.HairStyles, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// json書き込み
		err = json.NewEncoder(w).Encode(&hairStylesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairStylesJson.HairStyles, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&hairStylesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	// PUT以外は受け付けない
	if r.Method != http.MethodPut {
		return
	}
	var hairStylesJson HairStylesJson
	query := `
		UPDATE
			hairstyle
		SET
			style_id = :style_id
		WHERE
			entry_id = :entry_id
	`
	// json読み込み
	if err := json.NewDecoder(r.Body).Decode(&hairStylesJson); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// jsonバリデーション
	err := hairStylesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, hs := range hairStylesJson.HairStyles {
		// jsonバリデーション
		err = hs.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hs); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&hairStylesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	// DELETE以外は受け付けない
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			hairstyle
		WHERE
			entry_id IN (?)
	`
	// json読み込み
	if err := json.NewDecoder(r.Body).Decode(&delIDs); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// jsonバリデーション
	err := delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				hairstyle
			WHERE
				entry_id = $1
		`
		_, err := h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// json書き込み
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
	// json書き込み
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
