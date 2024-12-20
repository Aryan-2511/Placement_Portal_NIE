function Logo() {
  return (
    <div className="pt-[3.4rem] pb-[6.2rem] text-center">
      <div className=" mx-[auto] my-0 w-[4.6rem] h-[4.6rem] text-center p-2 border-2 border-[var(--color-blue-700)] rounded-full flex align-center justify-center">
        <img
          src="../../../public/images/nie.png"
          alt="nie logo"
          className="w-[2.4rem]"
        />
      </div>
      <p className="text-[1.6rem] mt-2 font-bold text-[var(--color-grey-600)]">
        The National Institute of Engineering, Mysuru{' '}
      </p>
    </div>
  );
}

export default Logo;
