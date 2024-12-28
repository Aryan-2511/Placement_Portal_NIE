function UserAvatar() {
  return (
    <div className="font-['Poppins']">
      <div className="flex justify-center gap-3">
        <div className="w-[4rem] h-[4rem] bg-slate-200 rounded-full"></div>
        <div>
          <p className="text-[1.6rem] text-[var(--color-grey-600) font-medium">
            John Doe
          </p>
          <p className="text-[1.2rem] text-[var(--color-grey-100)]">
            7th Semester
          </p>
        </div>
      </div>
    </div>
  );
}

export default UserAvatar;
