package eyecolor

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
	var eyeColorsJson EyeColorsJson
	query := `
		INSERT INTO eyecolor (
			entry_id,
			color_id
		) VALUES (
			:entry_id,
			:color_id
		)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &eyeColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, eyeColor := range eyeColorsJson.EyeColors {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColor)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var eyeColorsJson EyeColorsJson
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
		err := h.svc.DB.SelectContext(r.Context(), &eyeColorsJson.EyeColors, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&eyeColorsJson)
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
				eyecolor
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &eyeColorsJson.EyeColors, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&eyeColorsJson)
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
	err = h.svc.DB.SelectContext(r.Context(), &eyeColorsJson.EyeColors, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&eyeColorsJson)
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
	var eyeColorsJson EyeColorsJson
	query := `
		UPDATE
			eyecolor
		SET
			color_id = :color_id
		WHERE
			entry_id = :entry_id
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &eyeColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, eyeColor := range eyeColorsJson.EyeColors {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColor)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColorsJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
