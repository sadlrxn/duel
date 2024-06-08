import React from 'react';

import { BetButton as BaseBetButton } from 'components/GameTabs/styles';

export default function BetButton({ ...props }: any) {
  return (
    <BaseBetButton
      mt="8px"
      borderRadius="10px"
      fontSize="1em"
      fontWeight={700}
      color="black"
      {...props}
    >
      Place Bet
    </BaseBetButton>
  );
}
