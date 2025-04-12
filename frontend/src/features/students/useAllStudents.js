import { getAllStudentsByBatch } from '@/services/apiStudents';
import { useQuery } from '@tanstack/react-query';

export default function useAllStudents(batch) {
  const {
    isLoading,
    data: allStudents,
    error,
  } = useQuery({
    queryKey: ['allStudents', batch],
    queryFn: () => getAllStudentsByBatch(batch),
    retry: 1,
    enabled: !!batch,
  });
  return { isLoading, error, allStudents };
}
