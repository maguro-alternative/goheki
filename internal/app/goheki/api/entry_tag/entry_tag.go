package entry_tag

import (
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

	var entryTags []EntryTag
	query := `
		INSERT INTO entry_tag (
			entry_id,
			tag_id
		) VALUES (
			:entry_id,
			:tag_id
		)
	`
	err := json.NewDecoder(r.Body).Decode(&entryTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, entryTag := range entryTags {
		_, err := h.svc.DB.NamedExecContext(r.Context(), query, entryTag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	err = json.NewEncoder(w).Encode(&entryTags)
	if err != nil {
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
	if r.Method != http.MethodGet {
		return
	}
	var entryTags []EntryTag
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
		err := h.svc.DB.SelectContext(r.Context(), &entryTags, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = json.NewEncoder(w).Encode(&entryTags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		err := h.svc.DB.SelectContext(r.Context(), &entryTags, query, queryIDs[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = json.NewEncoder(w).Encode(&entryTags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &entryTags, query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(&entryTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	var entryTags []EntryTag
	query := `
		UPDATE
			entry_tag
		SET
			entry_id = :entry_id,
			tag_id = :tag_id
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&entryTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, entryTag := range entryTags {
		_, err := h.svc.DB.NamedExecContext(r.Context(), query, entryTag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	err = json.NewEncoder(w).Encode(&entryTags)
	if err != nil {
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
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
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
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query = db.Rebind(len(delIDs.IDs),query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
