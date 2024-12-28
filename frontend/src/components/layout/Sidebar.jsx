import Logo from '../shared/Logo';
import MainNav from '../shared/MainNav';

function Sidebar() {
  return (
    <aside className="row-span-full shadow-[var(--shadow-lg)] bg-[var(--color-grey-0)] h-full">
      <Logo />
      <MainNav />
    </aside>
  );
}

export default Sidebar;
