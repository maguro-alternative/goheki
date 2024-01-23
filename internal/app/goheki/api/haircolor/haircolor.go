package haircolor

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
	var hairColorsJson HairColorsJson
	query := `
		INSERT INTO haircolor (
			entry_id,
			color_id
		) VALUES (
			:entry_id,
			:color_id
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = hairColorsJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, hc := range hairColorsJson.HairColors {
		err = hc.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hc); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&hairColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	var hairColorsJson HairColorsJson
	query := `
		SELECT
			entry_id,
			color_id
		FROM
			haircolor
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
				haircolor
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorsJson.HairColors, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&hairColorsJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				color_id
			FROM
				haircolor
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorsJson.HairColors, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&hairColorsJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairColorsJson.HairColors, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&hairColorsJson)
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
	if r.Method != http.MethodPut {
		return
	}
	var hairColorsJson HairColorsJson
	query := `
		UPDATE
			haircolor
		SET
			color_id = :color_id
		WHERE
			entry_id = :entry_id
	`
	if err := json.NewDecoder(r.Body).Decode(&hairColorsJson); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := hairColorsJson.Validate(); err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, hc := range hairColorsJson.HairColors {
		if err := hc.Validate(); err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if _, err := h.svc.DB.NamedExecContext(r.Context(), query, hc); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err := json.NewEncoder(w).Encode(&hairColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	if err := json.NewDecoder(r.Body).Decode(&delIDs); err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query := `
		DELETE FROM
			haircolor
		WHERE
			entry_id IN (?)
	`
	err := delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				haircolor
			WHERE
				entry_id = $1
		`
		if _, err := h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0]); err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs), query)
	if _, err := h.svc.DB.ExecContext(r.Context(), query, args...); err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
