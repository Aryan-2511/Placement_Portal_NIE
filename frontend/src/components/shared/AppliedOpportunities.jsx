import Application from '@/features/applications/Application';
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
      <p className="pb-[2rem] border-b-[0.1rem]">Your application summary</p>
      {applications.map((application) => {
        return (
          <div
            key={application.id}
            className="w-full flex justify-between hover:bg-[var(--color-grey-50)] px-[0.4rem] py-[1.2rem]"
          >
            <p>{application.opportunity_id}</p>
            <p className="font-semibold">{application.company}</p>
            <p>{application.status}</p>
          </div>
        );
      })}
    </div>
  );
}

export default AppliedOpportunities;
