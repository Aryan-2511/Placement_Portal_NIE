import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { addEvent } from '@/services/apiSchedule';

export default function useAddEvent() {
  return useMutation({
    mutationFn: (event) => addEvent(event),
    onSuccess: () => {
      toast.success('Event added successfully!');
    },
    onError: () => {
      toast.error('Failed to add event!');
    },
  });
}