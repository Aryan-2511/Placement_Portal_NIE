import FormRow from '@/components/ui/FormRow';
import { Input } from '@/components/ui/input';
import { useState } from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select';
import HrBreak from '../ui/HrBreak';
// import { getStudentDetails } from '@/services/apiAuth';

function ProfileDetails() {
  // getStudentDetails('4NI21IS010')
  //   .then((res) => console.log(res))
  //   .catch((err) => console.error(err));

  const [isEditable, setIsEditable] = useState(false);

  return (
    <form className="p-[3.2rem] bg-[var(--color-grey-0)] shadow-[var(--shadow-lg)]">
      <p className="font-semibold mb-[2.4rem]">Personal details</p>
      <FormRow label="Fullname" className="">
        <Input
          type="text"
          value="John Doe"
          onChange={() => console.log('profile section')}
          id="fullName"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="Contact" className="">
        <Input
          type="text"
          value="8765326651"
          onChange={() => console.log('profile section')}
          id="contact"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="Address" className="">
        <Input
          type="text"
          value="32 street, JP nagar, Karnataka"
          onChange={() => console.log('profile section')}
          id="address"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="Gender" className="">
        <Select defaultValue="male" disabled={!isEditable}>
          <SelectTrigger className="">
            <SelectValue placeholder="Select" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="male">Male</SelectItem>
            <SelectItem value="female">Female</SelectItem>
          </SelectContent>
        </Select>
      </FormRow>
      <HrBreak />
      <p className="font-semibold mb-[2.4rem]">Academic details</p>
      <FormRow label="USN" className="">
        <Input
          type="text"
          value="4NI21CS000"
          onChange={() => console.log('profile section')}
          id="usn"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="College email" className="">
        <Input
          type="text"
          value="8765326651"
          onChange={() => console.log('profile section')}
          id="contact"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="Address" className="">
        <Input
          type="text"
          value="32 street, JP nagar, Karnataka"
          onChange={() => console.log('profile section')}
          id="address"
          disabled={!isEditable}
        />
      </FormRow>
      <FormRow label="Branch" className="">
        <Select defaultValue="" className="">
          <SelectTrigger className="">
            <SelectValue placeholder="Select" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="cse">CSE</SelectItem>
            <SelectItem value="ise">ISE</SelectItem>
            <SelectItem value="ece">ECE</SelectItem>
            <SelectItem value="eee">EEE</SelectItem>
            <SelectItem value="mech">MECH</SelectItem>
            <SelectItem value="civ">CIV</SelectItem>
          </SelectContent>
        </Select>
      </FormRow>
    </form>
  );
}

export default ProfileDetails;
