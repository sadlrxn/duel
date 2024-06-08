import React from 'react';

import { Button } from 'components';

export default function ModeButton({
  mode = '',
  selected = false,
  ...props
}: any) {
  return (
    <Button
      variant="secondary"
      height="38px"
      border={'1px solid #4FFF8B'}
      borderWidth={selected ? '1px' : '0px'}
      fontWeight={600}
      fontSize="14px"
      color={selected ? 'success' : 'text'}
      {...props}
    >
      {mode}
    </Button>
  );
}
