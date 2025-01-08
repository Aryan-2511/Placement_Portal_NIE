import { updateOpportunity } from '@/services/apiOpportunities';
import { useMutation } from 'react-query';
import { useParams } from 'react-router-dom';

export default function useUpdateOpportunity() {
  const { opportunityId } = useParams();
  const mutation = useMutation({
    mutationKey: ['opportunity', opportunityId],
    mutationFn: (updatedOpportunity) =>
      updateOpportunity(opportunityId, updatedOpportunity),
    retry: 1,
    enabled: !!opportunityId,
  });

  return mutation;
}
