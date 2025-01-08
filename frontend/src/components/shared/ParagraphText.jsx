function ParagraphText({ children }) {
  return (
    <p className="px-[1.2rem] text-[var(--color-grey-100)] text-[1.4rem] text-justify">
      {children}
    </p>
  );
}

export default ParagraphText;
