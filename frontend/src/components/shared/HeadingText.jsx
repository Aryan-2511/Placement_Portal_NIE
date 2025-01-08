function HeadingText({ children, className }) {
  return (
    <p
      className={`font-semibold text-[1.6rem] text-[var(--color-grey-600)] mb-[0.6rem] ${className}`}
    >
      {children}
    </p>
  );
}

export default HeadingText;
