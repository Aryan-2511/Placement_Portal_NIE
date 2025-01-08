import * as React from 'react';

import { cn } from '@/lib/utils';

const Input = React.forwardRef(({ className, type, ...props }, ref) => {
  return (
    <input
      type={type}
      className={cn(
        'flex min-h-9 w-full outline outline-1 outline-[var(--color-grey-100)] rounded-md border border-input bg-transparent px-3 py-2 text-[1.6rem] shadow-sm transition-colors file:border-0 file:bg-transparent file:text-2xl file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-[1.6rem]',
        className
      )}
      ref={ref}
      {...props}
    />
  );
});
Input.displayName = 'Input';

export { Input };
