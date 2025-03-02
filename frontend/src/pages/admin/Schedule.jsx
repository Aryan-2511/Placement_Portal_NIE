import { useState } from 'react';
// import CalendarNavigation from './CalendarNavigation';
import CalendarNavigation from '@/features/events/CalendarNavigation';
import CalendarView from '@/features/events/CalendarView';
import AddEvent from '@/features/events/AddEvent';

function Schedule() {
  const [activeTab, setActiveTab] = useState('calendar');
  const [selectedBatch, setSelectedBatch] = useState('2025');
  // Step 1: Set up state for the calendar date
  const [currentDate, setCurrentDate] = useState(new Date(2023, 0, 1)); // January 2023
  const batches = ['2023', '2024', '2025', '2026'];

  // Step 2: Define the onDateChange function
  const handleDateChange = (newDate) => {
    setCurrentDate(newDate);
  };

  return (
    <div className="min-w-[102.4rem] p-4">
      {/* Header */}
      <div className="header bg-[var(--color-grey-50)] p-4 flex justify-between items-center rounded-md">
        <h1 className="text-[2.4rem] font-semibold text-[var(--color-grey-600)]">
          Schedule for Batch {selectedBatch}
        </h1>
        <select
          value={selectedBatch}
          onChange={(e) => setSelectedBatch(e.target.value)}
          className="p-2 border border-[var(--color-grey-100)] rounded-md bg-white text-[var(--color-grey-600)]"
        >
          {batches.map((batch) => (
            <option key={batch} value={batch}>
              {batch}
            </option>
          ))}
        </select>
      </div>

      {/* Step 3: Pass onDateChange to CalendarNavigation */}
      <CalendarNavigation
        currentDate={currentDate}
        onDateChange={handleDateChange}
      />

      {/* Sub-tabs */}
      <div className="sub-tabs mt-4 flex gap-2">
        <button
          onClick={() => setActiveTab('add')}
          className={`py-2 px-4 rounded-md font-medium ${
            activeTab === 'add'
              ? 'bg-[var(--color-blue-600)] text-white'
              : 'bg-[var(--color-grey-50)] text-[var(--color-grey-600)] hover:bg-[var(--color-grey-100)]'
          }`}
        >
          Add Event
        </button>
        <button
          onClick={() => setActiveTab('calendar')}
          className={`py-2 px-4 rounded-md font-medium ${
            activeTab === 'calendar'
              ? 'bg-[var(--color-blue-600)] text-white'
              : 'bg-[var(--color-grey-50)] text-[var(--color-grey-600)] hover:bg-[var(--color-grey-100)]'
          }`}
        >
          Calendar View
        </button>
      </div>

      {/* Step 4: Pass currentDate to CalendarView */}
      <div className="mt-4">
        {activeTab === 'add' ? (
          <AddEvent selectedBatch={selectedBatch} />
        ) : (
          <CalendarView
            selectedBatch={selectedBatch}
            currentDate={currentDate}
          />
        )}
      </div>
    </div>
  );
}

export default Schedule;
