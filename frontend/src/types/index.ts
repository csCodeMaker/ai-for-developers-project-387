export interface Owner {
  id: string;
  name: string;
  email: string;
  description: string;
  timeZone: string;
}

export interface EventType {
  id: string;
  title: string;
  description: string;
  duration: number;
  isDisabled: boolean;
}

export interface Booking {
  id: string;
  eventTypeId: string;
  guestName: string;
  guestEmail: string;
  startTime: string;
  endTime: string;
  createdAt: string;
}

export interface AvailableSlot {
  startTime: string;
  endTime: string;
}

export interface Slot {
  startTime: string;
  endTime: string;
  isBusy: boolean;
}

export interface ErrorResponse {
  code: string;
  message: string;
}

export interface CreateBookingRequest {
  eventTypeId: string;
  guestName: string;
  guestEmail: string;
  startTime: string;
}

export interface CreateEventTypeRequest {
  title: string;
  description: string;
  duration: number;
}
