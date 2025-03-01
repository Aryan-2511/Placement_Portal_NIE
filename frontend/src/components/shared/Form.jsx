function Form({ onSubmit, children }) {
  return (
    <form
      onSubmit={onSubmit}
      className="w-full p-[3.2rem] bg-[var(--color-grey-0)] shadow-[var(--shadow-lg)] flex flex-col gap-[3.2rem]"
    >
      {children}
    </form>
  );
}

export default Form;
