import { useState, useEffect } from 'react';
import { Paper, Stack, TextInput, Textarea, Button, Text } from '@mantine/core';
import type { Owner } from '../../types';

interface OwnerFormProps {
  owner: Owner | null;
  loading: boolean;
  onSave: (data: Owner) => Promise<void>;
}

export function OwnerForm({ owner, loading, onSave }: OwnerFormProps) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [description, setDescription] = useState('');
  const [timeZone, setTimeZone] = useState('');

  useEffect(() => {
    if (owner) {
      setName(owner.name);
      setEmail(owner.email);
      setDescription(owner.description);
      setTimeZone(owner.timeZone);
    }
  }, [owner]);
  const [busy, setBusy] = useState(false);

  if (loading) return <Text c="dimmed" size="sm">Загрузка профиля...</Text>;
  if (!owner) return <Text c="dimmed" size="sm">Профиль не найден.</Text>;

  const handleSubmit = async () => {
    if (!name.trim() || !email.trim()) return;
    setBusy(true);
    try {
      await onSave({ ...owner, name: name.trim(), email: email.trim(), description: description.trim(), timeZone: timeZone.trim() });
    } finally {
      setBusy(false);
    }
  };

  return (
    <Paper p="md" withBorder data-testid="owner-form">
      <Stack gap="sm">
        <TextInput label="Имя" value={name} onChange={(e) => setName(e.target.value)} required data-testid="owner-name" />
        <TextInput label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required data-testid="owner-email" />
        <Textarea label="Описание" value={description} onChange={(e) => setDescription(e.target.value)} rows={2} data-testid="owner-description" />
        <TextInput label="Часовой пояс" value={timeZone} onChange={(e) => setTimeZone(e.target.value)} placeholder="Europe/Moscow" data-testid="owner-timezone" />
        <Button onClick={handleSubmit} loading={busy} data-testid="save-owner">Сохранить</Button>
      </Stack>
    </Paper>
  );
}
