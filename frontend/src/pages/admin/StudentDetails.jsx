import { useState } from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import AllPlacedStudentsTable from '@/features/students/AllPlacedStudentsTable';
import AllStudentsTable from '@/features/students/AllStudentsTable';
import useAllStudents from '../../features/students/useAllStudents';
import Spinner from '@/components/shared/Spinner';
import { Button } from '@/components/ui/button';

function StudentDetails() {
  const [isByBatchActive, setIsByBatchActive] = useState(true);
  const [batch, setBatch] = useState('2025');
  const { allStudents, isLoading } = useAllStudents(batch);
  if (isLoading) return <Spinner />;
  function handleTabChange() {
    setIsByBatchActive((curr) => !curr);
  }
  function handleOnAddPlacedStudent() {
    // open the modal window
  }
  return (
    <div className="min-w-[102.4rem]">
      {/* conditional rendering based on all students by
       batch or all placed students tab is active */}
      <div className="flex justify-between items-center">
        <h3>
          {isByBatchActive ? 'All students by batch' : 'All placed students'}
        </h3>
        {isByBatchActive && (
          <Select
            onValueChange={(value) => setBatch(value)}
            defaultValue={batch}
          >
            <SelectTrigger className="w-[10rem]">
              <SelectValue placeholder="Batch" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="2023">2023</SelectItem>
              <SelectItem value="2024">2024</SelectItem>
              <SelectItem value="2025">2025</SelectItem>
            </SelectContent>
          </Select>
        )}
      </div>
      <Tabs defaultValue="by-batch" onValueChange={handleTabChange}>
        <TabsList>
          <TabsTrigger value="by-batch">By Batch</TabsTrigger>
          <TabsTrigger value="all-placed">All Placed</TabsTrigger>
        </TabsList>
        <TabsContent value="by-batch">
          <AllStudentsTable students={allStudents} />
        </TabsContent>
        <TabsContent value="all-placed">
          <AllPlacedStudentsTable />
        </TabsContent>
      </Tabs>
      {!isByBatchActive && (
        <div className="text-end">
          <Button onClick={handleOnAddPlacedStudent}>Add Placed Student</Button>
        </div>
      )}
    </div>
  );
}

export default StudentDetails;
