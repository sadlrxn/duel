import { Text } from 'components/Text';
import styled from 'styled-components';
import { Link } from 'react-router-dom';
import { Button } from 'components/Button';
import { Flex } from '..';

export const StyledHeader = styled.header`
  background: ${({ theme }) => theme.colors.primaryDark};
  ${({ theme }) => theme.mediaQueries.md} {
    background: linear-gradient(90deg, #1a2a3e 0%, #0b141e 100%);
  }
  height: 65px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  padding: 0px 12px;
  z-index: 200;
  ${({ theme }) => theme.mediaQueries.md} {
    padding: 0px 15px;
  }
`;

export const AppTitle = styled(Text)`
  flex-grow: 1;
  font-family: 'Inter';
  font-style: normal;
  font-weight: 600;
  font-size: 18.802px;
  ${({ theme }) => theme.mediaQueries.md} {
    font-size: 22.802px;
  }
  line-height: 28px;
  /* identical to box height */

  letter-spacing: 0.195em;
  text-transform: uppercase;

  color: #768bad;
`;

export const DepositBtn = styled(Button)`
  background: #1a5032;
  border-radius: 5px;
  border: 0;
  font-family: 'Inter';
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  /* identical to box height */

  padding: 10px 18px;
  gap: 8px;

  color: #4fff8b;
  cursor: pointer;
`;

export const WithdrawBtn = styled(Button)`
  background: #242f42;
  border-radius: 5px;
  border: 0;
  font-family: 'Inter';
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  color: #768bad;
  /* identical to box height */
  padding: 10px 12px;
  cursor: pointer;
  gap: 8px;
`;

export const RewardsBtn = styled(WithdrawBtn)`
  display: none;
  ${({ theme }) => theme.mediaQueries.xxl} {
    display: flex;
  }
`;

export const StyledAvatarContainer = styled(Flex)`
  align-items: center;
  margin: 0px 10px;

  div img {
    cursor: pointer;
  }

  ${({ theme }) => theme.mediaQueries.md} {
    margin: 0px 18px;
  }
`;

export const UserName = styled(Text)`
  font-family: 'Inter';
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  background-color: #242f42;
  padding: 10px 20px 10px 30px;
  border-radius: 5px;
  /* identical to box height */
  margin-left: -20px;
  /* Secondary Text */

  color: #768bad;

  display: none;

  ${({ theme }) => theme.mediaQueries.md} {
    display: block;
  }
`;

export const UserBalance = styled.div`
  position: relative;
  display: flex;
  height: 38px;
  align-items: center;
  background: linear-gradient(90deg, #503b00 0%, #2f2814 100%);
  ::before {
    position: absolute;
    background: #ffbe5c;
    content: '';
    display: block;
    height: 45%;
    left: 0;
    pointer-events: none;

    /* transform: translateY(-50%); */
    width: 2px;
  }

  border-radius: 5px;
  padding: 0px 12px;

  img {
    width: 14px;
    height: 14px;
    margin-right: 8px;
  }

  span {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 600;
    font-size: 16px;
    line-height: 19px;

    color: #fff6ca;
  }
  &:hover {
    cursor: pointer;
  }
`;

export const IconButton = styled(Button)`
  position: relative;
  background: #242f42;
  border-radius: 5px;
  border: 0;
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;

  cursor: pointer;

  .badge {
    position: absolute;
    top: -4px;
    right: -4px;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: #4fff8b;
    color: white;
  }
`;

export const ConnectWalletBtn = styled(Button)`
  background: transparent;
  border: 1.5px solid #4fff8b;
  border-radius: 29px;

  font-family: 'Inter';
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  /* identical to box height */

  color: #ffffff;

  padding: 10px 25px;
`;

export const ToggleBtn = styled(Button)`
  width: 46px;
  height: 46px;
  background: transparent;
  border-radius: 5px;
  transition: 0.5s;
  ${({ theme }) => theme.mediaQueries.md} {
    display: none;
  }
`;

export const ButtonContainer = styled(Flex)`
  flex-direction: row;
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
    gap: 8px;
  }
`;

export const StyledLink = styled(Link)`
  display: flex;
  align-items: center;
  margin-left: 10px;
`;
