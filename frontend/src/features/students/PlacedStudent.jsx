function PlacedStudent({ student }) {
  return (
    <div className="text-[1.4rem] grid grid-cols-[1.2fr_0.8fr_1fr_0.6fr_0.8fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] border-b-[.1rem] cursor-pointer hover:bg-slate-50">
      <p>{student.name}</p>
      <p>{student.usn}</p>
      <p>{student.company}</p>
      <p>{student.branch}</p>
      <p>{student.placement_type}</p>
      <p>{student.package}</p>
    </div>
  );
}

export default PlacedStudent;
