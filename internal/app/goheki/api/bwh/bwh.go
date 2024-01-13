package bwh

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
	var bwhs []BWH
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
	err := json.NewDecoder(r.Body).Decode(&bwhs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, bwh := range bwhs {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, bwh)
		if err != nil {
			log.Fatal(fmt.Sprintf("insert error: %v", err))
		}
	}
	err = json.NewEncoder(w).Encode(&bwhs)
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
	var bwhs []BWH
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
		err := h.svc.DB.SelectContext(r.Context(), &bwhs, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("select error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&bwhs)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
		err := h.svc.DB.SelectContext(r.Context(), &bwhs, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("select error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&bwhs)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in error: %v", err), queryIDs)
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &bwhs, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&bwhs)
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
	var bwhs []BWH
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
	err := json.NewDecoder(r.Body).Decode(&bwhs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, bwh := range bwhs {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, bwh)
		if err != nil {
			log.Fatal(fmt.Sprintf("update error: %v", err))
		}
	}
	err = json.NewEncoder(w).Encode(&bwhs)
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
			bwh
		WHERE
			entry_id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
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
			log.Fatal(fmt.Sprintf("delete error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
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
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
	}
}
