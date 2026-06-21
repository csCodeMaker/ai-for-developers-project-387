package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := Load()
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, 8, cfg.SlotWorkStartHour)
	assert.Equal(t, 20, cfg.SlotWorkEndHour)
	assert.Equal(t, 30, cfg.SlotIntervalMin)
	assert.Equal(t, 14, cfg.BookingWindowDays)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("SLOT_WORK_START", "9")
	os.Setenv("SLOT_WORK_END", "18")
	os.Setenv("SLOT_INTERVAL_MIN", "15")
	os.Setenv("BOOKING_WINDOW_DAYS", "7")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("SLOT_WORK_START")
		os.Unsetenv("SLOT_WORK_END")
		os.Unsetenv("SLOT_INTERVAL_MIN")
		os.Unsetenv("BOOKING_WINDOW_DAYS")
	}()

	cfg := Load()
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, 9, cfg.SlotWorkStartHour)
	assert.Equal(t, 18, cfg.SlotWorkEndHour)
	assert.Equal(t, 15, cfg.SlotIntervalMin)
	assert.Equal(t, 7, cfg.BookingWindowDays)
}

func TestLoad_InvalidEnvFallsBack(t *testing.T) {
	os.Setenv("SLOT_WORK_START", "not-a-number")
	defer os.Unsetenv("SLOT_WORK_START")

	cfg := Load()
	assert.Equal(t, 8, cfg.SlotWorkStartHour)
}

func TestLoad_EmptyEnvFallsBack(t *testing.T) {
	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	cfg := Load()
	assert.Equal(t, "3000", cfg.Port)
}
