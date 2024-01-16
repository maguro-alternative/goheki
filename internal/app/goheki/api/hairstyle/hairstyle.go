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
	if r.Method != http.MethodPost {
		return
	}
	var hairStyles []HairStyle
	query := `
		INSERT INTO hairstyle (
			entry_id,
			style_id
		) VALUES (
			:entry_id,
			:style_id
		)
	`
	if err := json.NewDecoder(r.Body).Decode(&hairStyles); err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, hs := range hairStyles {
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hs); err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
	}
	err := json.NewEncoder(w).Encode(&hairStyles)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
	var hairStyles []HairStyle
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
		err := h.svc.DB.SelectContext(r.Context(), &hairStyles, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&hairStyles)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
		err := h.svc.DB.SelectContext(r.Context(), &hairStyles, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&hairStyles)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairStyles, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&hairStyles)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
	var hairStyles []HairStyle
	query := `
		UPDATE
			hairstyle
		SET
			style_id = :style_id
		WHERE
			entry_id = :entry_id
	`
	if err := json.NewDecoder(r.Body).Decode(&hairStyles); err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, hs := range hairStyles {
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hs); err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
	}
	err := json.NewEncoder(w).Encode(&hairStyles)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
			hairstyle
		WHERE
			entry_id IN (?)
	`
	if err := json.NewDecoder(r.Body).Decode(&delIDs); err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
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
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
	}
}
