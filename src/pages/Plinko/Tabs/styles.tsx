import styled from 'styled-components';

import { Flex, Button, BaseButton } from 'components';
import { InputBox } from 'components/InputBox';

export const BetButton = styled(Button)`
  font-size: 16px;
  font-weight: 600;
  color: black;
  width: 100%;
  height: 52px;
  margin-top: auto;
  border-radius: 5px;
`;

export const StyledInputBox = styled(InputBox)`
  input {
    color: white;
  }
`;

export const UpperBtn = styled(Button)`
  padding: 7.2px 5.5px;
  border-radius: 7px 7px 0px 0px;
  background: #24354d;
`;

export const DownBtn = styled(Button)`
  padding: 7.2px 5.5px;
  border-radius: 0px 0px 7px 7px;
  background: #24354d;
  svg {
    transform: rotate(180deg);
  }
`;

export const BetAllButton = styled(BaseButton)`
  background-color: #24354c;
  color: #526d90;
  font-size: 13px;
  font-weight: 500;

  border-radius: 37px;
  transition: all 0.3s ease-in;

  &:hover {
    background-color: ${({ theme }) => theme.coinflip.greenDark};
    color: ${({ theme }) => theme.colors.success};
  }
`;

export const HistoryListContainer = styled(Flex)`
  flex-direction: column;
  background: #0f1a26;
  border-radius: 13px;
  padding: 15px 5px 15px 0px;
  gap: 9px;
`;

export const HistoryList = styled(Flex)`
  flex-direction: column;
  max-height: 400px;
  gap: 15px;
  padding: 0px 5px 0px 15px;

  overflow: auto;
  scrollbar-width: 5px;
`;
