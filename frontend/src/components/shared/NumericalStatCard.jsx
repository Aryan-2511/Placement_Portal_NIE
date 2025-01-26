function NumericalStatCard({ path, title, color, data }) {
  return (
    <div className="h-[10.3rem] min-w-[20rem] bg-[var(--color-grey-0)] flex-1 p-3 flex items-center gap-[1.2rem] shadow-sm">
      <div
        className="h-[5.72rem] min-w-[5.72rem] flex items-center rounded-full"
        style={{ backgroundColor: color }}
      >
        <img src={path} alt="bag" className="w-[3rem] mx-auto" />
      </div>
      <div className="pt-3">
        <p className="text-[1.2rem] font-semibold text-[var(--color-grey-100)]">
          {title.toUpperCase()}
        </p>
        <p className="text-[2.4rem] font-semibold">{data}</p>
      </div>
    </div>
  );
}

export default NumericalStatCard;
