package slots

import (
	"time"

	"booking-backend/internal/api"
)

// Params задаёт параметры генерации слотов.
type Params struct {
	WorkStartHour int
	WorkEndHour   int
	IntervalMin   int
}

// Generate строит слоты на указанную дату (в локальной зоне date).
//
//   - окно WorkStartHour..WorkEndHour с шагом IntervalMin;
//   - для сегодняшней даты прошедшие слоты исключаются;
//   - каждый слот помечается isBusy при пересечении с любой бронью.
func Generate(date time.Time, now time.Time, bookings []api.Booking, p Params) []api.Slot {
	loc := date.Location()
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

	interval := time.Duration(p.IntervalMin) * time.Minute
	result := make([]api.Slot, 0)

	for h := p.WorkStartHour; h < p.WorkEndHour; h++ {
		for m := 0; m < 60; m += p.IntervalMin {
			start := dayStart.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
			end := start.Add(interval)

			// Скрыть прошедшие слоты
			if !start.After(now) {
				continue
			}

			isBusy := false
			for _, b := range bookings {
				if start.Before(b.EndTime) && b.StartTime.Before(end) {
					isBusy = true
					break
				}
			}

			result = append(result, api.Slot{
				StartTime: start.UTC(),
				EndTime:   end.UTC(),
				IsBusy:    isBusy,
			})
		}
	}

	return result
}
