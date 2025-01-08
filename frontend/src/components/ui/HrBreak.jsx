function HrBreak({ size = 'md' }) {
  const marginClasses = {
    sm: 'my-[0.8rem]',
    md: 'my-[2.4rem]',
    lg: 'my-[4.8rem]',
  };

  return <hr className={marginClasses[size] || marginClasses.md} />;
}

export default HrBreak;
