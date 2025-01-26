import Logo from '@/components/shared/Logo';
import CollegeAboutSidebar from '../components/shared/CollegeAboutSidebar';
import LandingPageTabs from '@/components/shared/LandingPageTabs';

const collegeData = {
  collegeName: 'the national instiute of engineering, 570008, mysuru',
  collegeInfo:
    'The National Institute of Engineering (NIE) is a grant-in-aid institution and approved by the All India Council for Technical Education (AICTE), New Delhi. NIE got autonomous status from Visvesvaraya Technological University, Belagavi in 2007. It has been accredited by NAAC.',
};

function LandingPage() {
  return (
    <div className="h-screen flex flex-col md:grid md:grid-cols-[2fr_3fr]">
      <div className="bg-[var(--color-blue-700)] text-[var(--color-grey-0)] bg-[url(/public/images/nie.png)] bg-no-repeat bg-cover bg-[center_right_8rem] bg-blend-overlay md:bg-[center_right_32rem]">
        <CollegeAboutSidebar
          collegeName={collegeData.collegeName}
          collegeInfo={collegeData.collegeInfo}
        />
      </div>
      <div className="bg-[var(--color-grey-50)] flex flex-col items-center px-[1.2rem] sm:px-[3.2rem] md:px-[4.8rem]">
        <Logo />
        <LandingPageTabs />
      </div>
    </div>
  );
}

export default LandingPage;
