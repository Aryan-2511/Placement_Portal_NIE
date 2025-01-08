import Logo from '../shared/Logo';
import MainNav from '../shared/MainNav';

function Sidebar() {
  return (
    <aside className="row-span-full shadow-lg bg-[var(--color-grey-0)] h-full  border-r-[0.1rem]">
      <Logo />
      <MainNav />
    </aside>
  );
}

export default Sidebar;
