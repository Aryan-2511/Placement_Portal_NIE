import { Button } from '@/components/ui/button';
import { useLogout } from './useLogout';

function Logout() {
  const { logout, isLoading } = useLogout();
  if (isLoading) return <p>Loading...</p>;
  return (
    <div>
      <Button onClick={logout}>Logout</Button>
    </div>
  );
}

export default Logout;
