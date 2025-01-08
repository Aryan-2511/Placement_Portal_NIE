import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';
import Sidebar from './Sidebar';

function MainLayout() {
  return (
    <div className="h-screen grid grid-cols-[27.4rem_1fr] grid-rows-[auto_1fr] text-[var(--color-grey-600)] overflow-hidden outline-dashed">
      <Sidebar />
      <Navbar />
      <main className="bg-[var(--color-grey-50)] py-[5rem] px-[7.2rem] font-['Poppins'] overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}

export default MainLayout;
