import { useState, useEffect } from 'react';
import { fetchEventTypes } from '../api/eventTypes';
import type { EventType } from '../types';

export function useEventTypes() {
  const [eventTypes, setEventTypes] = useState<EventType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchEventTypes()
      .then(setEventTypes)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  return { eventTypes, loading, error };
}
