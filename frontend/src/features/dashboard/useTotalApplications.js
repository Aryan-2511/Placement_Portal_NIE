import { getTotalApplications } from '@/services/apiDashboardFunctions';
import { useUser } from '../authentication/useUser';
import { useQuery } from '@tanstack/react-query';

export function useTotalApplications() {
  const { token, usn } = useUser();
  const { data: totalApplications, isLoading: isTotalApplicationsLoading } =
    useQuery({
      queryKey: ['total-applications', token, usn],
      queryFn: () => getTotalApplications(token, usn),
    });

  return { totalApplications, isTotalApplicationsLoading };
}
