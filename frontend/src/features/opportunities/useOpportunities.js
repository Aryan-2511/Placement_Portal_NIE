import { getOpportunitiesByBatch } from '@/services/apiOpportunities';
import { useQuery } from '@tanstack/react-query';

export default function useOpportunities(batch) {
  const {
    isLoading,
    data: opportunities,
    error,
  } = useQuery({
    queryKey: ['opportunities', batch],
    queryFn: () => getOpportunitiesByBatch(batch),
    retry: 1,
    enabled: !!batch,
  });
  return { isLoading, error, opportunities };
}
