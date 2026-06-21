import { useNavigate } from 'react-router-dom';
import { Card, Text, Group, Button } from '@mantine/core';
import type { EventType } from '../../types';

interface EventTypeCardProps {
  eventType: EventType;
}

export function EventTypeCard({ eventType }: EventTypeCardProps) {
  const navigate = useNavigate();

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder data-testid="event-type-card">
      <Text fw={600} fz="lg" mb="xs">
        {eventType.title}
      </Text>
      <Text size="sm" c="dimmed" mb="sm" style={{ flex: 1 }}>
        {eventType.description}
      </Text>
      <Group justify="space-between" align="center">
        <Text size="sm" c="blue" fw={500}>
          {eventType.duration} мин
        </Text>
        <Button onClick={() => navigate(`/book/${eventType.id}`)}>
          Записаться
        </Button>
      </Group>
    </Card>
  );
}
