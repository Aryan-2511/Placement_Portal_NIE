import HrBreak from '@/components/ui/HrBreak';
import Application from './application';
function ApplicationTable({ applications }) {
  return (
    <div className="pt-[1.6rem] px-[3.2rem] bg-[var(--color-grey-0)] shadow-md">
      <div className="grid grid-cols-[0.8fr_1.6fr_1.4fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] text-[1.4rem] font-semibold">
        <p>OPPORTUNITY ID</p>
        <p>JOB POST</p>
        <p>COMPANY</p>
        <p>STATUS</p>
      </div>
      <HrBreak size="sm" />
      <div>
        {applications.map((application) => {
          return <Application key={application.id} application={application} />;
        })}
      </div>
    </div>
  );
}

export default ApplicationTable;
