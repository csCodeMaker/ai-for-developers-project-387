package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"booking-backend/internal/api"
	"booking-backend/internal/config"
	"booking-backend/internal/domain"
	"booking-backend/internal/slots"
)

// Store описывает зависимость хендлеров от хранилища.
type Store interface {
	GetOwner() api.Owner
	UpdateOwner(api.Owner) api.Owner

	ListEventTypes(includeDisabled bool) []api.EventType
	GetEventType(id string) (api.EventType, error)
	CreateEventType(api.CreateEventTypeRequest) api.EventType
	UpdateEventType(id string, req api.CreateEventTypeRequest) (api.EventType, error)
	DisableEventType(id string) error

	ListBookings() []api.Booking
	CreateBooking(eventTypeID, guestName, guestEmail string, start, end time.Time) (api.Booking, error)
}

// Handler реализует HTTP-эндпоинты сервиса.
type Handler struct {
	store Store
	cfg   config.Config
}

func NewHandler(store Store, cfg config.Config) *Handler {
	return &Handler{store: store, cfg: cfg}
}

// Routes возвращает chi-роутер со всеми эндпоинтами.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		// Guest
		r.Get("/event-types", h.listEventTypes)
		r.Get("/event-types/{eventTypeId}/slots", h.listSlots)
		r.Post("/bookings", h.createBooking)

		// Admin
		r.Get("/admin/owner", h.getOwner)
		r.Put("/admin/owner", h.updateOwner)
		r.Get("/admin/event-types", h.adminListEventTypes)
		r.Post("/admin/event-types", h.adminCreateEventType)
		r.Put("/admin/event-types/{id}", h.adminUpdateEventType)
		r.Delete("/admin/event-types/{id}", h.adminDeleteEventType)
		r.Get("/admin/bookings", h.adminListBookings)
	})

	return r
}

// ── Guest ────────────────────────────────────────────────

func (h *Handler) listEventTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.ListEventTypes(false))
}

func (h *Handler) listSlots(w http.ResponseWriter, r *http.Request) {
	eventTypeID := chi.URLParam(r, "eventTypeId")

	et, err := h.store.GetEventType(eventTypeID)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Тип события не найден")
		return
	}
	if et.IsDisabled {
		writeError(w, http.StatusNotFound, "DISABLED", "Тип события отключён")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Не указана дата")
		return
	}
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Некорректный формат даты")
		return
	}

	result := slots.Generate(date, time.Now(), h.store.ListBookings(), slots.Params{
		WorkStartHour: h.cfg.SlotWorkStartHour,
		WorkEndHour:   h.cfg.SlotWorkEndHour,
		IntervalMin:   h.cfg.SlotIntervalMin,
	})

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	var req api.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Некорректное тело запроса")
		return
	}

	if req.GuestName == "" || req.GuestEmail == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Имя и email обязательны")
		return
	}

	et, err := h.store.GetEventType(req.EventTypeId)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Тип события не найден")
		return
	}
	if et.IsDisabled {
		writeError(w, http.StatusConflict, "DISABLED", "Тип события отключён")
		return
	}

	start := req.StartTime
	now := time.Now()

	// Валидация окна записи
	if !start.After(now) {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Нельзя записаться на прошедшее время")
		return
	}
	maxDate := now.AddDate(0, 0, h.cfg.BookingWindowDays)
	if start.After(maxDate) {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Время вне окна записи")
		return
	}

	end := start.Add(time.Duration(et.Duration) * time.Minute)

	booking, err := h.store.CreateBooking(req.EventTypeId, req.GuestName, req.GuestEmail, start, end)
	if err != nil {
		if errors.Is(err, domain.ErrSlotTaken) {
			writeError(w, http.StatusConflict, "SLOT_TAKEN", "Слот уже занят")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", "Внутренняя ошибка")
		return
	}

	writeJSON(w, http.StatusOK, booking)
}

// ── Admin ────────────────────────────────────────────────

func (h *Handler) getOwner(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.GetOwner())
}

func (h *Handler) updateOwner(w http.ResponseWriter, r *http.Request) {
	var o api.Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Некорректное тело запроса")
		return
	}
	writeJSON(w, http.StatusOK, h.store.UpdateOwner(o))
}

func (h *Handler) adminListEventTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.ListEventTypes(true))
}

func (h *Handler) adminCreateEventType(w http.ResponseWriter, r *http.Request) {
	var req api.CreateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Некорректное тело запроса")
		return
	}
	if req.Title == "" || req.Duration <= 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Название и длительность обязательны")
		return
	}
	writeJSON(w, http.StatusOK, h.store.CreateEventType(req))
}

func (h *Handler) adminUpdateEventType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req api.CreateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", "Некорректное тело запроса")
		return
	}
	et, err := h.store.UpdateEventType(id, req)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Тип события не найден")
		return
	}
	writeJSON(w, http.StatusOK, et)
}

func (h *Handler) adminDeleteEventType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.store.DisableEventType(id); err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Тип события не найден")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) adminListBookings(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.ListBookings())
}

// ── helpers ──────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, api.ErrorResponse{Code: code, Message: message})
}
