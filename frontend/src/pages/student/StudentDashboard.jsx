import GraphSummary from '../../components/shared/GraphSummary';
import NumericalStatsCards from '../../components/shared/NumericalStatsCards';
import OpportunitiesSummary from '../../components/shared/OpportunitiesSummary';

function StudentDashboard() {
  return (
    <div className="min-w-[102.4rem]">
      <h3 className="text-[2.4rem] font-semibold text-[var(--color-grey-600)] mb-9">
        Your Stats
      </h3>
      <div className="flex flex-col gap-[5.4rem]">
        <NumericalStatsCards />
        <OpportunitiesSummary />
        <GraphSummary />
      </div>
    </div>
  );
}

export default StudentDashboard;
