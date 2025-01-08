import UserAvatar from '../../features/authentication/UserAvatar';
import NavMenu from '../shared/NavMenu';

function Navbar() {
  return (
    <header className="h-[6.8rem] shadow-lg bg-[var(--color-grey-0)] flex gap-12 items-center justify-end pr-[7.2rem] border-b-[0.1rem]">
      <UserAvatar />
      <NavMenu />
    </header>
  );
}

export default Navbar;
