import { useState } from 'react';
import {
  Table,
  Button,
  Group,
  Stack,
  TextInput,
  Textarea,
  Modal,
  Text,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import type { EventType, CreateEventTypeRequest } from '../../types';

interface EventTypeManagerProps {
  eventTypes: EventType[];
  loading: boolean;
  onCreate: (data: CreateEventTypeRequest) => Promise<void>;
  onUpdate: (id: string, data: CreateEventTypeRequest) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

export function EventTypeManager({
  eventTypes,
  loading,
  onCreate,
  onUpdate,
  onDelete,
}: EventTypeManagerProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [duration, setDuration] = useState(30);
  const [busy, setBusy] = useState(false);

  const resetForm = () => {
    setEditingId(null);
    setTitle('');
    setDescription('');
    setDuration(30);
    close();
  };

  const openCreate = () => {
    resetForm();
    open();
  };

  const openEdit = (et: EventType) => {
    setEditingId(et.id);
    setTitle(et.title);
    setDescription(et.description);
    setDuration(et.duration);
    open();
  };

  const handleSubmit = async () => {
    if (!title.trim() || !description.trim()) return;
    setBusy(true);
    try {
      if (editingId) {
        await onUpdate(editingId, { title: title.trim(), description: description.trim(), duration });
      } else {
        await onCreate({ title: title.trim(), description: description.trim(), duration });
      }
      resetForm();
    } finally {
      setBusy(false);
    }
  };

  const handleDelete = async (id: string) => {
    setBusy(true);
    try {
      await onDelete(id);
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      <Group justify="space-between" mb="sm">
        <Text size="sm" c="dimmed">
          {loading ? 'Загрузка...' : `${eventTypes.length} типов событий`}
        </Text>
        <Button onClick={openCreate} size="xs" data-testid="create-event-type">Создать</Button>
      </Group>

      <Modal
        opened={opened}
        onClose={close}
        title={editingId ? 'Редактировать тип события' : 'Создать тип события'}
        size="sm"
      >
        <Stack gap="sm">
          <TextInput label="Название" value={title} onChange={(e) => setTitle(e.target.value)} required data-testid="event-type-title" />
          <Textarea label="Описание" value={description} onChange={(e) => setDescription(e.target.value)} required rows={2} data-testid="event-type-description" />
          <TextInput label="Длительность (мин)" type="number" min={5} max={120} value={duration} onChange={(e) => setDuration(Number(e.target.value))} required data-testid="event-type-duration" />
          <Group justify="flex-end" mt="md">
            <Button variant="default" onClick={close}>Отмена</Button>
            <Button onClick={handleSubmit} loading={busy} data-testid="save-event-type">{editingId ? 'Обновить' : 'Создать'}</Button>
          </Group>
        </Stack>
      </Modal>

        <Table striped highlightOnHover withTableBorder data-testid="event-type-list">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Название</Table.Th>
            <Table.Th>Длительность</Table.Th>
            <Table.Th>Статус</Table.Th>
            <Table.Th></Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {eventTypes.map((et) => (
            <Table.Tr key={et.id} data-testid="event-type-row">
              <Table.Td fw={500}>{et.title}</Table.Td>
              <Table.Td>{et.duration} мин</Table.Td>
              <Table.Td>{et.isDisabled ? 'Отключён' : 'Активен'}</Table.Td>
              <Table.Td>
                <Group gap="xs">
                  <Button variant="outline" size="xs" onClick={() => openEdit(et)} disabled={busy} data-testid="edit-event-type">
                    Ред.
                  </Button>
                  <Button variant="outline" color="red" size="xs" onClick={() => handleDelete(et.id)} disabled={busy} data-testid="delete-event-type">
                    {et.isDisabled ? 'Удалить' : 'Откл.'}
                  </Button>
                </Group>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </>
  );
}
