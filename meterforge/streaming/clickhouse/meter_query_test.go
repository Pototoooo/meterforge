package clickhouse

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/pkg/filter"
	"github.com/Pototoooo/meterforge/pkg/models"
)

func TestQueryMeter(t *testing.T) {
	subject := "subject1"
	from, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00.001Z")
	to, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")
	storedAtOffset, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00.001Z")
	tz, _ := time.LoadLocation("Asia/Shanghai")
	windowSize := meter.WindowSizeHour

	tests := []struct {
		name     string
		query    queryMeter
		wantSQL  string
		wantArgs []interface{}
	}{
		{
			name: "basic query",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterSubject: []string{subject},
				From:          &from,
				To:            &to,
				GroupBy:       []string{"subject", "group1", "group2"},
				WindowSize:    &windowSize,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalHour(1), 'UTC') AS windowstart, tumbleEnd(mf_events.time, toIntervalHour(1), 'UTC') AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, mf_events.subject, JSON_VALUE(mf_events.data, '$.group1') as group1, JSON_VALUE(mf_events.data, '$.group2') as group2 FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend, subject, group1, group2 ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1"}, from.Unix(), to.Unix()},
		},
		{
			name: "basic query with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterSubject:          []string{subject},
				From:                   &from,
				To:                     &to,
				GroupBy:                []string{"subject", "group1", "group2"},
				WindowSize:             &windowSize,
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalHour(1), 'UTC') AS windowstart, tumbleEnd(mf_events.time, toIntervalHour(1), 'UTC') AS windowend, sum(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19)) AS value, mf_events.subject, JSON_VALUE(mf_events.data, '$.group1') as group1, JSON_VALUE(mf_events.data, '$.group2') as group2 FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend, subject, group1, group2 ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1"}, from.Unix(), to.Unix()},
		},
		{
			name: "basic query with decimal stored at offset",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterStoredAt: &filter.FilterTimeUnix{
					FilterTime: filter.FilterTime{
						Lt: &storedAtOffset,
					},
				},
				FilterSubject:          []string{subject},
				From:                   &from,
				To:                     &to,
				GroupBy:                []string{"subject", "group1", "group2"},
				WindowSize:             &windowSize,
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalHour(1), 'UTC') AS windowstart, tumbleEnd(mf_events.time, toIntervalHour(1), 'UTC') AS windowend, sum(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19)) AS value, mf_events.subject, JSON_VALUE(mf_events.data, '$.group1') as group1, JSON_VALUE(mf_events.data, '$.group2') as group2 FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) AND mf_events.time >= ? AND mf_events.time < ? AND mf_events.stored_at < ? GROUP BY windowstart, windowend, subject, group1, group2 ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1"}, from.Unix(), to.Unix(), storedAtOffset.Unix()},
		},
		{
			name: "Aggregate all available data",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with count aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:         "meter1",
					EventType:   "event1",
					Aggregation: meter.MeterAggregationCount,
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, count(*) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with count aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:         "meter1",
					EventType:   "event1",
					Aggregation: meter.MeterAggregationCount,
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, count(*) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with unique count aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationUniqueCount,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, uniqExact(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null')) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with unique count aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationUniqueCount,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, uniqExact(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null')) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with AVG aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationAvg,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, avg(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with AVG aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationAvg,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, avg(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with MIN aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationMin,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, min(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with MIN aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationMin,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, min(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with MAX aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationMax,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, max(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with MAX aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationMax,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, max(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with LATEST aggregation",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationLatest,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, argMax(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null), mf_events.time) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate with LATEST aggregation with decimal precision",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationLatest,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnableDecimalPrecision: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, argMax(toDecimal128OrNull(nullIf(JSON_VALUE(mf_events.data, '$.value'), 'null'), 19), mf_events.time) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate data from start",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From: &from,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ?",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix()},
		},
		{
			name: "Aggregate data between period",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From: &from,
				To:   &to,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ? AND mf_events.time < ?",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix(), to.Unix()},
		},
		{
			name: "Aggregate data between period, groupped by window size",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From:       &from,
				To:         &to,
				WindowSize: &windowSize,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalHour(1), 'UTC') AS windowstart, tumbleEnd(mf_events.time, toIntervalHour(1), 'UTC') AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix(), to.Unix()},
		},
		{
			name: "Aggregate data between period in a different timezone, groupped by window size",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From:           &from,
				To:             &to,
				WindowSize:     &windowSize,
				WindowTimeZone: tz,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalHour(1), 'Asia/Shanghai') AS windowstart, tumbleEnd(mf_events.time, toIntervalHour(1), 'Asia/Shanghai') AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix(), to.Unix()},
		},
		{
			name: "Aggregate data between period, groupped by DAY window size",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From:       &from,
				To:         &to,
				WindowSize: lo.ToPtr(meter.WindowSizeDay),
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalDay(1), 'UTC') AS windowstart, windowstart + toIntervalDay(1) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix(), to.Unix()},
		},
		{
			name: "Aggregate data between period in a different timezone, groupped by DAY window size",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				From:           &from,
				To:             &to,
				WindowSize:     lo.ToPtr(meter.WindowSizeDay),
				WindowTimeZone: tz,
			},
			wantSQL:  "SELECT tumbleStart(mf_events.time, toIntervalDay(1), 'Asia/Shanghai') AS windowstart, windowstart + toIntervalDay(1) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.time >= ? AND mf_events.time < ? GROUP BY windowstart, windowend ORDER BY windowstart",
			wantArgs: []interface{}{"my_namespace", "event1", from.Unix(), to.Unix()},
		},
		{
			name: "Aggregate data for a single subject",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterSubject: []string{subject},
				GroupBy:       []string{"subject"},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, mf_events.subject FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) GROUP BY subject",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1"}},
		},
		{
			name: "Aggregate data for a single subject and group by additional fields",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterSubject: []string{subject},
				GroupBy:       []string{"subject", "group1", "group2"},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, mf_events.subject, JSON_VALUE(mf_events.data, '$.group1') as group1, JSON_VALUE(mf_events.data, '$.group2') as group2 FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) GROUP BY subject, group1, group2",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1"}},
		},
		{
			name: "Aggregate data for a multiple subjects",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				FilterSubject: []string{subject, "subject2"},
				GroupBy:       []string{"subject"},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, mf_events.subject FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) GROUP BY subject",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1", "subject2"}},
		},
		{
			name: "Select customer ID",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
				},
				FilterCustomer: []streaming.Customer{
					customer.Customer{
						ManagedResource: models.ManagedResource{
							NamespacedModel: models.NamespacedModel{
								Namespace: "my_namespace",
							},
							ID: "customer1",
						},
						UsageAttribution: &customer.CustomerUsageAttribution{
							SubjectKeys: []string{"subject1"},
						},
					},
					customer.Customer{
						ManagedResource: models.ManagedResource{
							NamespacedModel: models.NamespacedModel{
								Namespace: "my_namespace",
							},
							ID: "customer2",
						},
						UsageAttribution: &customer.CustomerUsageAttribution{
							SubjectKeys: []string{"subject2"},
						},
					},
				},
				GroupBy: []string{"customer_id"},
			},
			wantSQL:  "WITH map('subject1', 'customer1', 'subject2', 'customer2') as subject_to_customer_id SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, subject_to_customer_id[mf_events.subject] AS customer_id FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) GROUP BY customer_id",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1", "subject2"}},
		},
		{
			name: "Filter by customer ID without group by",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
				},
				FilterCustomer: []streaming.Customer{
					customer.Customer{
						ManagedResource: models.ManagedResource{
							NamespacedModel: models.NamespacedModel{
								Namespace: "my_namespace",
							},
							ID: "customer1",
						},
						Key: lo.ToPtr("customer-key-1"),
						UsageAttribution: &customer.CustomerUsageAttribution{
							SubjectKeys: []string{"subject1"},
						},
					},
					customer.Customer{
						ManagedResource: models.ManagedResource{
							NamespacedModel: models.NamespacedModel{
								Namespace: "my_namespace",
							},
							ID: "customer2",
						},
						UsageAttribution: &customer.CustomerUsageAttribution{
							SubjectKeys: []string{"subject2"},
						},
					},
				},
			},
			wantSQL: "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?)",
			wantArgs: []interface{}{"my_namespace", "event1", []string{
				// Only the first customer has a key
				"customer-key-1",
				// Usage attribution subjects of the first customer
				"subject1",
				// Usage attribution subjects of the second customer
				"subject2",
			}},
		},
		{ // Filter by both customer and subject
			name: "Filter by both customer and subject",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
				},
				FilterCustomer: []streaming.Customer{
					customer.Customer{
						ManagedResource: models.ManagedResource{
							NamespacedModel: models.NamespacedModel{
								Namespace: "my_namespace",
							},
							ID: "customer1",
						},
						UsageAttribution: &customer.CustomerUsageAttribution{
							SubjectKeys: []string{"subject1", "subject2"},
						},
					},
				},
				FilterSubject: []string{"subject1"},
				GroupBy:       []string{"customer_id"},
			},
			wantSQL:  "WITH map('subject1', 'customer1', 'subject2', 'customer1') as subject_to_customer_id SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value, subject_to_customer_id[mf_events.subject] AS customer_id FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND mf_events.subject IN (?) AND mf_events.subject IN (?) GROUP BY customer_id",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"subject1", "subject2"}, []string{"subject1"}},
		},
		{
			name: "Aggregate data with filtering for a single group and single value",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"g1": "$.group1",
						"g2": "$.group2",
					},
				},
				FilterGroupBy: map[string]filter.FilterString{"g1": {Eq: lo.ToPtr("g1v1")}},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND JSON_VALUE(mf_events.data, '$.group1') = ?",
			wantArgs: []interface{}{"my_namespace", "event1", "g1v1"},
		},
		{
			name: "Aggregate data with filtering for a single group and multiple values",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"g1": "$.group1",
						"g2": "$.group2",
					},
				},
				FilterGroupBy: map[string]filter.FilterString{"g1": {In: lo.ToPtr([]string{"g1v1", "g1v2"})}},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND JSON_VALUE(mf_events.data, '$.group1') IN (?)",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"g1v1", "g1v2"}},
		},
		{
			name: "Aggregate data with filtering for multiple groups and multiple values",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"g1": "$.group1",
						"g2": "$.group2",
					},
				},
				FilterGroupBy: map[string]filter.FilterString{
					"g1": {In: lo.ToPtr([]string{"g1v1", "g1v2"})},
					"g2": {In: lo.ToPtr([]string{"g2v1", "g2v2"})},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? AND JSON_VALUE(mf_events.data, '$.group1') IN (?) AND JSON_VALUE(mf_events.data, '$.group2') IN (?)",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"g1v1", "g1v2"}, []string{"g2v1", "g2v2"}},
		},
		{
			name: "Aggregate all available data, prewhere enabled (should not move anything to prewhere)",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"group1": "$.group1",
						"group2": "$.group2",
					},
				},
				EnablePrewhere: true,
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ?",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
		{
			name: "Aggregate data with with filtering for multiple groups and multiple values prewhere enabled",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				EnablePrewhere:  true,
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
					GroupBy: map[string]string{
						"g1": "$.group1",
						"g2": "$.group2",
					},
				},
				FilterGroupBy: map[string]filter.FilterString{
					"g1": {In: lo.ToPtr([]string{"g1v1", "g1v2"})},
					"g2": {In: lo.ToPtr([]string{"g2v1", "g2v2"})},
				},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events PREWHERE mf_events.namespace = ? AND mf_events.type = ? WHERE JSON_VALUE(mf_events.data, '$.group1') IN (?) AND JSON_VALUE(mf_events.data, '$.group2') IN (?) SETTINGS optimize_move_to_prewhere = 1, allow_reorder_prewhere_conditions = 1",
			wantArgs: []interface{}{"my_namespace", "event1", []string{"g1v1", "g1v2"}, []string{"g2v1", "g2v2"}},
		},
		{
			name: "Add query settings",
			query: queryMeter{
				Database:        "meterforge",
				EventsTableName: "mf_events",
				Namespace:       "my_namespace",
				Meter: meter.Meter{
					Key:           "meter1",
					EventType:     "event1",
					Aggregation:   meter.MeterAggregationSum,
					ValueProperty: lo.ToPtr("$.value"),
				},
				QuerySettings: map[string]string{"foo": "1"},
			},
			wantSQL:  "SELECT tumbleStart(min(mf_events.time), toIntervalMinute(1)) AS windowstart, tumbleEnd(max(mf_events.time), toIntervalMinute(1)) AS windowend, sum(ifNotFinite(toFloat64OrNull(JSON_VALUE(mf_events.data, '$.value')), null)) AS value FROM meterforge.mf_events WHERE mf_events.namespace = ? AND mf_events.type = ? SETTINGS foo = 1",
			wantArgs: []interface{}{"my_namespace", "event1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, gotArgs, err := tt.query.toSQL()
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tt.wantSQL, gotSql)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}
