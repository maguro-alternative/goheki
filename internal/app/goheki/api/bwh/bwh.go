package bwh

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
	var bwhsJson BWHsJson
	query := `
		INSERT INTO bwh (
			entry_id,
			bust,
			waist,
			hip,
			height,
			weight
		) VALUES (
			:entry_id,
			:bust,
			:waist,
			:hip,
			:height,
			:weight
		)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &bwhsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, bwh := range bwhsJson.BWHs {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, bwh)
		if err != nil {
			log.Println(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&bwhsJson)
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
	var bwhsJson BWHsJson
	query := `
		SELECT
			entry_id,
			bust,
			waist,
			hip,
			height,
			weight
		FROM
			bwh
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query := `
			SELECT
				entry_id,
				bust,
				waist,
				hip,
				height,
				weight
			FROM
				bwh
		`
		err := h.svc.DB.SelectContext(r.Context(), &bwhsJson.BWHs, query)
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&bwhsJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(queryIDs) == 1 {
		query := `
			SELECT
				entry_id,
				bust,
				waist,
				hip,
				height,
				weight
			FROM
				bwh
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &bwhsJson.BWHs, query, queryIDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&bwhsJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Println(fmt.Sprintf("in error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &bwhsJson.BWHs, query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&bwhsJson)
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
	var bwhsJson BWHsJson
	query := `
		UPDATE
			bwh
		SET
			bust = :bust,
			waist = :waist,
			hip = :hip,
			height = :height,
			weight = :weight
		WHERE
			entry_id = :entry_id
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &bwhsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, bwh := range bwhsJson.BWHs {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, bwh)
		if err != nil {
			log.Println(fmt.Sprintf("update error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&bwhsJson)
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
			bwh
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
		log.Println(fmt.Sprintf("json unmarshal error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query := `
			DELETE FROM
				bwh
			WHERE
				entry_id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			http.Error(w, fmt.Sprintf("delete error: %v", err), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			http.Error(w, fmt.Sprintf("json encode error: %v", err), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("in error: %v", err), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs),query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("delete error: %v", err), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("json encode error: %v", err), http.StatusInternalServerError)
	}
}
