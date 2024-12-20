import NumericalStatsCards from '../../components/shared/NumericalStatsCards';

function StudentDashboard() {
  return (
    <div>
      <h3 className="text-[2.4rem] font-semibold text-[var(--color-grey-600)] mb-9">
        Your Stats
      </h3>
      <div>
        <NumericalStatsCards />
      </div>
    </div>
  );
}

export default StudentDashboard;
