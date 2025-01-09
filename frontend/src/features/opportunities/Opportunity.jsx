import dateFormatter from '@/utils/dateFormatter';
import { useNavigate } from 'react-router-dom';

function Opportunity({ opportunity }) {
  const navigate = useNavigate();
  return (
    <div
      onClick={() => navigate(`${opportunity.id}`)}
      className="grid grid-cols-[2fr_1.2fr_1.2fr_0.4fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] border-b-[.1rem] cursor-pointer hover:bg-slate-50"
      key={opportunity.id}
    >
      <div>
        <p className="text-[1.6rem] font-semibold">{opportunity.company}</p>
        <p className="text-[1.4rem] font-normal">{opportunity.title}</p>
      </div>
      <p className="text-[1.4rem]">{opportunity.ctc}</p>
      <p className="text-[1.4rem]">
        {dateFormatter(opportunity.registration_date)}
      </p>
      <p className="text-[1.4rem] text-center">
        <span
          className={`font-semibold px-4 py-1 rounded-xl text-[1.4rem] ${
            opportunity.status === 'ACTIVE' ? 'bg-green-200' : 'bg-red-200'
          }`}
        >
          {opportunity.status.toUpperCase()}
        </span>
      </p>
    </div>
  );
}

export default Opportunity;
