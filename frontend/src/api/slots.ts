import { apiRequest } from './client';
import type { Slot } from '../types';

export function fetchSlots(
  eventTypeId: string,
  date: string,
): Promise<Slot[]> {
  return apiRequest(
    `/event-types/${eventTypeId}/slots?date=${encodeURIComponent(date)}`,
  );
}
