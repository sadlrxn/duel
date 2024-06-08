import styled from 'styled-components';

import { Box, Flex, Span, Button, BaseButton } from 'components';

export const Side = styled(Flex)<{ active: boolean }>`
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;

  &:hover {
    ${({ active }) => (active ? '' : 'opacity: 0.4;')}
  }

  opacity: ${({ active }) => (active ? '1' : '0.3')};

  transition: all 0.3s ease-in;
`;

export const Sides = styled(Box)`
  display: flex;
  gap: 40px;
  justify-content: space-between;
  align-items: center;
`;

export const Title = styled(Span)`
  font-size: 14px;
  color: ${({ theme }) => theme.coinflip.title};

  ${({ theme }) => theme.mediaQueries.md} {
    font-size: 16px;
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

export const SideContainer = styled(Box)`
  display: grid;
  justify-content: center;
  align-items: center;

  .width_1100 & {
    display: block;
    padding: 0px;
    grid-template-columns: 40% 60%;
    gap: 30px;
    justify-content: space-between;
  }
`;

export const PrivateButton = styled(Button)`
  width: 32px;
  height: 32px;
  & > img {
    transform: scale(0.6);
  }

  .width_1100 & {
    width: 52px;
    height: 52px;
    & > img {
      transform: scale(1);
    }
  }
`;

export const CreateButton = styled(Button)`
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: black;
  width: 230px;
  height: 48px;
  border-radius: 12px;
`;

export const DataContainer = styled(Box)`
  display: grid;
  gap: 24px;
  align-items: end;
  padding: 16px 20px;
  background-color: #1e3248;
  grid-template-columns: max-content 1fr;

  ${CreateButton} {
    justify-self: end;
  }

  .width_1100 & {
    background-color: transparent;
    padding: 0;
    gap: 40px;
  }
`;

export const Container = styled(Box)`
  display: flex;
  flex-direction: column;
  background: ${({ theme }) => theme.coinflip.gradients.background};
  border: 2px solid ${({ theme }) => theme.coinflip.border};
  border-radius: 16px;
  overflow: hidden;
  padding: 15px 20px;
  text-align: center;
  gap: 30px;

  .width_1100 & {
    display: flex;
    flex-direction: row;
    text-align: inherit;
    justify-content: space-between;
    align-items: start;
    padding: 30px 40px;
  }
`;

export const FlexBox = styled(Flex)`
  justify-content: center;

  .width_1100 & {
    justify-content: space-between;
  }
`;
