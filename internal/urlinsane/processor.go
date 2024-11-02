// Copyright (C) 2024 Rangertaha
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package urlinsane

// import (
// 	"time"


// )

// type processor struct {
// 	// metrics   chan<- interfaces.Metric
// 	precision time.Duration
// 	startTime time.Time
// 	endTime   time.Time
// 	nowTime   time.Time
// }

// func NewBackfillAccumulator(metrics chan<- interfaces.Metric, precision time.Duration, start, stop time.Time) interfaces.Accumulator {
// 	return &processor{
// 		metrics:   metrics,
// 		precision: precision,
// 		startTime: start,
// 		nowTime:   start,
// 		endTime:   time.Now().Round(precision),
// 	}
// }

// func NewAccumulator(metrics chan<- interfaces.Metric, precision time.Duration) interfaces.Accumulator {
// 	return &processor{
// 		metrics:   metrics,
// 		precision: precision,
// 	}
// }

// func (ac *processor) AddFields(
// 	measurement string,
// 	fields map[string]interface{},
// 	tags map[string]string,
// 	t ...time.Time,
// ) {
// 	ac.addMeasurement(measurement, tags, fields, interfaces.Untyped, t...)
// }

// // func (ac *processor) AddGauge(
// // 	measurement string,
// // 	fields map[string]interface{},
// // 	tags map[string]string,
// // 	t ...time.Time,
// // ) {
// // 	ac.addMeasurement(measurement, tags, fields, interfaces.Gauge, t...)
// // }

// // func (ac *processor) AddCounter(
// // 	measurement string,
// // 	fields map[string]interface{},
// // 	tags map[string]string,
// // 	t ...time.Time,
// // ) {
// // 	ac.addMeasurement(measurement, tags, fields, interfaces.Counter, t...)
// // }

// // func (ac *processor) AddSummary(
// // 	measurement string,
// // 	fields map[string]interface{},
// // 	tags map[string]string,
// // 	t ...time.Time,
// // ) {
// // 	ac.addMeasurement(measurement, tags, fields, interfaces.Summary, t...)
// // }

// // func (ac *processor) AddHistogram(
// // 	measurement string,
// // 	fields map[string]interface{},
// // 	tags map[string]string,
// // 	t ...time.Time,
// // ) {
// // 	ac.addMeasurement(measurement, tags, fields, interfaces.Histogram, t...)
// // }

// func (ac *processor) AddMetric(m interfaces.Metric) {
// 	if m != nil {
// 		m.SetTime(m.Time().Round(ac.precision))
// 		if ac.startTime.IsZero() {
// 			ac.startTime = m.Time()
// 		}

// 		ac.nowTime = m.Time()
// 		ac.metrics <- m

// 	}

// 	// if m := ac.maker.MakeMetric(m); m != nil {
// 	// 	ac.metrics <- m
// 	// }
// }

// func (ac *processor) StartTime() time.Time {
// 	return ac.startTime
// }

// func (ac *processor) SetStartTime(time time.Time) {
// 	ac.startTime = time.Round(ac.precision)
// }

// func (ac *processor) EndTime() time.Time {
// 	return ac.endTime
// }

// func (ac *processor) SetEndTime(time time.Time) {
// 	ac.endTime = time.Round(ac.precision)
// }

// func (ac *processor) SetTime(time time.Time) {
// 	ac.nowTime = time.Round(ac.precision)
// }

// func (ac *processor) Time() time.Time {
// 	return ac.nowTime
// }

// func (a *processor) Ticker(span time.Duration) (m <-chan time.Time) {
// 	if a.valid() {
// 		out := make(chan time.Time, 1)
// 		go func() {
// 			for d := a.startTime; !d.After(a.endTime); d = d.Add(span) {
// 				out <- d
// 				a.nowTime = d
// 			}
// 			defer close(out)
// 		}()
// 		return out
// 	}

// 	return
// }

// func (a *processor) TickerBar(span time.Duration) (m <-chan time.Time) {
// 	max := a.endTime.Round(span).Unix() - a.startTime.Round(span).Unix()
// 	bar := progressbar.Default(max)

// 	out := make(chan time.Time, 1)
// 	go func() {
// 		// fmt.Println("START", a.startTime)
// 		for d := a.startTime; !d.After(a.endTime); d = d.Add(span) {
// 			out <- d
// 			a.nowTime = d

// 			progress := max - (a.endTime.Round(span).Unix() - a.nowTime.Unix())
// 			bar.Set(int(progress))

// 		}
// 		// fmt.Println("END", a.endTime)
// 		bar.Finish()
// 		// fmt.Println("DONE")
// 		bar.Exit()
// 		defer close(out)
// 	}()
// 	return out
// }

// // func (a *processor) Minutes() (m <-chan time.Time) {
// // 	if a.valid() {
// // 		out := make(chan time.Time, 1)
// // 		go func() {
// // 			for d := a.startTime; !d.After(a.endTime); d = d.Add(1 * time.Minute) {
// // 				out <- d
// // 				a.nowTime = d
// // 			}
// // 			defer close(out)
// // 		}()
// // 		return out
// // 	}

