import { updateCurrentOpportunity } from '@/services/apiOpportunities';
import toast from 'react-hot-toast';
import { useMutation, useQueryClient } from '@tanstack/react-query';

export default function useUpdateOpportunity() {
  const queryClient = useQueryClient();

  const { mutate: updateOpportunity, isPending: isUpdating } = useMutation(
    {
      mutationFn: ({ opportunityId, opportunityData, role }) =>
        updateCurrentOpportunity(opportunityId, opportunityData, role),
      onSuccess: (opportunity) => {
        toast.success('Opporunity updated successfully!');
        queryClient.setQueryData(['opportunity', opportunity]);
      },
      onError: (err) => toast.error(err.message),
    },
    queryClient
  );

  return { updateOpportunity, isUpdating };
}
