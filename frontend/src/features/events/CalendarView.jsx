import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import Spinner from '@/components/shared/Spinner';
import useEvents from './useEvents';
import { useEffect, useRef } from 'react';

function CalendarView({ selectedBatch, currentDate }) {
  const { data: allEvents, isLoading, error } = useEvents();
  const calendarRef = useRef(null);

  useEffect(() => {
    if (calendarRef.current) {
      calendarRef.current.getApi().gotoDate(currentDate);
    }
  }, [currentDate]);

  if (isLoading) return <Spinner />;
  if (error) return <div className="text-red-500">Error: {error.message}</div>;

  const filteredEvents = allEvents?.filter(
    (event) => event.batch === selectedBatch
  );

  const calendarEvents = filteredEvents?.map((event) => ({
    id: event.schedule_id,
    title: event.title,
    start: event.start_time,
    end: event.end_time,
    extendedProps: {
      description: event.description,
      created_by: event.created_by,
      batch: event.batch,
    },
  }));

  return (
    <div className="bg-white p-4 rounded-md shadow-md">
      <FullCalendar
        ref={calendarRef}
        plugins={[dayGridPlugin]}
        initialView="dayGridMonth"
        initialDate={currentDate}
        events={calendarEvents}
        headerToolbar={{
          left: 'prev,next',
          center: 'title',
          right: '',
        }}
        height="auto"
        eventDidMount={(info) => {
          info.el.title = info.event.extendedProps.description || '';
        }}
      />
    </div>
  );
}

export default CalendarView;
