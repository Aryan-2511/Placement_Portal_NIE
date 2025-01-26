import { getOpportunity } from '@/services/apiOpportunities';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

export default function useOpportunity() {
  const { opportunityId } = useParams();
  const {
    isLoading: isOpportunityLoading,
    data: opportunity,
    error,
  } = useQuery({
    queryKey: ['opportunity', opportunityId],
    queryFn: () => getOpportunity(opportunityId),
    retry: 1,
    enabled: !!opportunityId,
  });

  return { isOpportunityLoading, opportunity, error };
}
