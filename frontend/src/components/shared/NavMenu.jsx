import { useState } from 'react';
import {
  HiMiniUser,
  HiSun,
  HiMoon,
  HiMiniArrowRightStartOnRectangle,
} from 'react-icons/hi2';
import { NavLink } from 'react-router-dom';

function NavMenu() {
  const [isDarkMode, setIsDarkMode] = useState(false);
  function handleToggleDarkMode() {
    setIsDarkMode((isDarkMode) => !isDarkMode);
  }
  return (
    <ul className="flex items-center gap-8 text-[var(--color-blue-700)]">
      <li>
        <NavLink to="student/profile">
          <HiMiniUser size={'2.4rem'} />
        </NavLink>
      </li>
      <li>
        <NavLink onClick={handleToggleDarkMode}>
          {isDarkMode ? <HiSun size={'2.4rem'} /> : <HiMoon size={'2.4rem'} />}
        </NavLink>
      </li>
      <li>
        <NavLink to="student/logout">
          <HiMiniArrowRightStartOnRectangle size={'2.4rem'} />
        </NavLink>
      </li>
    </ul>
  );
}

export default NavMenu;
