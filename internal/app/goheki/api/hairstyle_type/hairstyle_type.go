package hairstyletype

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
	// リクエストのメソッドがPOSTでなければ終了
	if r.Method != http.MethodPost {
		return
	}
	var hairStyleTypesJson HairStyleTypesJson
	query := `
		INSERT INTO hairstyle_type (
			style
		) VALUES (
			:style
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairStyleTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// リクエストボディのバリデーション
	err = hairStyleTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hairStyleType := range hairStyleTypesJson.HairStyleTypes {
		// リクエストボディのバリデーション
		err = hairStyleType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairStyleType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&hairStyleTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	// リクエストのメソッドがGETでなければ終了
	if r.Method != http.MethodGet {
		return
	}
	var hairStyleTypesJson HairStyleTypesJson
	query := `
		SELECT
			id,
			style
		FROM
			hairstyle_type
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				style
			FROM
				hairstyle_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairStyleTypesJson.HairStyleTypes, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// レスポンスの作成
		err = json.NewEncoder(w).Encode(&hairStyleTypesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				style
			FROM
				hairstyle_type
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairStyleTypesJson.HairStyleTypes, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// レスポンスの作成
		err = json.NewEncoder(w).Encode(&hairStyleTypesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairStyleTypesJson.HairStyleTypes, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&hairStyleTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	// リクエストのメソッドがPUTでなければ終了
	if r.Method != http.MethodPut {
		return
	}
	var hairStyleTypesJson HairStyleTypesJson
	query := `
		UPDATE
			hairstyle_type
		SET
			style = :style
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&hairStyleTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// リクエストボディのバリデーション
	err = hairStyleTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hairStyleType := range hairStyleTypesJson.HairStyleTypes {
		// リクエストボディのバリデーション
		err = hairStyleType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairStyleType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&hairStyleTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			hairstyle_type
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// リクエストボディのバリデーション
	err = delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				hairstyle_type
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
