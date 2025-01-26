import { getActiveOpportunities } from '@/services/apiDashboardFunctions';
import { useUser } from '../authentication/useUser';
import { useQuery } from '@tanstack/react-query';

export function useActiveOpportunities() {
  const { token, batch } = useUser();
  const { data: activeOpportunities, isLoading: isActiveOpportunitiesLoading } =
    useQuery({
      queryKey: ['active-opportunities', token, batch],
      queryFn: () => getActiveOpportunities(token, batch),
    });

  return { activeOpportunities, isActiveOpportunitiesLoading };
}
