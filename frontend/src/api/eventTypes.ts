import { apiRequest } from './client';
import type { EventType } from '../types';

export function fetchEventTypes(): Promise<EventType[]> {
  return apiRequest('/event-types');
}
