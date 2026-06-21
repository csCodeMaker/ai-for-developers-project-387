import { apiRequest } from './client';
import type {
  Owner,
  EventType,
  CreateEventTypeRequest,
} from '../types';

export function fetchOwner(): Promise<Owner> {
  return apiRequest('/admin/owner');
}

export function updateOwner(data: Owner): Promise<Owner> {
  return apiRequest('/admin/owner', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function fetchAdminEventTypes(): Promise<EventType[]> {
  return apiRequest('/admin/event-types');
}

export function createEventType(
  data: CreateEventTypeRequest,
): Promise<EventType> {
  return apiRequest('/admin/event-types', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function updateEventType(
  id: string,
  data: CreateEventTypeRequest,
): Promise<EventType> {
  return apiRequest(`/admin/event-types/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function deleteEventType(id: string): Promise<void> {
  return apiRequest(`/admin/event-types/${id}`, {
    method: 'DELETE',
  });
}
