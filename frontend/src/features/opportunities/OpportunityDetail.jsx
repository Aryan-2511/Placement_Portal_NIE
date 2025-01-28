import { HiMapPin, HiCurrencyDollar, HiMiniUserGroup } from 'react-icons/hi2';
import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Input } from '@/components/ui/input';
import HeadingText from '@/components/shared/HeadingText';
import ParagraphText from '@/components/shared/ParagraphText';
import Spinner from '@/components/shared/Spinner';
import { Button } from '@/components/ui/button';
import HrBreak from '@/components/ui/HrBreak';
import useOpportunity from './useOpportunity';
import dateFormatter from '@/utils/dateFormatter';
import useUpdateOpportunity from './useUpdateOpportunity';
import { fields } from './opportunityFields';
import { useUser } from '../authentication/useUser';
import { useApplyOpportunity } from '../applications/useApplyOpportunity';
// import toast from 'react-hot-toast';

function OpportunityDetail() {
  const [isEditable, setIsEditable] = useState(false);
  const { opportunity, isOpportunityLoading } = useOpportunity();
  const { updateOpportunity, isUpdating } = useUpdateOpportunity();
  const { apply, isApplicationLoading } = useApplyOpportunity();
  const [opportunityData, setOpportunityData] = useState({});
  const user = useUser();
  const { role } = user;
  const { opportunityId } = useParams();

  useEffect(() => {
    setOpportunityData(opportunity);
  }, [opportunity]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setOpportunityData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleApply = () => {
    apply();
    if (isApplicationLoading) return <p>Loading...</p>;
  };

  const handleSubmit = () => {
    // Submit logic here
    opportunityData.ctc = parseFloat(opportunityData.ctc);
    opportunityData.cgpa = parseFloat(opportunityData.cgpa);
    opportunityData.class_10_percentage_criteria = parseFloat(
      opportunityData.class_10_percentage_criteria
    );
    opportunityData.class_12_percentage_criteria = parseFloat(
      opportunityData.class_12_percentage_criteria
    );
    if (!Array.isArray(opportunityData.allowed_branches)) {
      opportunityData.allowed_branches =
        opportunityData.allowed_branches?.split(',');
    }
    if (!Array.isArray(opportunityData.allowed_genders)) {
      opportunityData.allowed_genders =
        opportunityData.allowed_genders?.split(',');
    }

    updateOpportunity(
      { opportunityId, opportunityData, role },
      {
        onSuccess: (updatedOpportunity) => {
          setOpportunityData((prevData) => ({
            ...prevData,
            ...updatedOpportunity,
          }));
          setIsEditable(false);
          // console.log('successfull updation');
        },
        onError: (error) => {
          console.error('Failed to update opportunity:', error);
        },
      }
    );
  };

  const handleCancel = () => {
    setIsEditable((prev) => !prev);
    setOpportunityData(opportunity);
  };

  if (isOpportunityLoading || !opportunityData || isUpdating)
    return <Spinner />;

  return (
    <div>
      <h3>Opportunity #{opportunity.id}</h3>
      <div className="shadow-lg">
        {/* Opportunity Header */}
        <div className="px-[3.2rem] py-[2.4rem] flex justify-between items-center bg-gradient-to-r from-[var(--color-brand-700)] to-[#6D12AF] text-[var(--color-white)]">
          <div className="">
            <p className="text-[var(--color-white)] font-semibold">
              {opportunity.company}
            </p>
            <p className="text-[var(--color-white)]">{opportunity.title}</p>
          </div>
          <div className="flex justify-end gap-[6.4rem]">
            <div className="flex flex-col items-center">
              <HiMapPin />
              <span>{opportunity.location}</span>
            </div>
            <div className="flex flex-col items-center">
              <HiCurrencyDollar />
              <span>₹{opportunity.ctc}</span>
            </div>
            <div className="flex flex-col items-center">
              <HiMiniUserGroup />
              <span>{opportunity.batch}</span>
            </div>
          </div>
        </div>

        {/* Opportunity Fields */}
        <div className="px-[3.2rem] py-[2.4rem] flex flex-col gap-[3rem] bg-[var(--color-grey-0)]">
          {fields.map(({ name, label, type }) => (
            <div key={name}>
              <HeadingText>{label}</HeadingText>
              {isEditable ? (
                type === 'textarea' ? (
                  <textarea
                    name={name}
                    value={opportunityData[name]}
                    onChange={handleInputChange}
                    className="w-full p-2 border rounded"
                  />
                ) : (
                  <Input
                    type={type}
                    name={name}
                    value={opportunityData[name]}
                    onChange={handleInputChange}
                  />
                )
              ) : (
                <ParagraphText>
                  {name === 'form_link' ? (
                    <Link to={opportunityData[name]}>
                      {opportunityData[name]}
                    </Link>
                  ) : name === 'registration_date' ? (
                    dateFormatter(opportunityData[name], 'string')
                  ) : (
                    opportunityData[name]
                  )}
                </ParagraphText>
              )}
            </div>
          ))}

          {role === 'ADMIN' && (
            <>
              <HrBreak size="sm" />
              <div>
                <HeadingText>Created by</HeadingText>
                <ParagraphText>{opportunity.created_by}</ParagraphText>
              </div>
            </>
          )}
        </div>
      </div>

      {/* Buttons */}
      <div className="flex gap-[1.2rem] py-[1.2rem] justify-end">
        {role === 'STUDENT' && (
          <Button
            disabled={opportunity.status !== 'ACTIVE'}
            onClick={handleApply}
          >
            Apply
          </Button>
        )}
        {role === 'ADMIN' && opportunity.status === 'ACTIVE' && (
          <>
            <Button
              disabled={isEditable}
              onClick={() => setIsEditable((curr) => !curr)}
            >
              Edit details
            </Button>
            <Button
              variant="destructive"
              disabled={!isEditable}
              onClick={handleCancel}
            >
              Cancel
            </Button>
          </>
        )}
        {isEditable && opportunity.status === 'ACTIVE' && (
          <Button onClick={handleSubmit}>Save changes</Button>
        )}
      </div>
    </div>
  );
}

export default OpportunityDetail;
