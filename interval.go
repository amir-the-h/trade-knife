package trade_knife

import (
	"time"
)

// Interval is the timeframe concept and determines duration of each candle.
type Interval string

// Returns actual duration of the interval.
func (i Interval) Duration() time.Duration {
	switch i {
	case Interval1m:
		return time.Minute
	case Interval3m:
		return time.Minute * 3
	case Interval5m:
		return time.Minute * 5
	case Interval15m:
		return time.Minute * 15
	case Interval30m:
		return time.Minute * 30
	case Interval1h:
		return time.Hour
	case Interval2h:
		return time.Hour * 2
	case Interval4h:
		return time.Hour * 4
	case Interval6h:
		return time.Hour * 6
	case Interval8h:
		return time.Hour * 8
	case Interval12h:
		return time.Hour * 12
	case Interval1d:
		return time.Hour * 24
	case Interval3d:
		return time.Hour * 24 * 3
	case Interval1w:
		return time.Hour * 24 * 7
	case Interval1M:
		return time.Hour * 24 * 30
	}

	return time.Duration(0)
}

// Returns the open and close time which the interval can fit in.
func (i Interval) GetPeriod(ts int64) (ot *time.Time, ct *time.Time, err error) {
	t := time.Unix(int64(ts/1000), 0)
	switch i {
	case Interval1m:
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	case Interval3m:
		m := t.Minute() - (t.Minute() % 3)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, 0, 0, time.UTC)
	case Interval5m:
		m := t.Minute() - (t.Minute() % 5)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, 0, 0, time.UTC)
	case Interval15m:
		m := t.Minute() - (t.Minute() % 15)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, 0, 0, time.UTC)
	case Interval30m:
		m := t.Minute() - (t.Minute() % 30)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, 0, 0, time.UTC)
	case Interval1h:
		*ot = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	case Interval2h:
		h := t.Hour() - (t.Hour() % 2)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), h, 0, 0, 0, time.UTC)
	case Interval4h:
		h := t.Hour() - (t.Hour() % 4)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), h, 0, 0, 0, time.UTC)
	case Interval6h:
		h := t.Hour() - (t.Hour() % 6)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), h, 0, 0, 0, time.UTC)
	case Interval8h:
		h := t.Hour() - (t.Hour() % 8)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), h, 0, 0, 0, time.UTC)
	case Interval12h:
		h := t.Hour() - (t.Hour() % 12)
		*ot = time.Date(t.Year(), t.Month(), t.Day(), h, 0, 0, 0, time.UTC)
	case Interval1d:
		*ot = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	case Interval3d:
		d := t.Hour() - (t.Hour() % 3)
		*ot = time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, time.UTC)
	case Interval1w:
		d := t.Hour() - (t.Hour() % 7)
		*ot = time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, time.UTC)
	case Interval1M:
		*ot = time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.UTC)
	}

	*ct = ot.Add(i.Duration())

	return
}
