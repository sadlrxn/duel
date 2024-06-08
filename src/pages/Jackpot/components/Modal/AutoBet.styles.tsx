import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import { Flex, Button, Notification } from 'components';

export const StyledNotification = styled(Notification)`
  display: block;
  position: relative;
  transform: none;
  top: auto;
  right: auto;
  margin-left: -15px;
  margin-top: -35px;
  min-width: max-content;
`;

export const StackedImage = styled(LazyLoadImage)`
  & + & {
    margin-left: -20px;
  }

  width: 45px;
  height: 45px;

  min-width: 45px;
  min-height: 45px;

  border: 2.48538px solid #0f1a26;
  border-radius: 12.2807px;
`;

export const NftAutoBetButton = styled(Button)`
  width: 100%;
  min-height: 52px;

  color: ${({ theme }) => theme.colors.text};
  font-weight: 600;
  font-size: 14px;
  line-height: 18px;
`;
NftAutoBetButton.defaultProps = {
  variant: 'secondary'
};

export const Container = styled(Flex)`
  flex-direction: column;
  background-color: ${({ theme }) => theme.jackpot.modal};
  color: white;
  padding: 40px 27px;
  overflow: auto;

  width: 100vw;
  height: calc(100vh - 65px);
  border-radius: 0px;

  ${({ theme }) => theme.mediaQueries.sm} {
    border-radius: 10px;
    max-width: 470px;
    width: auto;
    height: auto;
  }
`;
