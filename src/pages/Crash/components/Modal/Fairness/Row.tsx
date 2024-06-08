import React from 'react';

import { Chip } from 'components';
import { PaidBalanceType } from 'api/types/chip';

import { StyledTd } from './styles';
import { convertBalanceToChip } from 'utils/balance';

interface RowProps {
  name?: string;
  betAmount?: number;
  profit?: number;
  multiplier?: string;
  paidBalanceType: PaidBalanceType;
}

const Row: React.FC<RowProps> = ({
  name,
  betAmount,
  profit,
  multiplier,
  paidBalanceType
}) => {
  return (
    <tr>
      <td>{name}</td>
      <StyledTd>
        <Chip
          chipType={paidBalanceType}
          price={convertBalanceToChip(betAmount!).toFixed(2)}
          fontWeight={400}
          color="white"
        />
      </StyledTd>
      <StyledTd>
        <Chip
          chipType={paidBalanceType}
          price={convertBalanceToChip(profit!).toFixed(2)}
          fontWeight={400}
          color="white"
        />
      </StyledTd>
      <td>{multiplier}</td>
    </tr>
  );
};

export default React.memo(Row);