// // 	return
// // }

// // func (t *processor) Hours() (m <-chan time.Time) {
// // 	if t.valid() {
// // 		out := make(chan time.Time, 1)
// // 		go func() {
// // 			for d := t.startTime; !d.After(t.endTime); d = d.Add(1 * time.Hour) {
// // 				out <- d
// // 				t.nowTime = d
// // 			}
// // 			defer close(out)
// // 		}()
// // 		return out
// // 	}

// // 	return
// // }

// // func (t *processor) Days() (m <-chan time.Time) {
// // 	if t.valid() {
// // 		out := make(chan time.Time, 1)
// // 		go func() {
// // 			for d := t.startTime; !d.After(t.endTime); d = d.AddDate(0, 0, 1) {
// // 				out <- d
// // 				t.nowTime = d
// // 			}
// // 			defer close(out)
// // 		}()
// // 		return out
// // 	}

// // 	return
// // }

// func (a *processor) valid() bool {
// 	if a.startTime.IsZero() || a.endTime.IsZero() {
// 		return false
// 	}

// 	if a.startTime.Before(time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC)) || a.endTime.Before(time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC)) {
// 		return false
// 	}

// 	return true
// }

// // func (ac *processor) TotalHours() int64 {

// // 	start := ac.startTime

// // 	// If the start date is before year 2000 assume it not correct
// // 	if ac.startTime.Before(time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)) {
// // 		start = time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)
// // 	}

// // 	return int64(time.Since(start).Hours())
// // }

// // func (ac *processor) ShowProgress() {
// // 	max := ac.endTime.Unix() - ac.startTime.Unix()
// // 	bar := progressbar.Default(max)

// // 	for ac.nowTime.Compare(ac.endTime) < 0 {
// // 		progress := max - (ac.endTime.Unix() - ac.nowTime.Unix())
// // 		// fmt.Println(progress)
// // 		bar.Set(int(progress))

// // 	}
// // 	bar.Finish()
// // }

// // func (ac *processor) AddOrder(o interfaces.Order) {
// // 	if o != nil {
// // 		ac.metrics <- o.Metric()
// // 	}

// // 	// if m := ac.maker.MakeMetric(m); m != nil {
// // 	// 	ac.metrics <- m
// // 	// }
// // }

// func (ac *processor) addMeasurement(
// 	measurement string,
// 	tags map[string]string,
// 	fields map[string]interface{},
// 	tp interfaces.ValueType,
// 	t ...time.Time,
// ) {
// 	m := metric.New(measurement, tags, fields, ac.getTime(t), tp)
// 	ac.metrics <- m
// 	// if m := ac.maker.MakeMetric(m); m != nil {
// 	// 	ac.metrics <- m
// 	// }
// }

// // AddError passes a runtime error to the processor.
// // The error will be tagged with the plugin name and written to the log.
// func (ac *processor) AddError(err error) {
// 	if err == nil {
// 		return
// 	}
// 	// ac.maker.Log().Errorf("Error in plugin: %v", err)
// }

// func (ac *processor) SetPrecision(precision time.Duration) {
// 	ac.precision = precision
// }

// func (ac *processor) getTime(t []time.Time) time.Time {
// 	var timestamp time.Time
// 	if len(t) > 0 {
// 		timestamp = t[0]
// 	} else {
// 		timestamp = time.Now()
// 	}
// 	return timestamp.Round(ac.precision)
// }

// // func (ac *processor) WithTracking(maxTracked int) interfaces.TrackingAccumulator {
// // 	return &trackingAccumulator{
// // 		Accumulator: ac,
// // 		delivered:   make(chan interfaces.DeliveryInfo, maxTracked),
// // 	}
// // }

// // type trackingAccumulator struct {
// // 	interfaces.Accumulator
// // 	delivered chan interfaces.DeliveryInfo
// // }

// // func (a *trackingAccumulator) AddTrackingMetric(m string) interfaces.TrackingID {
// // 	dm, id := metric.WithTracking(m, a.onDelivery)
// // 	a.AddMetric(dm)
// // 	return id
// // }

// // func (a *trackingAccumulator) AddTrackingMetricGroup(group []string) interfaces.TrackingID {
// // 	db, id := metric.WithGroupTracking(group, a.onDelivery)
// // 	for _, m := range db {
// // 		a.AddMetric(m)
// // 	}
// // 	return id
// // }

// // func (a *trackingAccumulator) Delivered() <-chan interfaces.DeliveryInfo {
// // 	return a.delivered
// // }

// // func (a *trackingAccumulator) onDelivery(info interfaces.DeliveryInfo) {
// // 	select {
// // 	case a.delivered <- info:
// // 	default:
// // 		// This is a programming error in the input.  More items were sent for
// // 		// tracking than space requested.
// // 		panic("channel is full")
// // 	}
// // }
