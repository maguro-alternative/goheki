package eyecolor

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
	var eyeColors []EyeColor
	query := `
		INSERT INTO eyecolor (
			entry_id,
			color_id
		) VALUES (
			:entry_id,
			:color_id
		)
	`
	err := json.NewDecoder(r.Body).Decode(&eyeColors)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, eyeColor := range eyeColors {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColor)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColors)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
	if r.Method != http.MethodGet {
		return
	}
	var eyeColors []EyeColor
	query := `
		SELECT
			entry_id,
			color_id
		FROM
			eyecolor
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				color_id
			FROM
				eyecolor
		`
		err := h.svc.DB.SelectContext(r.Context(), &eyeColors, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&eyeColors)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				color_id
			FROM
				eyecolor
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &eyeColors, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&eyeColors)
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
	err = h.svc.DB.SelectContext(r.Context(), &eyeColors, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&eyeColors)
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
	var eyeColors []EyeColor
	query := `
		UPDATE
			eyecolor
		SET
			color_id = :color_id
		WHERE
			entry_id = :entry_id
	`
	err := json.NewDecoder(r.Body).Decode(&eyeColors)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v", err))
		return
	}
	for _, eyeColor := range eyeColors {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColor)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColors)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
			eyecolor
		WHERE
			entry_id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v", err))
		return
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				eyecolor
			WHERE
				entry_id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
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
}
