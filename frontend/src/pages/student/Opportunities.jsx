import Spinner from '@/components/shared/Spinner';
import { useUser } from '@/features/authentication/useUser';
import OpportunityTable from '@/features/opportunities/OpportunityTable';
import useOpportunities from '@/features/opportunities/useOpportunities';

function Opportunities() {
  const { batch } = useUser();
  const { opportunities, isLoading } = useOpportunities(batch);
  if (isLoading) return <Spinner />;

  return (
    <div className="min-w-[102.4rem]">
      <h3>Opportunities for the batch {batch}</h3>
      <OpportunityTable opportunities={opportunities} />
    </div>
  );
}

export default Opportunities;
