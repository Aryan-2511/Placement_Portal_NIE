import { useNavigate } from 'react-router-dom';

function Student({ student }) {
  const navigate = useNavigate();
  return (
    <div
      className="text-[1.4rem] grid grid-cols-[1.2fr_0.8fr_1.6fr_0.6fr_0.8fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] border-b-[.1rem] cursor-pointer hover:bg-slate-50"
      onClick={() => {
        navigate(`${student.usn}`);
      }}
    >
      <p>{student.name}</p>
      <p>{student.usn}</p>
      <p>{student.college_email}</p>
      <p>{student.branch}</p>
      <p>{student.contact}</p>
      <p>{student.current_cgpa}</p>
    </div>
  );
}

export default Student;
