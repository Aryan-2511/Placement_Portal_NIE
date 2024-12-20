import Logo from '../shared/Logo';
import MainNav from '../shared/MainNav';

function Sidebar() {
  return (
    <aside className="row-span-full shadow-lg bg-[var(--color-grey-0)] flex-col items-center justify-center">
      <Logo />
      <MainNav />
    </aside>
  );
}

export default Sidebar;
