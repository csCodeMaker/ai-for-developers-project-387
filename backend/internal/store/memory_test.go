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

func TestGetOwner(t *testing.T) {
	s := New()
	o := s.GetOwner()
	assert.Equal(t, "Владелец календаря", o.Name)
	assert.Equal(t, "owner@example.com", o.Email)
	assert.NotEmpty(t, o.Id)
}

func TestUpdateOwner(t *testing.T) {
	s := New()
	originalID := s.GetOwner().Id

	updated := s.UpdateOwner(api.Owner{
		Name: "Новый владелец", Email: "new@example.com",
		Description: "desc", TimeZone: "UTC",
	})

	assert.Equal(t, originalID, updated.Id, "ID должен сохраняться")
	assert.Equal(t, "Новый владелец", updated.Name)
	assert.Equal(t, "new@example.com", updated.Email)

	o := s.GetOwner()
	assert.Equal(t, "Новый владелец", o.Name)
}

func TestGetEventType_Found(t *testing.T) {
	s := New()
	active := s.ListEventTypes(false)
	require.Len(t, active, 1)

	et, err := s.GetEventType(active[0].Id)
	require.NoError(t, err)
	assert.Equal(t, active[0].Id, et.Id)
}

func TestGetEventType_NotFound(t *testing.T) {
	s := New()
	_, err := s.GetEventType("nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCreateEventType(t *testing.T) {
	s := New()
	et := s.CreateEventType(api.CreateEventTypeRequest{
		Title: "Новый тип", Description: "Описание", Duration: 60,
	})

	assert.NotEmpty(t, et.Id)
	assert.Equal(t, "Новый тип", et.Title)
	assert.Equal(t, int32(60), et.Duration)
	assert.False(t, et.IsDisabled)

	list := s.ListEventTypes(false)
	assert.Len(t, list, 2)
}

func TestUpdateEventType(t *testing.T) {
	s := New()
	active := s.ListEventTypes(false)
	originalID := active[0].Id

	et, err := s.UpdateEventType(originalID, api.CreateEventTypeRequest{
		Title: "Изменённый", Description: "Новое описание", Duration: 45,
	})
	require.NoError(t, err)
	assert.Equal(t, originalID, et.Id)
	assert.Equal(t, "Изменённый", et.Title)
	assert.Equal(t, int32(45), et.Duration)

	got, _ := s.GetEventType(originalID)
	assert.Equal(t, "Изменённый", got.Title)
}

func TestUpdateEventType_NotFound(t *testing.T) {
	s := New()
	_, err := s.UpdateEventType("nonexistent", api.CreateEventTypeRequest{
		Title: "X", Description: "d", Duration: 30,
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDisableEventType(t *testing.T) {
	s := New()
	created := s.CreateEventType(api.CreateEventTypeRequest{Title: "Тест", Description: "d", Duration: 30})
	require.NoError(t, s.DisableEventType(created.Id))

	et, _ := s.GetEventType(created.Id)
	assert.True(t, et.IsDisabled)
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

	for _, et := range s.ListEventTypes(false) {
		assert.NotEqual(t, created.Id, et.Id)
	}
	var found bool
	for _, et := range s.ListEventTypes(true) {
		if et.Id == created.Id {
			found = true
			assert.True(t, et.IsDisabled)
		}
	}
	assert.True(t, found)
}

func TestListEventTypes_IncludeDisabled(t *testing.T) {
	s := New()
	created := s.CreateEventType(api.CreateEventTypeRequest{Title: "D", Description: "d", Duration: 15})
	require.NoError(t, s.DisableEventType(created.Id))

	all := s.ListEventTypes(true)
	require.Len(t, all, 2)
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

	_, err = s.CreateBooking("et2", "Пётр", "petr@example.com", start, end)
	assert.ErrorIs(t, err, domain.ErrSlotTaken)
	assert.Len(t, s.ListBookings(), 1)
}

func TestCreateBooking_PartialOverlap(t *testing.T) {
	s := New()
	start := time.Now().AddDate(0, 0, 1).Truncate(time.Hour)
	end := start.Add(1 * time.Hour)

	_, err := s.CreateBooking("et1", "A", "a@a.com", start, end)
	require.NoError(t, err)

	overlapStart := start.Add(30 * time.Minute)
	overlapEnd := start.Add(90 * time.Minute)
	_, err = s.CreateBooking("et2", "B", "b@b.com", overlapStart, overlapEnd)
	assert.ErrorIs(t, err, domain.ErrSlotTaken)
}

func TestCreateBooking_NonOverlapping(t *testing.T) {
	s := New()
	base := time.Now().AddDate(0, 0, 1).Truncate(time.Hour)

	_, err := s.CreateBooking("et1", "A", "a@a.com", base, base.Add(30*time.Minute))
	require.NoError(t, err)

	_, err = s.CreateBooking("et2", "B", "b@b.com", base.Add(30*time.Minute), base.Add(1*time.Hour))
	require.NoError(t, err)

	assert.Len(t, s.ListBookings(), 2)
}

func TestListBookings_Empty(t *testing.T) {
	s := New()
	assert.Empty(t, s.ListBookings())
}

func TestListBookings_Multiple(t *testing.T) {
	s := New()
	start := time.Now().AddDate(0, 0, 1)
	for i := 0; i < 3; i++ {
		s.CreateBooking("et1", "A", "a@a.com", start.Add(time.Duration(i)*time.Hour), start.Add(time.Duration(i)*time.Hour+30*time.Minute))
	}
	assert.Len(t, s.ListBookings(), 3)
}
