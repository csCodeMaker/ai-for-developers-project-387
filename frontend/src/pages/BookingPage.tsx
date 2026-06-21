import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Container,
  Title,
  Text,
  Button,
  Stack,
  Paper,
  TextInput,
  Grid,
  Alert,
  ScrollArea,
} from '@mantine/core';
import { DatePicker } from '@mantine/dates';
import { useSlots } from '../hooks/useSlots';
import { createBooking } from '../api/bookings';
import { ApiError } from '../api/client';
import { Loader } from '../components/ui/Loader';

function formatTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
}

export function BookingPage() {
  const { eventTypeId } = useParams<{ eventTypeId: string }>();
  const navigate = useNavigate();
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<string | null>(null);
  const [guestName, setGuestName] = useState('');
  const [guestEmail, setGuestEmail] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const dateStr = selectedDate;
  const { slots, loading, error: slotsError, reload } = useSlots(eventTypeId!, dateStr);

  const today = new Date().toISOString().slice(0, 10);
  const maxDateObj = new Date();
  maxDateObj.setDate(maxDateObj.getDate() + 13);
  const maxDate = maxDateObj.toISOString().slice(0, 10);

  const handleSubmit = async () => {
    if (!eventTypeId || !selectedSlot || !guestName.trim() || !guestEmail.trim()) return;

    setSubmitting(true);
    setError(null);

    try {
      const booking = await createBooking({
        eventTypeId,
        guestName: guestName.trim(),
        guestEmail: guestEmail.trim(),
        startTime: selectedSlot,
      });
      navigate(`/booking/${booking.id}`);
    } catch (e: unknown) {
      if (e instanceof ApiError && e.status === 409) {
        // Слот заняли между выбором и отправкой — обновляем слоты.
        setError('Этот слот только что заняли. Выберите другое время.');
        setSelectedSlot(null);
        reload();
      } else {
        setError(e instanceof Error ? e.message : 'Ошибка бронирования');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Container size="lg" py="xl">
      <Button variant="subtle" onClick={() => navigate('/')} mb="lg">
        &larr; Назад
      </Button>

      <Title order={2} mb="xl">
        Выберите дату и время
      </Title>

      <Grid gap="xl">
        <Grid.Col span={{ base: 12, md: 5 }}>
          <Paper p="md" withBorder data-testid="date-picker">
            <DatePicker
              value={selectedDate}
              onChange={setSelectedDate}
              minDate={today}
              maxDate={maxDate}
            />
          </Paper>
        </Grid.Col>

        <Grid.Col span={{ base: 12, md: 7 }}>
          {!selectedDate && (
            <Text c="dimmed" ta="center" py="xl">
              Выберите дату, чтобы увидеть доступные слоты
            </Text>
          )}

          {loading && <Loader />}

          {slotsError && <Alert color="red">{slotsError}</Alert>}

          {selectedDate && !loading && !slotsError && (
            <Stack gap="sm">
              {!selectedSlot && error && <Alert color="red" data-testid="booking-error">{error}</Alert>}
              <Text fw={600} size="sm" c="dimmed">
                Доступное время на{' '}
                {new Date(selectedDate + 'T00:00:00').toLocaleDateString('ru-RU', {
                  day: 'numeric',
                  month: 'long',
                })}
              </Text>

              <ScrollArea.Autosize mah={320}>
                <Stack gap="xs">
                  {slots.length === 0 ? (
                    <Text c="dimmed" size="sm" py="md">
                      Нет доступных слотов на этот день
                    </Text>
                  ) : (
                    slots.map((slot) => {
                      const isSelected = slot.startTime === selectedSlot;
                      return (
                          <Button
                           key={slot.startTime}
                           fullWidth
                           disabled={slot.isBusy}
                           variant={
                             isSelected ? 'filled' : slot.isBusy ? 'default' : 'outline'
                           }
                           color={slot.isBusy ? 'gray' : undefined}
                           onClick={() =>
                             !slot.isBusy && setSelectedSlot(slot.startTime)
                           }
                           size="sm"
                           data-testid="slot"
                           data-busy={String(slot.isBusy)}
                           data-start-time={slot.startTime}
                         >
                          {formatTime(slot.startTime)}
                          {slot.isBusy ? ' · занято' : ''}
                        </Button>
                      );
                    })
                  )}
                </Stack>
              </ScrollArea.Autosize>

              {selectedSlot && (
                <Paper p="md" withBorder mt="md">
                  <Stack gap="sm">
                    <Text fw={600} size="sm">
                      Ваши данные
                    </Text>
                     <TextInput
                       label="Имя"
                       value={guestName}
                       onChange={(e) => setGuestName(e.target.value)}
                       placeholder="Иван Иванов"
                       required
                       data-testid="guest-name"
                     />
                     <TextInput
                       label="Email"
                       type="email"
                       value={guestEmail}
                       onChange={(e) => setGuestEmail(e.target.value)}
                       placeholder="ivan@example.com"
                       required
                       data-testid="guest-email"
                     />
                     {error && <Alert color="red" data-testid="booking-error">{error}</Alert>}
                     <Button
                       onClick={handleSubmit}
                       loading={submitting}
                       disabled={!guestName.trim() || !guestEmail.trim()}
                       data-testid="book-submit"
                     >
                      {submitting ? 'Бронирование...' : 'Записаться'}
                    </Button>
                  </Stack>
                </Paper>
              )}
            </Stack>
          )}
        </Grid.Col>
      </Grid>
    </Container>
  );
}
