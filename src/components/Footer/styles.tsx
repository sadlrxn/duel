import styled from 'styled-components';
import { Flex } from 'components/Box';
import { Button } from 'components/Button';

export const FlexFooter = styled(Flex)`
  background: linear-gradient(90.08deg, #0d141d 0.07%, #141e30 99.94%);
  border-top: 1px solid #242f42;
  justify-content: space-around;
  align-items: center;
  z-index: 1;

  padding: 46px 15px;
  flex-direction: column;
  gap: 24px;
  ${({ theme }) => theme.mediaQueries.md} {
    flex-direction: row;
  }
`;

export const PingBtn = styled(Button)`
  background: rgba(255, 236, 139, 0.2);
  border-radius: 5px;

  font-style: normal;
  font-weight: 400;
  font-size: 12px;
  line-height: 15px;

  padding: 3px 7px;
  color: #ffec8b;
`;
