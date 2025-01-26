import { useUser } from './useUser';

function UserAvatar() {
  const user = useUser();
  return (
    <div className="font-['Poppins']">
      <div className="flex justify-center gap-3">
        <div className="w-[4rem] h-[4rem] bg-[var(--color-grey-50)] border-[0.1rem] border-[var(--color-grey-100)] rounded-full flex items-center justify-center text-[2.8rem] font-semibold text-[var(--color-grey-100)]">
          {user.name[0]}
        </div>
        <div>
          <p className="text-[1.6rem] text-[var(--color-grey-600) font-medium">
            {user.name}
          </p>
          <p className="text-[1.2rem] text-[var(--color-grey-100)]">
            {user.role === 'ADMIN' ? 'ADMIN' : `BATCH - ${user.batch}`}
          </p>
        </div>
      </div>
    </div>
  );
}

export default UserAvatar;
