import { useNavigate } from 'react-router-dom';

function Application({ application }) {
  const navigate = useNavigate();
  return (
    <div
      onClick={() =>
        navigate(`/student/opportunities/${application.opportunity_id}`)
      }
      className="grid grid-cols-[1.2fr_1.2fr_1.8fr_0.6fr] gap-[1.2rem] items-center px-[2rem] py-[1.2rem] border-b-[.1rem] cursor-pointer hover:bg-slate-50"
    >
      <p>{application.student_name}</p>
      <p className="text-[1.4rem]">{application.opportunity_id}</p>
      <div>
        <p className="font-semibold text-[var(--color-grey-600)]">
          {application.job_post}
        </p>
        <p className="text-[1.4rem] text-[var(--color-grey-100)]">
          {application.company}
        </p>
      </div>
      <p
        className={`px-4 py-1 font-semibold text-[1.4rem] text-[var(--color-brand-700)] rounded-xl text-center ${
          application.status === 'PROCESSING' && 'bg-yellow-200'
        } ${application.status === 'ACTIVE' ? 'bg-green-200' : 'bg-red-200'}`}
      >
        {application.status}
      </p>
    </div>
  );
}

export default Application;
