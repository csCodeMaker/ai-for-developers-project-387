package config

import (
	"os"
	"strconv"
)

// Config содержит настройки сервиса, читаемые из переменных окружения.
type Config struct {
	Port              string
	SlotWorkStartHour int
	SlotWorkEndHour   int
	SlotIntervalMin   int
	BookingWindowDays int
}

// Load читает конфигурацию из env с дефолтными значениями.
func Load() Config {
	return Config{
		Port:              getEnv("PORT", "3000"),
		SlotWorkStartHour: getEnvInt("SLOT_WORK_START", 8),
		SlotWorkEndHour:   getEnvInt("SLOT_WORK_END", 20),
		SlotIntervalMin:   getEnvInt("SLOT_INTERVAL_MIN", 30),
		BookingWindowDays: getEnvInt("BOOKING_WINDOW_DAYS", 14),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
