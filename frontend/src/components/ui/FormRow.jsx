import { Label } from '@/components/ui/label';

function FormRow({ children, label, className }) {
  return (
    <div
      className={`max-h-[20rem] grid grid-cols-[24rem_2fr_1.2fr] items-center px-[2rem] mb-[3rem] ${className}`}
    >
      <Label className="">{label}</Label>
      {children}
    </div>
  );
}

export default FormRow;
