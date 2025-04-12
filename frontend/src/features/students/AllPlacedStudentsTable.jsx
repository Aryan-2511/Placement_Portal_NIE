import HrBreak from '@/components/ui/HrBreak';
import { placedStudents } from '../../../public/dummy-data/student/studentData';
import PlacedStudent from './PlacedStudent';

function AllPlacedStudentsTable() {
  return (
    <div className="pt-[1.6rem] px-[3.2rem] bg-[var(--color-grey-0)] mb-[1.2rem] shadow-md">
      <div className="grid grid-cols-[1.2fr_0.8fr_1fr_0.6fr_0.8fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] text-[1.4rem] font-semibold">
        <p>NAME</p>
        <p>USN</p>
        <p>COMPANY</p>
        <p>BRANCH</p>
        <p>PLACEMENT TYPE</p>
        <p>PACKAGE</p>
      </div>
      <HrBreak size="sm" />
      <div>
        {placedStudents.map((student) => {
          return <PlacedStudent key={student.usn} student={student} />;
        })}
      </div>
    </div>
  );
}

export default AllPlacedStudentsTable;
