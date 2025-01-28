import HrBreak from '@/components/ui/HrBreak';
import Opportunity from '@/features/opportunities/Opportunity';

function OpportunityTable({ opportunities }) {
  return (
    <div className="pt-[1.6rem] px-[3.2rem] bg-[var(--color-grey-0)] shadow-md">
      <div className="grid grid-cols-[1.6fr_1fr_0.6fr_1fr_0.4fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] text-[1.4rem] font-semibold">
        <p>NAME</p>
        <p>COMPENSATION(in ₹)</p>
        <p>TYPE</p>
        <p>REGISTRATION DATE</p>
        <p>STATUS</p>
      </div>
      <HrBreak size="sm" />
      <div>
        {opportunities.map((opportunity) => {
          return <Opportunity key={opportunity.id} opportunity={opportunity} />;
        })}
      </div>
    </div>
  );
}

export default OpportunityTable;
