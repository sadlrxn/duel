import React from "react";

import { Chip } from "components";

import { StyledTd } from "./styles";

interface RowProps {
  name?: string;
  usdAmount?: number;
  nftAmount?: number;
  chance?: number;
}

const Row: React.FC<RowProps> = ({ name, usdAmount, nftAmount, chance }) => {
  return (
    <tr>
      <td>{name}</td>
      <StyledTd>
        <Chip price={usdAmount} fontWeight={400} color="white" />
      </StyledTd>
      <StyledTd>
        <Chip price={nftAmount} fontWeight={400} color="white" />
      </StyledTd>
      <td>{chance}%</td>
    </tr>
  );
};

export default React.memo(Row);
