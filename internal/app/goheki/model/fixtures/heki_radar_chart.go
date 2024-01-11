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

func NewHekiRadarChart(ctx context.Context) *ModelConnector {
	hekiRadarChart := &HekiRadarChart{
		AI: 1,
		NU: 1,
	}

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
			result, err := f.dbv1.NamedExecContext(ctx, "INSERT INTO heki_radar_chart (entry_id, ai, nu) VALUES (:entry_id, :ai, :nu)", hekiRadarChart)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			hekiRadarChart.EntryID = &id
		},
	}
}
