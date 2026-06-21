import { Center, Loader as MantineLoader } from '@mantine/core';

export function Loader() {
  return (
    <Center h={200}>
      <MantineLoader />
    </Center>
  );
}
