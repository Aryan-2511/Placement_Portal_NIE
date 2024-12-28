import { NavLink } from 'react-router-dom';

import {
  HiMiniHome,
  HiClipboardDocumentList,
  HiPencilSquare,
  HiMiniUserCircle,
  HiMiniChatBubbleLeft,
  HiMiniQuestionMarkCircle,
} from 'react-icons/hi2';

function MainNav() {
  return (
    <nav>
      <ul>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="dashboard"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiMiniHome size={'2.4rem'} />
            <span>Home</span>
          </NavLink>
        </li>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="opportunities"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiClipboardDocumentList size={'2.4rem'} />
            <span>Opportunities</span>
          </NavLink>
        </li>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="applications"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiPencilSquare size={'2.4rem'} />
            <span>Applications</span>
          </NavLink>
        </li>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="profile"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiMiniUserCircle size={'2.4rem'} />
            <span>Profile</span>
          </NavLink>
        </li>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="feedback"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiMiniChatBubbleLeft size={'2.4rem'} />
            <span>Feedback</span>
          </NavLink>
        </li>
        <li className="w-full h-[6.4rem]">
          <NavLink
            to="faq"
            className={({ isActive }) =>
              `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                isActive
                  ? 'bg-[var(--color-blue-700)] text-white'
                  : 'text-[var(--color-grey-600)] hover:bg-[var(--color-blue-0)]'
              }`
            }
          >
            <HiMiniQuestionMarkCircle size={'2.4rem'} />
            <span>FAQ</span>
          </NavLink>
        </li>
      </ul>
    </nav>
  );
}

export default MainNav;
