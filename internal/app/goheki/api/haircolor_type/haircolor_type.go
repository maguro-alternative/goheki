package haircolortype

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
	var hairColorTypesJson HairColorTypesJson
	query := `
		INSERT INTO haircolor_type (
			color
		) VALUES (
			:color
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = hairColorTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hairColorType := range hairColorTypesJson.HairColorTypes {
		// jsonのバリデーションを通過したデータをDBに登録
		err = hairColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairColorType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&hairColorTypesJson)
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
	// GET以外は受け付けない
	if r.Method != http.MethodGet {
		return
	}
	var hairColorTypesJson HairColorTypesJson
	query := `
		SELECT
			id,
			color
		FROM
			haircolor_type
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	// idが指定されていない場合は全件取得
	if !ok {
		query = `
			SELECT
				id,
				color
			FROM
				haircolor_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorTypesJson.HairColorTypes, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// レスポンスの作成
		err = json.NewEncoder(w).Encode(&hairColorTypesJson)
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
				color
			FROM
				haircolor_type
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorTypesJson.HairColorTypes, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// レスポンスの作成
		err = json.NewEncoder(w).Encode(&hairColorTypesJson)
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
	err = h.svc.DB.SelectContext(r.Context(), &hairColorTypesJson.HairColorTypes, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// レスポンスの作成
	err = json.NewEncoder(w).Encode(&hairColorTypesJson)
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
	// PUT以外は受け付けない
	if r.Method != http.MethodPut {
		return
	}
	var hairColorTypesJson HairColorTypesJson
	query := `
		UPDATE
			haircolor_type
		SET
			color = :color
		WHERE
			id = :id
	`
	// jsonのデコード
	err := json.NewDecoder(r.Body).Decode(&hairColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = hairColorTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hairColorType := range hairColorTypesJson.HairColorTypes {
		// jsonのバリデーションを通過したデータをDBに登録
		err = hairColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairColorType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&hairColorTypesJson)
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
			haircolor_type
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				haircolor_type
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
