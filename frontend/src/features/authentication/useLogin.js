import { login as loginApi } from '@/services/apiAuth';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import Cookies from 'js-cookie';
import { toast } from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

export default function useLogin() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { mutate: login, isLoading } = useMutation({
    mutationFn: ({ email, password, role }) =>
      loginApi({ email, password, role }),
    onSuccess: (user) => {
      queryClient.setQueryData(['user'], user);
      Cookies.set('user', JSON.stringify(user), { expires: (20 / 24) * 60 });
      toast.success('Login successful');
      navigate('student/dashboard', { replace: true });
    },
    onError: (err) => {
      console.log('ERROR', err);
      toast.error('Provided email or password are incorrect!');
    },
  });

  return { login, isLoading };
}
