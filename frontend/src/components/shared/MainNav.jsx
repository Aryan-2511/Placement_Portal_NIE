import {
  HiMiniHome,
  HiClipboardDocumentList,
  HiPencilSquare,
  HiMiniUserCircle,
  HiMiniChatBubbleLeft,
  HiMiniQuestionMarkCircle,
  HiCalendarDays,
  HiAcademicCap,
  HiChatBubbleBottomCenterText,
  HiMiniUsers,
} from 'react-icons/hi2';

import { NavLink } from 'react-router-dom';
import { useUser } from '@/features/authentication/useUser';
import Spinner from './Spinner';

const studentFields = [
  { label: 'Home', to: 'dashboard', icon: <HiMiniHome size={'2.4rem'} /> },
  {
    label: 'Opportunities',
    to: '/student/opportunities',
    icon: <HiClipboardDocumentList size={'2.4rem'} />,
  },
  {
    label: 'Applications',
    to: '/student/applications',
    icon: <HiPencilSquare size={'2.4rem'} />,
  },
  {
    label: 'Profile',
    to: '/student/profile',
    icon: <HiMiniUserCircle size={'2.4rem'} />,
  },
  {
    label: 'Feedback',
    to: '/student/feedback',
    icon: <HiMiniChatBubbleLeft size={'2.4rem'} />,
  },
  {
    label: 'FAQ',
    to: '/student/faq',
    icon: <HiMiniQuestionMarkCircle size={'2.4rem'} />,
  },
];
const adminFields = [
  { label: 'Home', to: 'dashboard', icon: <HiMiniHome size={'2.4rem'} /> },
  {
    label: 'Opportunities',
    to: 'manage_opportunities',
    icon: <HiClipboardDocumentList size={'2.4rem'} />,
  },
  {
    label: 'Student details',
    to: 'student_details',
    icon: <HiAcademicCap size={'2.4rem'} />,
  },
  {
    label: 'Schedule',
    to: 'schedule',
    icon: <HiCalendarDays size={'2.4rem'} />,
  },
  {
    label: 'Announcements',
    to: 'annoucements',
    icon: <HiChatBubbleBottomCenterText size={'2.4rem'} />,
  },
  {
    label: 'Admins panel',
    to: 'admin_panel',
    icon: <HiMiniUsers size={'2.4rem'} />,
  },
];
const coordinatorFields = [
  { label: 'Home', to: 'dashboard', icon: <HiMiniHome size={'2.4rem'} /> },
  {
    label: 'Opportunities',
    to: 'opportunities',
    icon: <HiClipboardDocumentList size={'2.4rem'} />,
  },
  {
    label: 'Applications',
    to: 'applications',
    icon: <HiPencilSquare size={'2.4rem'} />,
  },
  {
    label: 'Profile',
    to: 'profile',
    icon: <HiMiniUserCircle size={'2.4rem'} />,
  },
  {
    label: 'Feedback',
    to: 'feedback',
    icon: <HiMiniChatBubbleLeft size={'2.4rem'} />,
  },
  {
    label: 'FAQ',
    to: 'faq',
    icon: <HiMiniQuestionMarkCircle size={'2.4rem'} />,
  },
];

function MainNav() {
  const user = useUser();
  const { role } = user;
  let tabs;
  if (role === 'ADMIN') {
    tabs = adminFields;
  } else if (role === 'STUDENT') {
    tabs = studentFields;
  } else if (role === 'COORDINATOR') {
    tabs = coordinatorFields;
  }
  if (!tabs) return <Spinner />;
  return (
    <nav>
      <ul>
        {tabs.map((tab) => {
          return (
            <li key={tab.label} className="w-full h-[6.4rem]">
              <NavLink
                to={tab.to}
                className={({ isActive }) =>
                  `flex items-center justify-start gap-3 h-full w-full pl-[5.4rem] font-bold transition-all duration-300 ease-in ${
                    isActive
                      ? 'bg-[var(--color-brand-700)] text-[var(--color-brand-50)]'
                      : 'text-[var(--color-grey-600)] hover:bg-[var(--color-grey-50)]'
                  }`
                }
              >
                {tab.icon}
                <span>{tab.label}</span>
              </NavLink>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}

export default MainNav;
