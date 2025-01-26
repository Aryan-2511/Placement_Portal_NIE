import GraphSummary from '../../features/dashboard/GraphSummary';
import NumericalStatsCards from '../../components/shared/NumericalStatsCards';
import OpportunitiesSummary from '../../components/shared/OpportunitiesSummary';

function StudentDashboard() {
  return (
    <div className="min-w-[102.4rem]">
      <h3>Your Stats</h3>
      <div className="flex flex-col gap-[5.4rem]">
        <NumericalStatsCards />
        <OpportunitiesSummary />
        <GraphSummary />
      </div>
    </div>
  );
}

export default StudentDashboard;
