import styled from 'styled-components';

import { Flex, Grid, BaseButton } from 'components';
import privateBackground from 'assets/imgs/icons/private.svg';

export const FairnessButton = styled(BaseButton)`
  display: flex;
  align-items: center;
  justify-content: center;

  width: 26px;
  height: 26px;
  padding-top: 4px;

  border-radius: 50%;

  background: ${({ theme }) => theme.coinflip.border};

  display: none;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    display: block;
  }

  path {
    transition: all 0.3s;
  }

  &:hover {
    path {
      fill: ${({ theme }) => theme.colors.success};
    }
  }
`;

export const CrownIcon = styled.img<{ side: string }>`
  position: absolute;
  width: 15px;
  height: 15px;
  top: 50%;
  transform: translateY(-50%);

  ${({ side }) => (side === 'duel' ? 'left: 0' : 'right: 0')};

  display: none;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    display: block;
  }
`;

export const PrivateGame = styled(Flex)`
  justify-content: center;
  align-items: center;
  gap: 7px;
  width: 210px;
  height: 22px;
  background: url(${privateBackground});
  color: ${({ theme }) => theme.coinflip.private};
  font-size: 12px;
  line-height: 15px;
`;

export const MiddleContainer = styled(Flex)`
  flex-direction: column;
  justify-content: center;
  height: 100%;
`;

export const RowContainer = styled(Grid)`
  /* position: relative; */
  grid-template-columns: 1fr max-content 1fr;
  width: 100%;
  /* overflow: hidden; */
  position: relative;

  &:before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    background: ${({ theme }) => theme.coinflip.border};
    z-index: -1;
  }
`;

export const Container = styled(Flex)`
  flex-direction: row;
  align-items: center;
  gap: 11px;

  border-bottom: 1px solid #0d141e;

  &:first-child > ${RowContainer} {
    border-top-left-radius: 16px;
    border-top-right-radius: 16px;
    &:before {
      border-top-left-radius: 16px;
      border-top-right-radius: 16px;
    }
  }

  &:last-child > ${RowContainer} {
    border-bottom-left-radius: 16px;
    border-bottom-right-radius: 16px;
    border-bottom: 0;
    &:before {
      border-bottom-left-radius: 16px;
      border-bottom-right-radius: 16px;
    }
  }
`;
