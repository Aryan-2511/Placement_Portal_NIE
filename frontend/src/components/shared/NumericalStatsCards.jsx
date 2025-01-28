import { useActiveOpportunities } from '@/features/dashboard/useActiveOpportunities';
import NumericalStatCard from './NumericalStatCard';
import Spinner from './Spinner';
import { useRecentOpportunities } from '@/features/dashboard/useRecentOpportunities';
import { useTotalPlacedStudents } from '@/features/dashboard/useTotalPlacedStudents';
import { useTotalApplications } from '@/features/dashboard/useTotalApplications';

function NumericalStatsCards() {
  const { activeOpportunities, isActiveOpportunitiesLoading } =
    useActiveOpportunities();
  const { recentOpportunities, isRecentOpportunitiesLoading } =
    useRecentOpportunities();
  const { totalPlacedStudents, isTotalPlacedStudentsLoading } =
    useTotalPlacedStudents();
  const { totalApplications, isTotalApplicationsLoading } =
    useTotalApplications();
  if (
    isActiveOpportunitiesLoading ||
    isRecentOpportunitiesLoading ||
    isTotalPlacedStudentsLoading ||
    isTotalApplicationsLoading
  )
    return <Spinner />;
  return (
    <div className="w-full flex items-center justify-between gap-8">
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/bag.png'}
        title={'companies applied for'}
        color={'var(--color-blue-100)'}
        data={activeOpportunities.active_opportunities}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/verified.png'}
        title={'Recent Opportunities'}
        color={'var(--color-brown-100)'}
        data={recentOpportunities.recent_opportunities}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/dart.png'}
        title={'total placed in batch'}
        color={'var(--color-green-100)'}
        data={totalPlacedStudents.placed_students}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/bag.png'}
        title={'total applications'}
        color={'var(--color-blue-100)'}
        data={totalApplications.total_applications}
      />
    </div>
  );
}

export default NumericalStatsCards;
