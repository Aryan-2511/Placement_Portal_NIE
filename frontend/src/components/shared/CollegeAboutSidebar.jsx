function CollegeAboutSidebar({ collegeName, collegeInfo }) {
  // capitalize college name
  collegeName = collegeName
    .split(' ') // Split the string into an array of words
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1)) // Capitalize the first letter of each word
    .join(' ');

  return (
    <div className={'w-[75%] mx-auto mt-[17.2rem] md:w-[80%] pb-6'}>
      <p className="text-[1.6rem] text-[#BCBCBC]">
        Welcome to the placement portal of,
      </p>
      <h3 className="text-[2rem] text-[#EAEAEA] leading-[3.6rem] mb-[1.6rem] font-semibold md:text-[3.2rem] md:leading-[4.2rem]">
        {collegeName}
      </h3>
      <p className="text-justify text-[#D9D9D9]">{collegeInfo}</p>
    </div>
  );
}

export default CollegeAboutSidebar;
