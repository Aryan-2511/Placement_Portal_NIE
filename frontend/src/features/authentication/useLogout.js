import { useMutation, useQueryClient } from '@tanstack/react-query';
import { logout as logoutApi } from '../../services/apiAuth';
import { useNavigate } from 'react-router-dom';
import { clearItem } from '@/utils/localStorageServices';
import Cookies from 'js-cookie';

export function useLogout() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { mutate: logout, isLoading } = useMutation({
    mutationFn: logoutApi,
    onSuccess: () => {
      clearItem('currentUser');
      Cookies.remove('user');
      queryClient.removeQueries();
      navigate('/', { replace: true });
    },
  });

  return { logout, isLoading };
}
