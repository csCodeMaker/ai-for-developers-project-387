import { apiRequest } from './client';
import type { Booking, CreateBookingRequest } from '../types';

export function createBooking(
  data: CreateBookingRequest,
): Promise<Booking> {
  return apiRequest('/bookings', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function fetchAllBookings(): Promise<Booking[]> {
  return apiRequest('/admin/bookings');
}
