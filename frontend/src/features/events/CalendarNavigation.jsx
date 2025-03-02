import { useState } from 'react';

function CalendarNavigation({ currentDate, onDateChange }) {
  const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth());
  const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());

  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];

  const years = Array.from(
    { length: 11 },
    (_, i) => currentDate.getFullYear() - 5 + i
  );

  const handleMonthChange = (e) => {
    const newMonth = parseInt(e.target.value, 10);
    setSelectedMonth(newMonth);
    onDateChange(new Date(selectedYear, newMonth, 1));
  };

  const handleYearChange = (e) => {
    const newYear = parseInt(e.target.value, 10);
    setSelectedYear(newYear);
    onDateChange(new Date(newYear, selectedMonth, 1));
  };

  return (
    <div className="flex gap-2 mb-4">
      <select
        value={selectedMonth}
        onChange={handleMonthChange}
        className="p-2 border border-[#AAAAAA] rounded-md bg-[#1A1A1A] text-[#FFFFFF]"
      >
        {months.map((month, index) => (
          <option key={index} value={index}>
            {month}
          </option>
        ))}
      </select>
      <select
        value={selectedYear}
        onChange={handleYearChange}
        className="p-2 border border-[#AAAAAA] rounded-md bg-[#1A1A1A] text-[#FFFFFF]"
      >
        {years.map((year) => (
          <option key={year} value={year}>
            {year}
          </option>
        ))}
      </select>
    </div>
  );
}

export default CalendarNavigation;
