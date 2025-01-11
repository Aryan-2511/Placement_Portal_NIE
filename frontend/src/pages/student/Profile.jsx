import ProfileDetails from '@/components/shared/ProfileDetails';
import { Button } from '@/components/ui/button';
import { useQuery } from '@tanstack/react-query';

function Profile() {
  const { data: user } = useQuery(['user'], { initialData: null });
  console.log(user);
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
      {/* </div> */}
    </div>
  );
}

export default Profile;
