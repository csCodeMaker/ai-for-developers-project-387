import { useNavigate, useParams } from 'react-router-dom';
import { Container, Paper, ThemeIcon, Title, Text, Button } from '@mantine/core';

export function ConfirmationPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  return (
    <Container size="xs" py="xl">
      <Paper
        p="xl"
        withBorder
        ta="center"
        style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 12 }}
      >
        <ThemeIcon size="xl" radius="xl" color="green">
          ✓
        </ThemeIcon>
        <Title order={2}>Бронь подтверждена!</Title>
        <Text c="dimmed" size="sm" data-testid="booking-id">
          Ваша бронь <strong>#{id}</strong> создана.
        </Text>
        <Text c="dimmed" size="sm">
          Вы получите письмо с подтверждением.
        </Text>
        <Button onClick={() => navigate('/')} mt="md">
          На главную
        </Button>
      </Paper>
    </Container>
  );
}
