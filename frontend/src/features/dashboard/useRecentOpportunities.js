import { getRecentOpportunities } from '@/services/apiDashboardFunctions';
import { useUser } from '../authentication/useUser';
import { useQuery } from '@tanstack/react-query';

export function useRecentOpportunities() {
  const { token, batch } = useUser();
  const { data: recentOpportunities, isLoading: isRecentOpportunitiesLoading } =
    useQuery({
      queryKey: ['recent-opportunities', token, batch],
      queryFn: () => getRecentOpportunities(token, batch),
    });

  return { recentOpportunities, isRecentOpportunitiesLoading };
}
