import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';
import Sidebar from './Sidebar';

function MainLayout() {
  return (
    <div className="grid grid-cols-[27.4rem_1fr] grid-rows-[auto_1fr] h-screen">
      <Sidebar />
      <Navbar />
      <main className="bg-slate-200 pt-[5rem] px-[7.2rem] font-['Poppins']">
        <Outlet />
      </main>
    </div>
  );
}

export default MainLayout;
