package source

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type Source struct {
	ID      *int64 `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Url     string `db:"url" json:"url"`
	Type    string `db:"type" json:"type"`
}

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
	var sources []Source
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
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, source := range sources {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Fatal(fmt.Sprintf("insert error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(sources)
}

type AllReadHandler struct {
	svc *service.IndexService
}

func NewAllReadHandler(svc *service.IndexService) *AllReadHandler {
	return &AllReadHandler{
		svc: svc,
	}
}

func (h *AllReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	var sources []Source
	query := `
		SELECT
			id,
			name,
			url,
			type
		FROM
			source
	`
	err := h.svc.DB.SelectContext(r.Context(), &sources, query)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	json.NewEncoder(w).Encode(sources)
}

type GetReadHandler struct {
	svc *service.IndexService
}

func NewGetReadHandler(svc *service.IndexService) *GetReadHandler {
	return &GetReadHandler{
		svc: svc,
	}
}

func (h *GetReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	var sources Source
	query := `
		SELECT
			id,
			name,
			url,
			type
		FROM
			source
		WHERE
			id = $1
	`
	id := r.URL.Query().Get("id")
	err := h.svc.DB.GetContext(r.Context(), &sources, query, id)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	json.NewEncoder(w).Encode(sources)
}

type MultipleReadHandler struct {
	svc *service.IndexService
}

func NewMultipleReadHandler(svc *service.IndexService) *MultipleReadHandler {
	return &MultipleReadHandler{
		svc: svc,
	}
}

func (h *MultipleReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	var sources []Source
	query := `
		SELECT
			id,
			name,
			url,
			type
		FROM
			source
		WHERE
			id IN (?)
	`
	entryIDs, ok := r.URL.Query()["id"]
	if !ok {
		log.Fatal(fmt.Sprintf("id not found: %v", r.URL.Query()))
	}
	query, args, err := db.In(query, entryIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in error: %v", err), entryIDs)
	}
	query = db.Rebind(len(entryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &sources, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	json.NewEncoder(w).Encode(sources)
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
	var sources []Source
	query := `
		UPDATE
			source
		SET
			name = :name,
			url = :url,
			type = :type
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, source := range sources {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Fatal(fmt.Sprintf("update error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(sources)
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
	var sources []Source
	query := `
		DELETE FROM
			source
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, source := range sources {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Fatal(fmt.Sprintf("delete error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(sources)
}
