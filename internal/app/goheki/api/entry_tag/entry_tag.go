package entry_tag

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
	if r.Method != http.MethodPost {
		return
	}

	var entryTagsJson EntryTagsJson
	query := `
		INSERT INTO entry_tag (
			entry_id,
			tag_id
		) VALUES (
			:entry_id,
			:tag_id
		)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &entryTagsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, entryTag := range entryTagsJson.EntryTags {
		_, err := h.svc.DB.NamedExecContext(r.Context(), query, entryTag)
		if err != nil {
			log.Println(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&entryTagsJson)
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
	var entryTagsJson EntryTagsJson
	query := `
		SELECT
			id,
			entry_id,
			tag_id
		FROM
			entry_tag
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				entry_id,
				tag_id
			FROM
				entry_tag
		`
		err := h.svc.DB.SelectContext(r.Context(), &entryTagsJson, query)
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&entryTagsJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				entry_id,
				tag_id
			FROM
				entry_tag
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &entryTagsJson, query, queryIDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&entryTagsJson)
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
	err = h.svc.DB.SelectContext(r.Context(), &entryTagsJson, query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&entryTagsJson)
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

	var entryTagsJson EntryTagsJson
	query := `
		UPDATE
			entry_tag
		SET
			entry_id = :entry_id,
			tag_id = :tag_id
		WHERE
			id = :id
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &entryTagsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, entryTag := range entryTagsJson.EntryTags {
		_, err := h.svc.DB.NamedExecContext(r.Context(), query, entryTag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&entryTagsJson)
	if err != nil {
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
			entry_tag
		WHERE
			id IN (?)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				entry_tag
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs),query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
