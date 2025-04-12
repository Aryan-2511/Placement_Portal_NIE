import { getStudentDetails } from '@/services/apiAuth';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';

export default function useStudent() {
  const { usn } = useParams();
  const {
    isLoading,
    data: student,
    error,
  } = useQuery({
    queryKey: ['allStudents', usn],
    queryFn: () => getStudentDetails(usn),
    retry: 1,
    enabled: !!usn,
  });
  return { isLoading, error, student };
}
