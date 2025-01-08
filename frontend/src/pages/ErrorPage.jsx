import Lottie from 'lottie-react';
import hangingAnimation from '../assets/hangingAnimation.json';
function ErrorPage() {
  return (
    <div className="h-screen flex flex-col items-center justify-center">
      <div className="h-[40rem] w-[40rem]">
        <Lottie loop={true} animationData={hangingAnimation} />
      </div>
      <h3>Oops! Looks like we encountered an error</h3>
      <p>
        The page you&apos;re looking for can&apos;t be found. Make sure
        you&apos;ve entered the correct URL and try again
      </p>
    </div>
  );
}

export default ErrorPage;
