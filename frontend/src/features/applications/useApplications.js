import { getApplications } from '@/services/apiApplications';
import { useQuery } from '@tanstack/react-query';
import { useUser } from '../authentication/useUser';

export function useApplications() {
  const { usn, token } = useUser();
  const {
    data: applications,
    error,
    isLoading,
  } = useQuery({
    queryKey: ['applications', usn, token],
    queryFn: () => getApplications(usn, token),
    retry: 1,
    enabled: !!usn,
  });
  return { applications, error, isLoading };
}
