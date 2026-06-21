package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"booking-backend/internal/api"
	"booking-backend/internal/domain"
)

func TestSeed(t *testing.T) {
	s := New()

	assert.NotEmpty(t, s.GetOwner().Id)
	active := s.ListEventTypes(false)
	require.Len(t, active, 1)
	assert.Equal(t, "30-минутный звонок", active[0].Title)
}

func TestCreateBooking_OK(t *testing.T) {
	s := New()
	start := time.Now().AddDate(0, 0, 1)
	end := start.Add(30 * time.Minute)

	b, err := s.CreateBooking("et1", "Иван", "ivan@example.com", start, end)
	require.NoError(t, err)
	assert.NotEmpty(t, b.Id)
	assert.Len(t, s.ListBookings(), 1)
}

func TestCreateBooking_SlotTaken(t *testing.T) {
	s := New()
	start := time.Now().AddDate(0, 0, 1)
	end := start.Add(30 * time.Minute)

	_, err := s.CreateBooking("et1", "Иван", "ivan@example.com", start, end)
	require.NoError(t, err)

	// Пересекающаяся бронь (другой тип события) → занято.
	_, err = s.CreateBooking("et2", "Пётр", "petr@example.com", start, end)
	assert.ErrorIs(t, err, domain.ErrSlotTaken)
	assert.Len(t, s.ListBookings(), 1)
}

func TestGetEventType_OK(t *testing.T) {
	s := New()
	active := s.ListEventTypes(false)
	require.Len(t, active, 1)

	et, err := s.GetEventType(active[0].Id)
	require.NoError(t, err)
	assert.Equal(t, active[0].Title, et.Title)
}

func TestGetEventType_NotFound(t *testing.T) {
	s := New()
	_, err := s.GetEventType("nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUpdateOwner(t *testing.T) {
	s := New()
	originalID := s.GetOwner().Id

	updated := s.UpdateOwner(api.Owner{
		Name:  "Новое имя",
		Email: "new@example.com",
	})

	assert.Equal(t, originalID, updated.Id)
	assert.Equal(t, "Новое имя", updated.Name)
	assert.Equal(t, "new@example.com", updated.Email)

	// Проверяем, что хранилище обновилось
	current := s.GetOwner()
	assert.Equal(t, "Новое имя", current.Name)
}

func TestUpdateEventType_OK(t *testing.T) {
	s := New()
	active := s.ListEventTypes(true)
	require.Len(t, active, 1)

	updated, err := s.UpdateEventType(active[0].Id, api.CreateEventTypeRequest{
		Title: "Новый заголовок", Description: "Новое описание", Duration: 60,
	})
	require.NoError(t, err)
	assert.Equal(t, "Новый заголовок", updated.Title)
	assert.Equal(t, "Новое описание", updated.Description)
	assert.Equal(t, int32(60), updated.Duration)
}

func TestUpdateEventType_NotFound(t *testing.T) {
	s := New()
	_, err := s.UpdateEventType("nonexistent", api.CreateEventTypeRequest{})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDisableEventType_NotFound(t *testing.T) {
	s := New()
	err := s.DisableEventType("nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDisableEventType_FiltersGuest(t *testing.T) {
	s := New()
	created := s.CreateEventType(api.CreateEventTypeRequest{Title: "Тест", Description: "d", Duration: 30})

	require.NoError(t, s.DisableEventType(created.Id))

	// Гостю — без отключённого.
	for _, et := range s.ListEventTypes(false) {
		assert.NotEqual(t, created.Id, et.Id)
	}
	// Админу — со всеми.
	var found bool
	for _, et := range s.ListEventTypes(true) {
		if et.Id == created.Id {
			found = true
			assert.True(t, et.IsDisabled)
		}
	}
	assert.True(t, found)
}
