import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import { Box, Badge, Text, Notification, Flex } from 'components';

export const Divider = styled.div`
  background-color: #3f4f63;
  width: 1px;
  height: 100%;

  .user_status_max_640 & {
    display: none;
  }
`;

export const DividerHorizontal = styled.div`
  background-color: #3f4f63;
  height: 1px;
  margin-top: 5px;
  width: 100%;

  display: none;
  .user_status_max_640 & {
    display: block;
  }
`;

export const UsdAmount = styled(Badge)`
  height: max-content;
  background: rgba(255, 186, 48, 0.1);
`;
UsdAmount.defaultProps = { variant: 'secondary' };

export const NftAmount = styled(Badge)`
  height: max-content;
`;

export const UserInfo = styled(Box)`
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  gap: 5px;
  color: ${({ theme }) => theme.colors.textWhite};
  font-weight: 500;
  font-size: 14px;
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

export const User = styled(Box)`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const UserContainer = styled(Box)`
  display: flex;
  justify-content: space-between;
  align-items: center;

  ${UsdAmount} {
    display: none;
    .user_status_max_640 & {
      display: block;
    }
  }
`;

export const StyledText = styled(Text)`
  border-radius: 999px;
  padding: 1px 0.5rem;
  width: max-content;
  height: max-content;
  color: #0b141e;
  font-size: 12px;
  background-color: ${({ theme }) => theme.colors.success};
  font-weight: 600;
`;

export const StyledIntro = styled(Text)`
  font-size: 12px;
  font-weight: 500;
  color: #768bad;

  display: none;
  .user_status_max_640 & {
    display: block;
  }
`;

export const Nft = styled(Box)`
  display: grid;
  grid-auto-flow: column;
  align-items: center;
  gap: 10px;

  .user_status_max_640 & {
    /* width: 100%; */
    grid-template-columns: max-content auto max-content;
  }
`;

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

export const NftContainer = styled(Box)`
  /* display: flex;
  justify-content: flex-start; */

  display: grid;
  grid-template-columns: max-content max-content max-content 1fr;

  align-items: center;

  .user_status_max_640 &.hide {
    display: none;
  }

  .user_status_max_640 & {
    grid-template-columns: 1fr;
    ${StyledText} {
      display: none;
    }
  }

  ${NftAmount} {
    display: none;
    .user_status_max_640 & {
      display: block;
    }
  }
`;

export const AmountContainer = styled(Box)`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const Container = styled(Box)`
  display: grid;
  grid-template-columns: 155px max-content auto max-content max-content;
  gap: 30px;
  padding: 15px;

  &.user_status_max_640 {
    grid-template-columns: auto;
    gap: 10px;

    ${AmountContainer} {
      display: none;
    }
  }

  border-radius: 10px;
`;

export const FlexBox = styled(Flex)`
  display: block;
  .user_status_max_640 & {
    display: none;
  }
`;
