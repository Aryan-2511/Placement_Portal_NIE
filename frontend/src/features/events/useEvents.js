import { useQuery } from '@tanstack/react-query';
import { getAllEvents } from '@/services/apiSchedule';

export default function useEvents() {
  return useQuery({
    queryKey: ['events'],
    queryFn: getAllEvents,
  });
}
