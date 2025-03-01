import { Label } from '@/components/ui/label';

function FormRow({ children, label, className, columns = 3 }) {
  const columnClass =
    columns === 2
      ? 'grid-cols-[1.5fr_2fr]'
      : columns === 3
      ? 'grid-cols-[24rem_2fr_1.2fr]'
      : columns === 4
      ? 'grid-cols-[24rem_2fr_1.2fr_1fr]'
      : '';

  return (
    <div
      className={`max-h-[20rem] grid ${columnClass} items-center px-[2rem] ${className}`}
    >
      <Label className="">{label}</Label>
      {children}
    </div>
  );
}

export default FormRow;
