// import Application from '@/features/applications/Application';
import { useApplications } from '@/features/applications/useApplications';
import Spinner from './Spinner';

function AppliedOpportunities() {
  const {
    applications,
    error,
    isLoading: isApplicationsLoading,
  } = useApplications();
  if (isApplicationsLoading) return <Spinner />;
  return (
    <div className="min-h-[34.6rem] flex-1 py-[3rem] px-[3.6rem] shadow-md bg-[var(--color-grey-0)]">
      <p className="pb-[2rem]">Your applications</p>
      <div className="grid grid-cols-[1.2fr_1.6fr_0.8fr] gap-[1.6rem] border-b-[0.1rem] text-[1.4rem] font-semibold">
        <p>OPP. ID</p>
        <p>COMPANY</p>
        <p>STATUS</p>
      </div>
      {applications.map((application) => {
        return (
          <div
            key={application.id}
            className="w-full grid grid-cols-[1.2fr_1.6fr_0.8fr] gap-[1.6rem] hover:bg-[var(--color-grey-50)] px-[0.4rem] py-[1.2rem] text-[1.4rem]"
          >
            <p>{application.opportunity_id}</p>
            <p>{application.company}</p>
            <p>{application.status}</p>
          </div>
        );
      })}
    </div>
  );
}

export default AppliedOpportunities;
