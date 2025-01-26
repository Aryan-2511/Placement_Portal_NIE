import { getTotalPlacedStudents } from '@/services/apiDashboardFunctions';
import { useUser } from '../authentication/useUser';
import { useQuery } from '@tanstack/react-query';

export function useTotalPlacedStudents() {
  const { token, batch } = useUser();
  const { data: totalPlacedStudents, isLoading: isTotalPlacedStudentsLoading } =
    useQuery({
      queryKey: ['total-placed-students', token, batch],
      queryFn: () => getTotalPlacedStudents(token, batch),
    });

  return { totalPlacedStudents, isTotalPlacedStudentsLoading };
}
