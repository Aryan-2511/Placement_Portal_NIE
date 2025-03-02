import Form from '@/components/shared/Form';
import FormRow from '@/components/ui/FormRow';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useForm } from 'react-hook-form';
import { useUser } from '../authentication/useUser';
import Spinner from '@/components/shared/Spinner';
import useAddEvent from './useAddEvent';

function AddEvent({ selectedBatch }) {
  const { role } = useUser();
  const { mutate: addEvent, isLoading } = useAddEvent();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    defaultValues: {
      batch: selectedBatch,
    },
  });

  const batches = ['2023', '2024', '2025', '2026'];

  function onSubmit(data) {
    data.start_time = new Date(data.start_time).toISOString();
    data.end_time = new Date(data.end_time).toISOString();
    data.created_by = role;
    if (isLoading) return;
    addEvent(data, {
      onSuccess: () => reset(),
    });
  }

  if (isLoading) return <Spinner />;

  return (
    <div>
      <h3 className="text-[2.4rem] font-semibold text-[var(--color-grey-600)] mb-6">
        Add a New Event
      </h3>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <FormRow label="Title" error={errors?.title?.message}>
          <Input
            type="text"
            id="title"
            {...register('title', { required: 'Title is required!' })}
          />
        </FormRow>

        <FormRow label="Description" error={errors?.description?.message}>
          <Textarea
            id="description"
            {...register('description', {
              required: 'Description is required!',
            })}
          />
        </FormRow>

        <FormRow label="Start Time" error={errors?.start_time?.message}>
          <Input
            type="datetime-local"
            id="start_time"
            {...register('start_time', { required: 'Start time is required!' })}
          />
        </FormRow>

        <FormRow label="End Time" error={errors?.end_time?.message}>
          <Input
            type="datetime-local"
            id="end_time"
            {...register('end_time', { required: 'End time is required!' })}
          />
        </FormRow>

        <FormRow label="Batch" error={errors?.batch?.message}>
          <select
            id="batch"
            className="p-2 border border-[var(--color-grey-100)] rounded-md bg-white w-full"
            {...register('batch', { required: 'Batch is required!' })}
          >
            {batches.map((batch) => (
              <option key={batch} value={batch}>
                {batch}
              </option>
            ))}
          </select>
        </FormRow>

        <Button type="submit" disabled={isLoading}>
          Add Event
        </Button>
      </Form>
    </div>
  );
}

export default AddEvent;
