import HrBreak from '@/components/ui/HrBreak';
import Student from './Student';

function AllStudentsTable({ students }) {
  return (
    <div className="pt-[1.6rem] px-[3.2rem] bg-[var(--color-grey-0)] mb-[1.2rem] shadow-md">
      <div className="grid grid-cols-[1.2fr_0.8fr_1.6fr_0.6fr_0.8fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] text-[1.4rem] font-semibold">
        <p>NAME</p>
        <p>USN</p>
        <p>COLLEGE EMAIL</p>
        <p>BRANCH</p>
        <p>CONTACT</p>
        <p>CGPA</p>
      </div>
      <HrBreak size="sm" />
      <div>
        {students.map((student) => {
          return <Student key={student.usn} student={student} />;
        })}
      </div>
    </div>
  );
}

export default AllStudentsTable;
