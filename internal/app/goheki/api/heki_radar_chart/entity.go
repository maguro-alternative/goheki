package hekiradarchart

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HekiRadarChart struct {
	EntryID int64 `db:"entry_id"`
	AI      int64 `db:"ai"`
	NU      int64 `db:"nu"`
}

func (h *HekiRadarChart) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.EntryID, validation.Required),
		validation.Field(&h.AI, validation.Required),
		validation.Field(&h.NU, validation.Required),
	)
}

type HekiRadarChartsJson struct {
	HekiRadarCharts []HekiRadarChart `json:"heki_radar_charts"`
}

func (h *HekiRadarChartsJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HekiRadarCharts, validation.Required),
	)
}

type IDs struct {
	IDs []int64 `json:"ids"`
}

func (i *IDs) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.IDs, validation.Required),
	)
}
