import AppliedOpportunities from './AppliedOpportunities';
import ToApplyOpportunities from './ToApplyOpportunities';

function OpportunitiesSummary() {
  return (
    <div className="flex gap-8">
      <AppliedOpportunities />
      <ToApplyOpportunities />
    </div>
  );
}

export default OpportunitiesSummary;
