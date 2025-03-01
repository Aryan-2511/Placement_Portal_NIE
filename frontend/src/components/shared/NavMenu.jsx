import { useDarkMode } from '@/context/DarkModeContext';
import Logout from '@/features/authentication/Logout';
import { HiMiniUser, HiSun, HiMoon } from 'react-icons/hi2';
import { NavLink } from 'react-router-dom';

function NavMenu() {
  const { isDarkMode, toggleDarkMode } = useDarkMode();

  return (
    <ul className="flex items-center gap-12 text-[var(--color-blue-700)]">
      <li>
        <NavLink to="profile">
          <HiMiniUser size={'2.4rem'} />
        </NavLink>
      </li>
      <li>
        <button onClick={toggleDarkMode}>
          {isDarkMode ? <HiSun size={'2.4rem'} /> : <HiMoon size={'2.4rem'} />}
        </button>
      </li>
      <li>
        <Logout />
      </li>
    </ul>
  );
}

export default NavMenu;
