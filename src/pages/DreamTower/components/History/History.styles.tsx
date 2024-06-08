import { Box, Flex } from 'components';
import styled from 'styled-components';

export const HistoryItemContainer = styled(Flex)<{ selected: boolean }>`
  flex-direction: column;
  padding: 20px;
  gap: 9px;
  background: #182738;
  border-radius: 8px;
  min-width: 480px;

  cursor: pointer;
  ${({ selected }) =>
    selected &&
    `
    border: 1px solid #FFE87F;
    box-shadow: 0px 0px 16px rgba(255, 232, 127, 0.75);
  `}
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
