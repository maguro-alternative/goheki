package hekiradarchart

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
	// POST以外は受け付けない
	if r.Method != http.MethodPost {
		return
	}
	var HekiRadarChartsJson HekiRadarChartsJson
	query := `
		INSERT INTO heki_radar_chart (
			entry_id,
			ai,
			nu
		) VALUES (
			:entry_id,
			:ai,
			:nu
		)
	`
	err := json.NewDecoder(r.Body).Decode(&HekiRadarChartsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = HekiRadarChartsJson.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hrc := range HekiRadarChartsJson.HekiRadarCharts {
		// jsonバリデーション
		err = hrc.Validate()
		if err != nil {
			log.Println(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hrc)
		if err != nil {
			log.Println(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&HekiRadarChartsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
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
	// GET以外は受け付けない
	if r.Method != http.MethodGet {
		return
	}
	var hekiRadarChartsJson HekiRadarChartsJson
	query := `
		SELECT
			entry_id,
			ai,
			nu
		FROM
			heki_radar_chart
		WHERE
			entry_id IN (?)
	`
	queryIDs, ok := r.URL.Query()["entry_id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				ai,
				nu
			FROM
				heki_radar_chart
		`
		err := h.svc.DB.SelectContext(r.Context(), &hekiRadarChartsJson.HekiRadarCharts, query)
		if err != nil {
			log.Println(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// json書き込み
		err = json.NewEncoder(w).Encode(&hekiRadarChartsJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				ai,
				nu
			FROM
				heki_radar_chart
			WHERE
				entry_id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hekiRadarChartsJson.HekiRadarCharts, query, queryIDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// json書き込み
		err = json.NewEncoder(w).Encode(&hekiRadarChartsJson)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Println(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hekiRadarChartsJson.HekiRadarCharts, query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&hekiRadarChartsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	// PUT以外は受け付けない
	if r.Method != http.MethodPut {
		return
	}
	var hekiRadarChartsJson HekiRadarChartsJson
	query := `
		UPDATE
			heki_radar_chart
		SET
			ai = :ai,
			nu = :nu
		WHERE
			entry_id = :entry_id
	`
	err := json.NewDecoder(r.Body).Decode(&hekiRadarChartsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = hekiRadarChartsJson.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	for _, hrc := range hekiRadarChartsJson.HekiRadarCharts {
		// jsonバリデーション
		err = hrc.Validate()
		if err != nil {
			log.Println(fmt.Sprintf("json validate error: %v", err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		_, err := h.svc.DB.NamedExecContext(r.Context(), query, hrc)
		if err != nil {
			log.Println(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&hekiRadarChartsJson)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
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
	// DELETE以外は受け付けない
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			heki_radar_chart
		WHERE
			entry_id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// jsonバリデーション
	err = delIDs.Validate()
	if err != nil {
		log.Println(fmt.Sprintf("json validate error: %v", err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				heki_radar_chart
			WHERE
				entry_id = $1
		`
		_, err := h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Println(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// json書き込み
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Println(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// idの数だけ置換文字を作成
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Println(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Postgresの場合は置換文字を$1, $2, ...とする必要がある
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Println(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json書き込み
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Println(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
