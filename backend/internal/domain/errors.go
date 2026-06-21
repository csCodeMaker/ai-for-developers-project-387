package domain

import "errors"

var (
	// ErrNotFound — сущность не найдена.
	ErrNotFound = errors.New("not found")
	// ErrSlotTaken — слот уже занят другой бронью.
	ErrSlotTaken = errors.New("slot taken")
	// ErrEventTypeDisabled — тип события отключён.
	ErrEventTypeDisabled = errors.New("event type disabled")
	// ErrValidation — некорректные данные запроса.
	ErrValidation = errors.New("validation failed")
)
