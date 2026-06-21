import { Outlet, useLocation } from 'react-router-dom';
import { AppShell, Group, Anchor } from '@mantine/core';

export function Layout() {
  const { pathname } = useLocation();
  const isAdmin = pathname.startsWith('/admin');

  return (
    <AppShell header={{ height: 56 }} padding="md">
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Anchor href="/" underline="hover" c="dark" fw={700} fz="lg">
            Календарь звонков
          </Anchor>
          <Anchor
            href={isAdmin ? '/' : '/admin'}
            underline="hover"
            c="dimmed"
            size="sm"
          >
            {isAdmin ? 'Гостевой режим' : 'Админка'}
          </Anchor>
        </Group>
      </AppShell.Header>
      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
