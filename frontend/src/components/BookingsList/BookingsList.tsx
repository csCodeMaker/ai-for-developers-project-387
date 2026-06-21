import { Table, Text } from '@mantine/core';
import type { Booking } from '../../types';

interface BookingsListProps {
  bookings: Booking[];
  loading: boolean;
}

function formatDateTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleString('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  });
}

export function BookingsList({ bookings, loading }: BookingsListProps) {
  if (loading) return <Text c="dimmed" size="sm">Загрузка броней...</Text>;
  if (bookings.length === 0) return <Text c="dimmed" size="sm">Нет предстоящих броней.</Text>;

  return (
    <Table striped highlightOnHover withTableBorder data-testid="bookings-table">
      <Table.Thead>
        <Table.Tr>
          <Table.Th>Гость</Table.Th>
          <Table.Th>Email</Table.Th>
          <Table.Th>Начало</Table.Th>
          <Table.Th>Конец</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {bookings.map((b) => (
          <Table.Tr key={b.id} data-testid="booking-row">
            <Table.Td fw={500}>{b.guestName}</Table.Td>
            <Table.Td>{b.guestEmail}</Table.Td>
            <Table.Td>{formatDateTime(b.startTime)}</Table.Td>
            <Table.Td>{formatDateTime(b.endTime)}</Table.Td>
          </Table.Tr>
        ))}
      </Table.Tbody>
    </Table>
  );
}
