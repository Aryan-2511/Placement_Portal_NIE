import NumericalStatCard from './NumericalStatCard';

function NumericalStatsCards() {
  return (
    <div className="w-full flex items-center justify-between gap-8">
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/bag.png'}
        title={'companies applied for'}
        color={'var(--color-blue-100)'}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/verified.png'}
        title={'eligible for'}
        color={'var(--color-brown-100)'}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/dart.png'}
        title={'active opportunities'}
        color={'var(--color-green-100)'}
      />
      <NumericalStatCard
        path={'../../../public/icons/dashboard-icons/bag.png'}
        title={'lorem ipsum'}
        color={'var(--color-blue-100)'}
      />
    </div>
  );
}

export default NumericalStatsCards;
