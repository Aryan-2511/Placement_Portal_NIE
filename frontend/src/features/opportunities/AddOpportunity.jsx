import Form from '@/components/shared/Form';
import FormRow from '@/components/ui/FormRow';
import { Input } from '@/components/ui/input';
// import { formFields } from './opportunityFields';
import { Button } from '@/components/ui/button';
import { useForm, Controller, useFieldArray } from 'react-hook-form';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { useUser } from '../authentication/useUser';
import { HiMiniXMark } from 'react-icons/hi2';
import useAddOpportunity from './useAddOpportunity';
import Spinner from '@/components/shared/Spinner';

function AddOpportunity() {
  const user = useUser();
  const { role } = user;
  const { mutate: addOpportunity, isLoading } = useAddOpportunity();
  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm({
    defaultValues: {
      allowed_branches: [],
      allowed_genders: [],
      coordinators: [],
      attached_documents: [''],
    },
  });
  const { fields, append, remove } = useFieldArray({
    control,
    name: 'attached_documents',
  });

  function onSubmit(data) {
    data.created_by = role;
    data.status = 'ACTIVE';
    data.completed = 'No';
    data.registration_date = new Date(data.registration_date).toISOString();
    console.log(data);
    if (isLoading) return <Spinner />;
    addOpportunity(data);
    reset();
  }
  const branches = [
    { id: 'cse', value: 'CSE' },
    { id: 'ise', value: 'ISE' },
    { id: 'ece', value: 'ECE' },
    { id: 'eee', value: 'EEE' },
    { id: 'mech', value: 'MECH' },
    { id: 'civ', value: 'CIV' },
  ];
  const genders = [
    { id: 'male', value: 'Male' },
    { id: 'female', value: 'Female' },
    { id: 'other', value: 'Other' },
  ];
  const coordinators = [
    { id: '1', value: { name: 'John Wick', contact: '9918656496' } },
    { id: '2', value: { name: 'Suryakumar Yadav', contact: '9918656496' } },
    { id: '3', value: { name: 'Virat Kohli', contact: '9918656496' } },
    { id: '4', value: { name: 'Steve Smith', contact: '9918656496' } },
    { id: '5', value: { name: 'James Anderson', contact: '9918656496' } },
    { id: '6', value: { name: 'Beth Mooney', contact: '9918656496' } },
  ];

  return (
    <div className="min-w-[102.4rem]">
      <h3>Add a new opporunity</h3>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <FormRow label="Company" error={errors?.company?.message}>
          <Input
            type="text"
            id="company"
            disabled={false}
            {...register('company', { required: 'Company is required!' })}
          />
        </FormRow>
        <FormRow label="Title" error={errors?.title?.message}>
          <Input
            type="text"
            id="title"
            disabled={false}
            {...register('title', { required: 'Title is required!' })}
          />
        </FormRow>
        <FormRow label="Location" error={errors?.location?.message}>
          <Input
            type="text"
            id="location"
            disabled={false}
            {...register('location', { required: 'Location is required!' })}
          />
        </FormRow>
        <FormRow
          label="Job description"
          error={errors?.job_description?.message}
        >
          <Textarea
            type="text"
            id="job_description"
            disabled={false}
            {...register('job_description', {
              required: 'Job description is required!',
            })}
          />
        </FormRow>
        <FormRow label="Category" error={errors?.category?.message}>
          <Input
            type="text"
            id="category"
            disabled={false}
            {...register('category', { required: 'Category is required!' })}
          />
        </FormRow>
        <FormRow
          label="Opportunity type"
          error={errors?.opportunity_type?.message}
        >
          <Input
            type="text"
            id="opportunity_type"
            disabled={false}
            {...register('opportunity_type', {
              required: 'Opportunity type is required!',
            })}
          />
        </FormRow>
        <FormRow label="Additional information" error={errors?.name?.message}>
          <Textarea
            type="text"
            id="additional_info"
            disabled={false}
            {...register('additional_info', {
              required: 'Additional information is required!',
            })}
          />
        </FormRow>
        <FormRow label="CTC" error={errors?.name?.message}>
          <Input
            type="number"
            step="any"
            id="ctc"
            disabled={false}
            {...register('ctc', {
              required: 'Location is required!',
              min: {
                value: 0,
                message: "CTC can't be negative",
              },
              setValueAs: (value) => parseFloat(value) || 0,
            })}
          />
        </FormRow>
        <FormRow label="CTC information" error={errors?.name?.message}>
          <Input
            type="text"
            id="ctc_info"
            disabled={false}
            {...register('ctc_info', {
              required: 'CTC information is required!',
            })}
          />
        </FormRow>

        <FormRow
          label="Allowed branches"
          error={errors?.allowed_branches?.message}
        >
          <div className="flex flex-wrap justify-between">
            {branches.map((branch) => (
              <div key={branch.id} className="flex items-center gap-[0.4rem]">
                <Controller
                  name="allowed_branches"
                  control={control}
                  render={({ field }) => (
                    <Checkbox
                      id={branch.id}
                      checked={field.value?.includes(branch.value) || false}
                      onCheckedChange={(checked) => {
                        // Toggle selection
                        const newValue = checked
                          ? [...(field.value || []), branch.value]
                          : (field.value || []).filter(
                              (v) => v !== branch.value
                            );

                        field.onChange(newValue);
                      }}
                    />
                  )}
                />
                <label htmlFor={branch.id}>{branch.value}</label>
              </div>
            ))}
          </div>
        </FormRow>
        <FormRow
          label="Allowed genders"
          error={errors?.allowed_genders?.message}
        >
          <div className="flex flex-wrap justify-between">
            {genders.map((gender) => (
              <div key={gender.id} className="flex items-center gap-[0.4rem]">
                <Controller
                  name="allowed_genders"
                  control={control}
                  render={({ field }) => (
                    <Checkbox
                      id={gender.id}
                      checked={field.value?.includes(gender.value) || false}
                      onCheckedChange={(checked) => {
                        // Toggle selection
                        const newValue = checked
                          ? [...(field.value || []), gender.value]
                          : (field.value || []).filter(
                              (v) => v !== gender.value
                            );

                        field.onChange(newValue);
                      }}
                    />
                  )}
                />
                <label htmlFor={gender.id}>{gender.value}</label>
              </div>
            ))}
          </div>
        </FormRow>

        <FormRow label="Coordinators" error={errors?.coordinators?.message}>
          <div className="flex flex-wrap justify-between">
            {coordinators.map((coordinator) => (
              <div
                key={coordinator.id}
                className="flex items-center gap-[0.4rem]"
              >
                <Controller
                  name="coordinators"
                  control={control}
                  render={({ field }) => (
                    <Checkbox
                      id={coordinator.id}
                      checked={
                        field.value?.some(
                          (c) => c.name === coordinator.value.name
                        ) || false
                      }
                      onCheckedChange={(checked) => {
                        // Toggle selection
                        const newValue = checked
                          ? [...(field.value || []), coordinator.value]
                          : (field.value || []).filter(
                              (v) => v.name !== coordinator.value.name
                            );

                        field.onChange(newValue);
                      }}
                    />
                  )}
                />
                <label htmlFor={coordinator.id}>{coordinator.value.name}</label>
              </div>
            ))}
          </div>
        </FormRow>

        <FormRow
          label="Class 10 percentage"
          error={errors?.class_10_percentage_criteria?.message}
        >
          <Input
            type="number"
            step="any"
            id="class_10_percentage_criteria"
            disabled={false}
            {...register('class_10_percentage_criteria', {
              required: 'Class 10 % is required!',
              min: 1,
              max: 100,
              setValueAs: (value) => parseFloat(value) || 0,
            })}
          />
        </FormRow>

        <FormRow label="Class 12 percentage" error={errors?.name?.message}>
          <Input
            type="number"
            step="any"
            id="class_12_percentage_criteria"
            disabled={false}
            {...register('class_12_percentage_criteria', {
              required: 'Class 12 % is required!',
              min: 1,
              max: 100,
              setValueAs: (value) => parseFloat(value) || 0,
            })}
          />
        </FormRow>

        <FormRow label="CGPA" error={errors?.name?.message}>
          <Input
            type="number"
            step="any"
            id="cgpa"
            disabled={false}
            {...register('cgpa', {
              required: 'CGPA is required!',
              min: 1,
              max: 10,
              setValueAs: (value) => parseFloat(value) || 0,
            })}
          />
        </FormRow>

        <FormRow label="Batch" error={errors?.name?.message}>
          <Input
            type="text"
            id="batch"
            disabled={false}
            {...register('batch', {
              required: 'Batch is required!',
            })}
          />
        </FormRow>

        <FormRow label="Backlogs" error={errors?.backlog?.message}>
          <Input
            type="number"
            id="backlog"
            disabled={false}
            {...register('backlog', {
              required: 'Backlog field is required!',
              min: 0,
              setValueAs: (value) => parseFloat(value),
            })}
          />
        </FormRow>

        <FormRow label="Registration date" error={errors?.name?.message}>
          <Input
            type="date"
            id="registration_date"
            disabled={false}
            {...register('registration_date', {
              required: 'Registration date is required!',
            })}
          />
        </FormRow>

        <FormRow label="Form link" error={errors?.name?.message}>
          <Input
            type="url"
            id="form_link"
            disabled={false}
            {...register('form_link', {
              required: 'Form link date is required!',
            })}
          />
        </FormRow>

        <FormRow label="Attached documents">
          <div className="flex flex-col gap-[1.2rem]">
            {fields.map((field, index) => (
              <div key={field.id} className="flex gap-2">
                <Input
                  type="url"
                  {...register(`attached_documents.${index}`, {
                    // required: 'URL is required!',
                  })}
                  placeholder="Enter document URL"
                />
                {/* <Button
                  variant="destructive"
                  onClick={() => remove(index)}
                ></Button> */}
                <HiMiniXMark
                  size={'2.4rem'}
                  onClick={() => remove(index)}
                  disabled={fields.length === 0}
                  className="cursor-pointer hover:bg-[var(--color-grey-50)]"
                />
              </div>
            ))}
            <Button
              className="self-start"
              onClick={(e) => {
                e.preventDefault();
                append('');
              }}
            >
              Add+
            </Button>
          </div>
        </FormRow>

        <Button className="self-start">Add Opportunity</Button>
      </Form>
    </div>
  );
}

export default AddOpportunity;
