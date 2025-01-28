import { useLogout } from './useLogout';
import { HiMiniArrowRightStartOnRectangle } from 'react-icons/hi2';
import Spinner from '@/components/shared/Spinner';

function Logout() {
  const { logout, isLoading } = useLogout();
  return (
    <button onClick={logout}>
      {isLoading ? (
        <Spinner />
      ) : (
        <HiMiniArrowRightStartOnRectangle size={'2.4rem'} />
      )}
    </button>
  );
}

export default Logout;
