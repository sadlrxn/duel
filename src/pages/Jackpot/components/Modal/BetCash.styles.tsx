import styled from 'styled-components';

import { Box, Flex, Grid, Text, Button } from 'components';
import { Input } from 'components/Input';

export const MaxText = styled(Text)`
  font-size: 14px;
  padding: 6px 18px;
  color: ${({ theme }) => theme.coinflip.private};
`;

export const DepositButton = styled(Button)`
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: black;
  width: 100%;
  border-radius: 5px;
  min-height: 52px;
  margin-top: auto;
`;

export const StyledInput = styled(Input)`
  font-size: 16px;
  /* color: ${({ theme }) => theme.coinflip.private}; */
  color: white;
`;

export const InputContainer = styled(Grid)`
  gap: 12px;
  margin-top: 12px;
  padding: 10px 14px 10px 25px;
  grid-auto-flow: column;
  grid-template-columns: max-content auto max-content;
  align-items: center;
  border-radius: 11px;
  font-weight: 400;
  background: ${({ theme }) => theme.jackpot.input};
`;

export const HorizontalDivider = styled(Box)`
  width: 100%;
  height: 1px;
  background: black;
  margin-top: 28px;
  margin-bottom: 24px;
`;

export const Container = styled(Flex)`
  flex-direction: column;
  background-color: ${({ theme }) => theme.jackpot.modal};
  color: white;
  padding: 40px 27px;
  overflow: auto;

  width: 100vw;
  height: calc(100vh - 65px);
  border-radius: 0px;

  & :nth-child(4) {
    order: 5;
  }

  & :nth-child(5) {
    order: 4;
    margin-top: 0px;
  }

  ${({ theme }) => theme.mediaQueries.sm} {
    border-radius: 10px;
    width: auto;
    height: auto;

    & :nth-child(4) {
      order: 4;
    }

    & :nth-child(5) {
      order: 5;
      margin-top: 32px;
    }
  }
`;
