package entry

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	//"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type ReadHandler struct {
	svc *service.IndexService
}

func NewReadHandler(svc *service.IndexService) *ReadHandler {
	return &ReadHandler{
		svc: svc,
	}
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	//err := json.NewDecoder(r.Body).Decode(&entrys)
	//if err != nil {
	//log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	//}
	//query, args, err := db.In(query, entrys)
	//if err != nil {
	//return
	//}
	err := h.svc.DB.SelectContext(r.Context(), &entrys, query)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&entrys)
}
