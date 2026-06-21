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

func TestListSlots_OK(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id
	date := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots?date="+date, nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var slots []api.Slot
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &slots))
	assert.NotEmpty(t, slots)
}

func TestListSlots_MissingDate(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestListSlots_InvalidDate(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+etID+"/slots?date=bad-date", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestListSlots_EventTypeNotFound(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/nonexistent/slots?date=2030-01-01", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "NOT_FOUND", e.Code)
}

func TestListSlots_DisabledEventType(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "D", Description: "d", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/event-types/"+created.Id+"/slots?date=2030-01-01", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "DISABLED", e.Code)
}

func TestCreateBooking_ComputesEndTime(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

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

func TestCreateBooking_InvalidBody(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader([]byte("{invalid")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestCreateBooking_MissingFields(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: "x", GuestName: "", GuestEmail: "",
		StartTime: time.Now().Add(48 * time.Hour),
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestCreateBooking_PastTime(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: etID, GuestName: "N", GuestEmail: "n@n.com",
		StartTime: time.Now().Add(-1 * time.Hour),
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestCreateBooking_EventTypeNotFound(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: "nonexistent", GuestName: "N", GuestEmail: "n@n.com",
		StartTime: time.Now().Add(48 * time.Hour),
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "NOT_FOUND", e.Code)
}

func TestCreateBooking_DisabledEventType(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "D", Description: "d", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	body, _ := json.Marshal(api.CreateBookingRequest{
		EventTypeId: created.Id, GuestName: "N", GuestEmail: "n@n.com",
		StartTime: time.Now().Add(48 * time.Hour),
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "DISABLED", e.Code)
}

func TestGetOwner(t *testing.T) {
	srv, st := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/owner", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var o api.Owner
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &o))
	assert.Equal(t, st.GetOwner().Id, o.Id)
}

func TestUpdateOwner(t *testing.T) {
	srv, st := newTestServer()

	body, _ := json.Marshal(api.Owner{
		Name: "Updated", Email: "u@u.com", Description: "d", TimeZone: "UTC",
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/owner", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var o api.Owner
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &o))
	assert.Equal(t, st.GetOwner().Id, o.Id)
	assert.Equal(t, "Updated", o.Name)
}

func TestAdminListEventTypes(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "X", Description: "d", Duration: 30})
	require.NoError(t, st.DisableEventType(created.Id))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/event-types", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var list []api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
	require.Len(t, list, 2)
}

func TestAdminCreateEventType(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "New", Description: "d", Duration: 45,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var et api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &et))
	assert.Equal(t, "New", et.Title)
	assert.Equal(t, int32(45), et.Duration)
	assert.False(t, et.IsDisabled)
}

func TestAdminCreateEventType_MissingTitle(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "", Description: "d", Duration: 30,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestAdminCreateEventType_InvalidDuration(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "T", Description: "d", Duration: 0,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestAdminCreateEventType_InvalidBody(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/event-types", bytes.NewReader([]byte("{bad")))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "VALIDATION", e.Code)
}

func TestAdminUpdateEventType(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "Updated", Description: "d", Duration: 60,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/event-types/"+etID, bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var et api.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &et))
	assert.Equal(t, "Updated", et.Title)
	assert.Equal(t, int32(60), et.Duration)
}

func TestAdminUpdateEventType_NotFound(t *testing.T) {
	srv, _ := newTestServer()

	body, _ := json.Marshal(api.CreateEventTypeRequest{
		Title: "X", Description: "d", Duration: 30,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/event-types/nonexistent", bytes.NewReader(body))
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "NOT_FOUND", e.Code)
}

func TestAdminDeleteEventType(t *testing.T) {
	srv, st := newTestServer()
	created := st.CreateEventType(api.CreateEventTypeRequest{Title: "D", Description: "d", Duration: 30})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/event-types/"+created.Id, nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminDeleteEventType_NotFound(t *testing.T) {
	srv, _ := newTestServer()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/event-types/nonexistent", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var e api.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &e))
	assert.Equal(t, "NOT_FOUND", e.Code)
}

func TestAdminListBookings(t *testing.T) {
	srv, st := newTestServer()
	etID := st.ListEventTypes(false)[0].Id
	start := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Minute)

	makeBooking := func(name string, offset time.Duration) {
		body, _ := json.Marshal(api.CreateBookingRequest{
			EventTypeId: etID, GuestName: name, GuestEmail: name + "@e.com",
			StartTime: start.Add(offset),
		})
		req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}
	makeBooking("A", 0)
	makeBooking("B", 1*time.Hour)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/bookings", nil)
	srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var bookings []api.Booking
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &bookings))
	assert.Len(t, bookings, 2)
}
