package fixtures

import (
	"context"
	"testing"
)

type HekiRadarChart struct {
	EntryID *int64 `db:"entry_id"`
	AI      int64  `db:"ai"`
	NU      int64  `db:"nu"`
}

func NewHekiRadarChart(ctx context.Context, setter func(h *HekiRadarChart)) *ModelConnector {
	hekiRadarChart := &HekiRadarChart{
		AI: 1,
		NU: 1,
	}

	setter(hekiRadarChart)

	return &ModelConnector{
		Model: hekiRadarChart,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HekiRadarCharts = append(f.HekiRadarCharts, hekiRadarChart)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				hekiRadarChart.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hekiRadarChart)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			// 連番されるIDをセットする
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO heki_radar_chart (
					entry_id,
					ai,
					nu
				) VALUES (
					$1,
					$2,
					$3
				) RETURNING entry_id`,
				hekiRadarChart.EntryID,
				hekiRadarChart.AI,
				hekiRadarChart.NU,
			).Scan(&hekiRadarChart.EntryID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
