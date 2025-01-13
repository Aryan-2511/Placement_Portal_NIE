import { useUser } from './useUser';

function UserAvatar() {
  const user = useUser();
  return (
    <div className="font-['Poppins']">
      <div className="flex justify-center gap-3">
        <div className="w-[4rem] h-[4rem] bg-slate-200 rounded-full"></div>
        <div>
          <p className="text-[1.6rem] text-[var(--color-grey-600) font-medium">
            {user.name}
          </p>
          <p className="text-[1.2rem] text-[var(--color-grey-100)]">
            BATCH - {user.batch}
          </p>
        </div>
      </div>
    </div>
  );
}

export default UserAvatar;
