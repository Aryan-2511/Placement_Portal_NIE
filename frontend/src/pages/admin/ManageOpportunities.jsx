import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

import OpportunityTable from '@/features/opportunities/OpportunityTable';
import { useState } from 'react';
import useOpportunities from '@/features/opportunities/useOpportunities';
import Spinner from '@/components/shared/Spinner';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

function ManageOpportunities() {
  const [batch, setBatch] = useState('2025');
  const { opportunities, isLoading } = useOpportunities(batch);
  if (isLoading) return <Spinner />;
  return (
    <div className="min-w-[102.4rem]">
      <div className="flex justify-between items-center">
        <h3>Opportunities for the batch {batch}</h3>
        <Select onValueChange={(value) => setBatch(value)} defaultValue={batch}>
          <SelectTrigger className="w-[10rem]">
            <SelectValue placeholder="Batch" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="2023">2023</SelectItem>
            <SelectItem value="2024">2024</SelectItem>
            <SelectItem value="2025">2025</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <OpportunityTable opportunities={opportunities} />
      <Link to="add_new_opportunity">
        <Button>Add new opportunity</Button>
      </Link>
    </div>
  );
}

export default ManageOpportunities;
