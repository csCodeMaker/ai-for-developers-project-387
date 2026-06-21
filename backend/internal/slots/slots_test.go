package slots

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"booking-backend/internal/api"
)

var defaultParams = Params{WorkStartHour: 8, WorkEndHour: 20, IntervalMin: 30}

func TestGenerate_Count(t *testing.T) {
	// Дата завтра, now — сейчас → все слоты в будущем.
	date := time.Now().AddDate(0, 0, 1)
	now := time.Now()

	got := Generate(date, now, nil, defaultParams)

	// Окно 08:00–20:00 с шагом 30 мин = 24 слота.
	assert.Len(t, got, 24)
}

func TestGenerate_HidesPastSlots(t *testing.T) {
	loc := time.Local
	// Сегодня в 12:30 — слоты до 12:30 должны скрыться.
	now := time.Date(2030, 1, 10, 12, 30, 0, 0, loc)
	date := time.Date(2030, 1, 10, 0, 0, 0, 0, loc)

	got := Generate(date, now, nil, defaultParams)

	for _, s := range got {
		assert.True(t, s.StartTime.After(now), "слот %v должен быть в будущем", s.StartTime)
	}
	// С 13:00 до 19:30 включительно = 14 слотов.
	assert.Len(t, got, 14)
}

func TestGenerate_MarksBusy(t *testing.T) {
	loc := time.Local
	now := time.Date(2030, 1, 10, 0, 0, 0, 0, loc)
	date := time.Date(2030, 1, 11, 0, 0, 0, 0, loc)

	// Бронь на 10:00–10:30.
	bStart := time.Date(2030, 1, 11, 10, 0, 0, 0, loc)
	bookings := []api.Booking{{StartTime: bStart, EndTime: bStart.Add(30 * time.Minute)}}

	got := Generate(date, now, bookings, defaultParams)

	var busyCount int
	for _, s := range got {
		if s.IsBusy {
			busyCount++
			assert.True(t, s.StartTime.Equal(bStart.UTC()))
		}
	}
	assert.Equal(t, 1, busyCount, "ровно один слот должен быть занят")
}
