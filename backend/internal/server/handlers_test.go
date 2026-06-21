package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"booking-backend/internal/api"
	"booking-backend/internal/config"
	"booking-backend/internal/store"
)

func newTestServer() (http.Handler, *store.MemoryStore) {
	st := store.New()
	cfg := config.Load()
	h := NewHandler(st, cfg)
	return h.Routes(), st
}

func TestCreateBooking_ComputesEndTime(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id // duration 30

	start := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Minute)
	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: etID,
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
		StartTime:   start,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var b api.Booking
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &b))
	assert.Equal(t, start.Add(30*time.Minute), b.EndTime.UTC())
}

func TestCreateBooking_Conflict409(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	start := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Minute)
	mk := func() *httptest.ResponseRecorder {
		body, _ := json.Marshal(api.CreateBookingRequest{
			EventTypeId: etID, GuestName: "Г", GuestEmail: "g@example.com", StartTime: start,
		})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
		srv.ServeHTTP(rec, req)
		return rec
	}

	require.Equal(t, http.StatusOK, mk().Code)

	rec := mk()
	require.Equal(t, http.StatusConflict, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "SLOT_TAKEN", e.Code)
}

// ── Admin: Owner ─────────────────────────────────────────

func TestGetOwner(t *testing.T) {
	srv, st := newTestServer()
	expected := st.GetOwner()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/owner", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var o api.Owner
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &o))
	assert.Equal(t, expected.Id, o.Id)
	assert.Equal(t, expected.Name, o.Name)
}

func TestUpdateOwner(t *testing.T) {
	srv, st := newTestServer()

	body, _ := json.Marshal(api.Owner{
		Name: "Обновлённый", Email: "updated@example.com",
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/owner", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var o api.Owner
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &o))
	assert.Equal(t, "Обновлённый", o.Name)
	assert.Equal(t, st.GetOwner().Id, o.Id)
}

func TestUpdateOwner_BadRequest(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/owner", bytes.NewReader([]byte("{")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// ── Admin: Event Types ───────────────────────────────────

func TestAdminListEventTypes(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/event-types", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var list []api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
	assert.Len(t, list, 1)
}

func TestAdminCreateEventType_OK(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "Новый тип", Description: "Описание", Duration: 45,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var et api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &et))
	assert.Equal(t, "Новый тип", et.Title)
	assert.Equal(t, int32(45), et.Duration)
	assert.False(t, et.IsDisabled)
}

func TestAdminCreateEventType_BadBody(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader([]byte("{")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminCreateEventType_Validation(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{Title: "", Duration: 0})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminUpdateEventType_OK(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "Обновлённый", Description: "Новое описание", Duration: 60,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/event-types/"+etID, bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var et api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &et))
	assert.Equal(t, "Обновлённый", et.Title)
	assert.Equal(t, int32(60), et.Duration)
}

func TestAdminUpdateEventType_BadBody(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/event-types/some-id", bytes.NewReader([]byte("{")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminUpdateEventType_NotFound(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{Title: "X", Duration: 30})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/event-types/nonexistent", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminDeleteEventType(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "Для удаления", Duration: 30})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/event-types/"+created.Id, nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)

	et, err := st.GetEventType(created.Id)
	require.NoError(t, err)
	assert.True(t, et.IsDisabled)
}

func TestAdminDeleteEventType_NotFound(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/event-types/nonexistent", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Admin: Bookings ──────────────────────────────────────

func TestAdminListBookings(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/bookings", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var list []api.Booking
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
	assert.Empty(t, list)
}

// ── Guest: listSlots ─────────────────────────────────────

func TestListSlots_OK(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots?date=2030-01-11", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var slots []api.Slot
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &slots))
	assert.NotEmpty(t, slots)
}

func TestListSlots_NotFound(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/nonexistent/slots?date=2030-01-11", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListSlots_Disabled(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "Откл", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+created.Id+"/slots?date=2030-01-11", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListSlots_MissingDate(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListSlots_InvalidDate(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots?date=invalid", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// ── Guest: createBooking validations ─────────────────────

func TestCreateBooking_BadBody(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader([]byte("{")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBooking_MissingFields(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateBookingRequest{GuestName: "", GuestEmail: ""})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBooking_EventTypeNotFound(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: "nonexistent",
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
		StartTime:   time.Now().Add(48 * time.Hour),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateBooking_DisabledEventType(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "Откл", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: created.Id,
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
		StartTime:   time.Now().Add(48 * time.Hour),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateBooking_PastTime(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: etID,
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
		StartTime:   time.Now().Add(-1 * time.Hour),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBooking_OutsideWindow(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: etID,
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
		StartTime:   time.Now().AddDate(0, 0, 30),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListEventTypes_GuestExcludesDisabled(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "X", Description: "d", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var list []api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
	for _, et := range list {
		assert.NotEqual(t, created.Id, et.Id)
	}
}
