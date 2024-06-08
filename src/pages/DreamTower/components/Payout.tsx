import { Chip, Flex, Span } from 'components';
import { FC } from 'react';
import styled from 'styled-components';
import { convertBalanceToChip } from 'utils/balance';

const Payout: FC<{
  multiplier: number;
  profit: number;
  chipType: string;
}> = ({ multiplier, profit, chipType }) => {
  return (
    <StyledPayout>
      <Span
        color="#FFFFFF"
        fontWeight={700}
        textAlign="center"
        fontSize="30px"
        lineHeight="36.31px"
      >
        {multiplier.toFixed(2)}x
      </Span>
      <Flex justifyContent="space-around" width="100%">
        <Chip
          color="#FFFFFF"
          chipType={chipType}
          price={convertBalanceToChip(profit).toFixed(2)}
          fontWeight={600}
          fontSize="20px"
          lineHeight="24.2px"
        />
      </Flex>
    </StyledPayout>
  );
};

const StyledPayout = styled(Flex)`
  position: absolute;
  left: 20%;
  right: 20%;
  top: 25%;
  bottom: 64%;
  padding: 16px;
  flex-direction: column;
  align-items: center;
  z-index: 20;
  background: rgba(75, 39, 150, 0.55);
  border: 2px solid #7428c0;
  backdrop-filter: blur(3.5px);
  box-shadow: rgb(0 0 0 / 19%) 0px 10px 20px, rgb(0 0 0 / 23%) 0px 6px 6px;
  border-radius: 10px;
`;
export default Payout;
