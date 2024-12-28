import GraphData from '../../components/shared/GraphData';

function GraphSummary() {
  return (
    <div className="bg-red-100 py-[3rem] px-[3.6rem] shadow-[var(--shadow-lg)]">
      <h3 className="text-[2rem] font-semibold mb-[1.6rem]">
        Monthwise Frequency
      </h3>
      <GraphData />
    </div>
  );
}

export default GraphSummary;
