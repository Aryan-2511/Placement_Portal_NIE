import GraphData from './GraphData';

function GraphSummary() {
  return (
    <div className="bg-[var(--color-grey-0)] py-[3rem] px-[3.6rem] shadow-md">
      <h3 className="text-[2rem] font-semibold mb-[1.6rem]">
        Monthwise Frequency
      </h3>
      <GraphData />
    </div>
  );
}

export default GraphSummary;
