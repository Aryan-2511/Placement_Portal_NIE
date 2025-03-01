import { addNewOpportunity } from '@/services/apiOpportunities';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';

export default function useAddOpportunity() {
  return useMutation({
    mutationFn: (opportunity) => addNewOpportunity(opportunity),
    onSuccess: () => {
      toast.success('Opportunity added successfully!');
    },
    onError: () => {
      toast.error('Failed to add opportunity!');
    },
  });
}
