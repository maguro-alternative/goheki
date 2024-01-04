package entry

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
	"time"
)

type Entry struct {
	ID        *int64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Image     string    `db:"image" json:"image"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ID struct {
	ID int64 `json:"id"`
}

type IDs struct {
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
	var entrys []Entry
	query := `
		INSERT INTO entry (
			name,
			image,
			content,
			created_at
		) VALUES (
			:name,
			:image,
			:content,
			:created_at
		)
	`
	err := json.NewDecoder(r.Body).Decode(&entrys)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, entry := range entrys {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
		}
	}
	json.NewEncoder(w).Encode(&entrys)
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
	var entrys []Entry
	query := `
		SELECT
			name,
			image,
			content,
			created_at
		FROM
			entry
	`

	err := h.svc.DB.SelectContext(r.Context(), &entrys, query)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&entrys)
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
	var entry Entry
	var id ID
	query := `
		SELECT
			name,
			image,
			content,
			created_at
		FROM
			entry
		WHERE
			id = ?
	`
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}

	err = h.svc.DB.GetContext(r.Context(), &entry, query, id.ID)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&entry)
}

type MultipleReadHandler struct {
	svc *service.IndexService
}

func NewMultipleReadHandler(svc *service.IndexService) *AllReadHandler {
	return &AllReadHandler{
		svc: svc,
	}
}

func (h *MultipleReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	var entrys []Entry
	var ids IDs
	query := `
		SELECT
			name,
			image,
			content,
			created_at
		FROM
			entry
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	query, args, err := db.In(query, ids.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in error: %v", err), ids.IDs)
	}
	err = h.svc.DB.SelectContext(r.Context(), &entrys, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&entrys)
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
	var entrys []Entry
	query := `
		UPDATE
			entry
		SET
			name = :name,
			image = :image,
			content = :content,
			created_at = :created_at
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&entrys)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	//query, args, err := db.In(query, entrys)
	//if err != nil {
	//return
	//}
	for _, entry := range entrys {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
		}
	}
	json.NewEncoder(w).Encode(&entrys)
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
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in error: %v", err), delIDs.IDs)
	}
	query = db.Rebind(len(delIDs.IDs),query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&delIDs)
}
