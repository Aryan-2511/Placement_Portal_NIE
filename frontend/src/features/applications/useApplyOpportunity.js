import { useUser } from '../authentication/useUser';
import { applyToOpportunity } from '@/services/apiApplications';
import { useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { useParams } from 'react-router-dom';

export function useApplyOpportunity() {
  const { opportunityId } = useParams();
  const { usn, token } = useUser();
  const {
    mutate: apply,
    isLoading: isApplicationLoading,
    error,
  } = useMutation({
    mutationKey: ['opportunity', 'apply', opportunityId],
    mutationFn: () => applyToOpportunity(usn, opportunityId, token),
    onSuccess: () => {
      toast.success('Application successful!');
    },
    onError: (error) => {
      toast.error(error.message || 'error');
    },
  });
  return { apply, isApplicationLoading, error };
}
