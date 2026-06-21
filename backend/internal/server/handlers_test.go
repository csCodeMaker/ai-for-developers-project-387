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
