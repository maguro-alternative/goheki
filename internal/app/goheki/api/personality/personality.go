package personality

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
	var personalitiesJson PersonalitiesJson
	query := `
		INSERT INTO personality (
			entry_id,
			type_id
		) VALUES (
			:entry_id,
			:type_id
		)
	`
	// json decode
	err := json.NewDecoder(r.Body).Decode(&personalitiesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// validation
	err = personalitiesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, personality := range personalitiesJson.Personalities {
		// validation
		err = personality.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, personality)
		if err != nil {
			log.Printf(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	// json encode
	err = json.NewEncoder(w).Encode(&personalitiesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	var personalitiesJson PersonalitiesJson
	query := `
		SELECT
			entry_id,
			type_id
		FROM
			personality
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				type_id
			FROM
				personality
		`
		err := h.svc.DB.SelectContext(r.Context(), &personalitiesJson.Personalities, query)
		if err != nil {
			log.Printf(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// json encode
		err = json.NewEncoder(w).Encode(&personalitiesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				type_id
			FROM
				personality
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &personalitiesJson.Personalities, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// json encode
		err = json.NewEncoder(w).Encode(&personalitiesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("in query error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &personalitiesJson.Personalities, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// json encode
	err = json.NewEncoder(w).Encode(&personalitiesJson)
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
	var personalitiesJson PersonalitiesJson
	query := `
		UPDATE
			personality
		SET
			type_id = :type_id
		WHERE
			entry_id = :entry_id
	`
	err := json.NewDecoder(r.Body).Decode(&personalitiesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// validation
	err = personalitiesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	for _, personality := range personalitiesJson.Personalities {
		// validation
		err = personality.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, personality)
		if err != nil {
			log.Printf(fmt.Sprintf("update error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	// json encode
	err = json.NewEncoder(w).Encode(&personalitiesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			personality
		WHERE
			entry_id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// validation
	err = delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				personality
			WHERE
				entry_id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("delete error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// json encode
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
		log.Printf(fmt.Sprintf("in query error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("delete error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// json encode
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
