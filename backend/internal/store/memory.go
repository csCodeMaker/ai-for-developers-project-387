package store

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"booking-backend/internal/api"
	"booking-backend/internal/domain"
)

// MemoryStore — потокобезопасное in-memory хранилище.
// Данные сбрасываются при перезапуске сервиса.
type MemoryStore struct {
	mu         sync.RWMutex
	owner      api.Owner
	eventTypes map[string]api.EventType
	bookings   map[string]api.Booking
}

// New создаёт хранилище с сидом: один владелец и дефолтный тип события.
func New() *MemoryStore {
	s := &MemoryStore{
		eventTypes: make(map[string]api.EventType),
		bookings:   make(map[string]api.Booking),
	}
	s.seed()
	return s
}

func (s *MemoryStore) seed() {
	s.owner = api.Owner{
		Id:          uuid.NewString(),
		Name:        "Владелец календаря",
		Email:       "owner@example.com",
		Description: "Записывайтесь на звонок в удобное время",
		TimeZone:    "Europe/Moscow",
	}

	def := api.EventType{
		Id:          uuid.NewString(),
		Title:       "30-минутный звонок",
		Description: "Быстрый созвон на 30 минут",
		Duration:    30,
		IsDisabled:  false,
	}
	s.eventTypes[def.Id] = def
}

// ── Owner ────────────────────────────────────────────────

func (s *MemoryStore) GetOwner() api.Owner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.owner
}

func (s *MemoryStore) UpdateOwner(o api.Owner) api.Owner {
	s.mu.Lock()
	defer s.mu.Unlock()
	o.Id = s.owner.Id
	s.owner = o
	return s.owner
}

// ── EventTypes ───────────────────────────────────────────

// ListEventTypes возвращает все типы (includeDisabled=true) либо только активные.
func (s *MemoryStore) ListEventTypes(includeDisabled bool) []api.EventType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]api.EventType, 0, len(s.eventTypes))
	for _, et := range s.eventTypes {
		if includeDisabled || !et.IsDisabled {
			result = append(result, et)
		}
	}
	return result
}

func (s *MemoryStore) GetEventType(id string) (api.EventType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	et, ok := s.eventTypes[id]
	if !ok {
		return api.EventType{}, domain.ErrNotFound
	}
	return et, nil
}

func (s *MemoryStore) CreateEventType(req api.CreateEventTypeRequest) api.EventType {
	s.mu.Lock()
	defer s.mu.Unlock()
	et := api.EventType{
		Id:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		IsDisabled:  false,
	}
	s.eventTypes[et.Id] = et
	return et
}

func (s *MemoryStore) UpdateEventType(id string, req api.CreateEventTypeRequest) (api.EventType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	et, ok := s.eventTypes[id]
	if !ok {
		return api.EventType{}, domain.ErrNotFound
	}
	et.Title = req.Title
	et.Description = req.Description
	et.Duration = req.Duration
	s.eventTypes[id] = et
	return et, nil
}

// DisableEventType помечает тип события как отключённый (soft-delete).
func (s *MemoryStore) DisableEventType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	et, ok := s.eventTypes[id]
	if !ok {
		return domain.ErrNotFound
	}
	et.IsDisabled = true
	s.eventTypes[id] = et
	return nil
}

// ── Bookings ─────────────────────────────────────────────

func (s *MemoryStore) ListBookings() []api.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]api.Booking, 0, len(s.bookings))
	for _, b := range s.bookings {
		result = append(result, b)
	}
	return result
}

// CreateBooking атомарно проверяет занятость и создаёт бронь.
// Правило: на одно время нельзя две брони (по всем типам событий).
func (s *MemoryStore) CreateBooking(eventTypeID, guestName, guestEmail string, start, end time.Time) (api.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, b := range s.bookings {
		if start.Before(b.EndTime) && b.StartTime.Before(end) {
			return api.Booking{}, domain.ErrSlotTaken
		}
	}

	b := api.Booking{
		Id:          uuid.NewString(),
		EventTypeId: eventTypeID,
		GuestName:   guestName,
		GuestEmail:  guestEmail,
		StartTime:   start,
		EndTime:     end,
		CreatedAt:   time.Now().UTC(),
	}
	s.bookings[b.Id] = b
	return b, nil
}
