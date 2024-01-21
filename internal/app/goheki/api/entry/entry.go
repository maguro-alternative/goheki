package entry

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
	err := json.NewDecoder(r.Body).Decode(&entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, entry := range entriesJson.Entries {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Println(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
	queryIDs, ok := r.URL.Query()["id"]
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
		err := h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query)
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
		err := h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query, queryIDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&entriesJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Println(fmt.Sprintf("db.In error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &entriesJson.Entries, query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	err := json.NewDecoder(r.Body).Decode(&entriesJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, entry := range entriesJson.Entries {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Println(fmt.Sprintf("update error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				entry
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("delete error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Println(fmt.Sprintf("db.In error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("delete error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
