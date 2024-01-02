package tag

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type Tag struct {
	ID        *int64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
}

type DeleteIDs struct {
	IDs []int64 `json:"ids"`
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
	var tags []Tag
	query := `
		INSERT INTO tag (
			name
		) VALUES (
			:name
		)
	`
	err := json.NewDecoder(r.Body).Decode(&tags)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, tag := range tags {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, tag)
		if err != nil {
			log.Fatal(fmt.Sprintf("insert error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(tags)
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
	var tags []Tag
	query := `
		SELECT
			id,
			name
		FROM
			tag
	`
	err := h.svc.DB.SelectContext(r.Context(), &tags, query)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	json.NewEncoder(w).Encode(tags)
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
	var tags []Tag
	query := `
		UPDATE
			tag
		SET
			name = :name
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&tags)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, tag := range tags {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, tag)
		if err != nil {
			log.Fatal(fmt.Sprintf("update error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(tags)
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
	var delIDs DeleteIDs
	query := `DELETE FROM tag WHERE id IN (?)`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in error: %v", err), delIDs.IDs)
	}
	query = db.Rebind(len(delIDs.IDs),query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("delete error: %v", err), query, args)
	}
	json.NewEncoder(w).Encode(&delIDs)
}