import ProfileDetails from '@/components/shared/ProfileDetails';
import { Button } from '@/components/ui/button';

function Profile() {
  return (
    <div className="min-w-[102.4rem]">
      <h3>Profile details</h3>
      {/* <div className="flex flex-col gap-[5.4rem]"> */}
      <ProfileDetails />
      <div className="flex justify-end px-[2rem] py-[1.2rem]">
        <Button variant="ghost" size="lg">
          Cancel
        </Button>
        <Button size="lg">Update</Button>
      </div>
    </div>
  );
}

export default Profile;
