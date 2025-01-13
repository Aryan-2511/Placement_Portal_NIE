import { useApplications } from '@/features/applications/useApplications';
import ApplicationTable from '../../features/applications/ApplicationTable';
import Spinner from '@/components/shared/Spinner';

function Applications() {
  const { applications, isLoading } = useApplications();

  if (isLoading) return <Spinner />;

  return (
    <div className="min-w-[102.4rem]">
      <div className="flex justify-between items-center">
        <h3>Applied opportunities by you</h3>
      </div>
      <ApplicationTable applications={applications} />
    </div>
  );
}

export default Applications;
