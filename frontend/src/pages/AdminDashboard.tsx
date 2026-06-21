import { useState, useEffect } from 'react';
import { Container, Title, Stack, Alert } from '@mantine/core';
import { useBookings } from '../hooks/useBookings';
import {
  fetchAdminEventTypes,
  createEventType,
  updateEventType,
  deleteEventType,
  fetchOwner,
  updateOwner,
} from '../api/admin';
import { BookingsList } from '../components/BookingsList/BookingsList';
import { EventTypeManager } from '../components/EventTypeManager/EventTypeManager';
import { OwnerForm } from '../components/OwnerForm/OwnerForm';
import type { EventType, Owner, CreateEventTypeRequest } from '../types';

export function AdminDashboard() {
  const { bookings, loading: bookingsLoading, error: bookingsError } = useBookings();
  const [eventTypes, setEventTypes] = useState<EventType[]>([]);
  const [eventTypesLoading, setEventTypesLoading] = useState(true);
  const [eventTypesError, setEventTypesError] = useState<string | null>(null);
  const [owner, setOwner] = useState<Owner | null>(null);
  const [ownerLoading, setOwnerLoading] = useState(true);
  const [ownerError, setOwnerError] = useState<string | null>(null);

  const loadEventTypes = () => {
    setEventTypesLoading(true);
    setEventTypesError(null);
    fetchAdminEventTypes()
      .then(setEventTypes)
      .catch((e) => setEventTypesError(e.message))
      .finally(() => setEventTypesLoading(false));
  };

  const loadOwner = () => {
    setOwnerLoading(true);
    setOwnerError(null);
    fetchOwner()
      .then(setOwner)
      .catch((e) => setOwnerError(e.message))
      .finally(() => setOwnerLoading(false));
  };

  useEffect(() => {
    loadEventTypes();
    loadOwner();
  }, []);

  const handleCreate = async (data: CreateEventTypeRequest) => {
    await createEventType(data);
    loadEventTypes();
  };

  const handleUpdate = async (id: string, data: CreateEventTypeRequest) => {
    await updateEventType(id, data);
    loadEventTypes();
  };

  const handleDelete = async (id: string) => {
    await deleteEventType(id);
    loadEventTypes();
  };

  const handleOwnerSave = async (data: Owner) => {
    const updated = await updateOwner(data);
    setOwner(updated);
  };

  return (
    <Container size="lg" py="xl">
      <Title order={1} mb="lg">Панель администратора</Title>

      <Stack gap="xl">
        <div>
          <Title order={3} mb="sm">Предстоящие брони</Title>
          {bookingsError && <Alert color="red">{bookingsError}</Alert>}
          <BookingsList bookings={bookings} loading={bookingsLoading} />
        </div>

        <div>
          <Title order={3} mb="sm">Типы событий</Title>
          {eventTypesError && <Alert color="red">{eventTypesError}</Alert>}
          <EventTypeManager
            eventTypes={eventTypes}
            loading={eventTypesLoading}
            onCreate={handleCreate}
            onUpdate={handleUpdate}
            onDelete={handleDelete}
          />
        </div>

        <div>
          <Title order={3} mb="sm">Профиль владельца</Title>
          {ownerError && <Alert color="red">{ownerError}</Alert>}
          <OwnerForm owner={owner} loading={ownerLoading} onSave={handleOwnerSave} />
        </div>
      </Stack>
    </Container>
  );
}
