import {
  Line,
  LineChart,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from 'recharts';

const data = [
  { name: 'Page A', uv: 400, pv: 2400, amt: 2400 },
  { name: 'Page B', uv: 300, pv: 2200, amt: 2000 },
  { name: 'Page C', uv: 200, pv: 2800, amt: 2400 },
  { name: 'Page D', uv: 278, pv: 2600, amt: 2300 },
  { name: 'Page E', uv: 189, pv: 2900, amt: 2500 },
];

function GraphData() {
  return (
    <LineChart
      width={1024}
      height={428}
      data={data}
      margin={{ top: 5, right: 36, bottom: 5, left: 0 }}
    >
      <Line type="monotone" dataKey="uv" stroke="#8884d8" />
      <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
      <XAxis dataKey="name" />
      <YAxis />
      <Tooltip />
    </LineChart>
  );
}

export default GraphData;
