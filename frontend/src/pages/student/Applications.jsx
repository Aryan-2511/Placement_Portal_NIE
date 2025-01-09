import { applications } from '../../../public/dummy-data/application/applicationData';
import ApplicationTable from '../../features/applications/ApplicationTable';
// import { useState } from 'react';

function Applications() {
  // const [batch, setBatch] = useState('2025');
  // const { opportunities, isLoading } = useOpportunities(batch);

  // below is temporary code just for making the layout

  // if (isLoading) return <Spinner />;

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
