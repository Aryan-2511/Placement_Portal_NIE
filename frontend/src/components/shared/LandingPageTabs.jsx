import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import OurRecruiters from './OurRecruiters';
import PlacementStats from './PlacementStats';
import OurTeam from './OurTeam';
import Login from './Login';

function LandingPageTabs() {
  return (
    <div className="w-full max-w-[59.5rem]">
      <Tabs
        defaultValue="our-recruiters"
        className="text-[var(--color-grey-600)]"
      >
        <TabsList className="grid w-full grid-cols-4 gap-2 mb-[7.4rem]">
          <TabsTrigger value="our-recruiters">Our Recruiters</TabsTrigger>
          <TabsTrigger value="placement-statistics">
            Placement Stats
          </TabsTrigger>
          <TabsTrigger value="our-team">Our Team</TabsTrigger>
          <TabsTrigger value="login">Login</TabsTrigger>
        </TabsList>
        <TabsContent value="our-recruiters">
          <OurRecruiters />
        </TabsContent>
        <TabsContent value="placement-statistics">
          <PlacementStats />
        </TabsContent>
        <TabsContent value="our-team">
          <OurTeam />
        </TabsContent>
        <TabsContent value="login">
          <Login />
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default LandingPageTabs;
