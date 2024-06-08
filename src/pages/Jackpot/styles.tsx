import styled, { css } from 'styled-components';
import { Tabs } from 'react-tabs';

import { Box, Flex, Grid, Button } from 'components';

export const HistoryButton = styled(Button)<{
  prev?: boolean;
}>`
  min-width: max-content;
  height: 28px;
  color: ${({ theme }) => theme.colors.text};
  padding: 5px 7px;

  ${({ prev }) =>
    prev
      ? css`
          &::after {
            content: 'Prev Game';
          }
        `
      : css`
          &::before {
            content: 'Next Game';
          }
        `}
`;
HistoryButton.defaultProps = {
  variant: 'secondary',
  prev: false
};

export const HistoryButtonContainer = styled(Flex)`
  display: none;
  .width_700 & {
    display: flex;
  }
`;

export const TotalChip = styled(Flex)`
  display: flex;
  align-items: center;
  padding: 5px 7px 4px;
  font-weight: 600;
  border-radius: 4.5px;
  border: 1.04px solid #c0b264;

  color: #fff6ca;
  font-size: 10px;
  font-weight: 600;
  line-height: 13px;
  gap: 6px;

  position: absolute;
  top: 0px;
  left: 50%;
  transform: translate(-50%, calc(-50% - 2px));
  overflow: hidden;

  &::before {
    position: absolute;
    left: 0;
    top: 0;
    content: '';
    z-index: -1;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      360deg,
      rgba(53, 54, 31, 0.5) 0%,
      rgba(11, 15, 19, 0.5) 100%
    );
    backdrop-filter: blur(5px);
  }
`;

export const TimeLine = styled(Box)`
  position: absolute;
  bottom: 0px;
  left: 50%;
  transform: translate(-50%, 0);
  width: 100%;
  height: 2px;
  overflow: hidden;
`;

export const ChipIconContainer = styled(Flex)``;

export const BetButtonContainer = styled(Grid)<{ roundId?: number }>`
  width: 100%;
  align-items: center;

  grid-template-columns: 1fr;
  gap: 20px;

  margin-top: 10px;
  margin-bottom: 25px;

  & > button {
    order: 3;
  }

  & > div:nth-child(1) {
    ${({ roundId }) =>
      roundId
        ? css`
            flex-direction: column;
            gap: 12px;
          `
        : ''}
  }

  .width_900 & {
    margin-top: 20px;
    grid-template-columns: 153px 1fr 153px;

    & > div:nth-child(1) {
      order: 2;
      ${({ roundId }) =>
        roundId
          ? css`
              flex-direction: column;
              gap: 12px;
            `
          : ''}
    }
    & > div:nth-child(2) {
      order: 1;
      flex-direction: column;
      gap: 12px;
    }
  }
`;

export const Container = styled(Box)`
  display: flex;
  flex-direction: column-reverse;
  gap: 25px;
  width: 100%;

  .width_1100 & {
    display: grid;
    grid-template-columns: 310px 1fr;
  }
`;

export const StyledTabs = styled(Tabs)<{ isgameend: boolean }>`
  .react-tabs__tab-list {
    display: flex;
    /* width: 800px; */
    margin: 0;
    border: none;
    & > div {
      z-index: 1;
      margin-bottom: 2px;
      background: linear-gradient(
        5.33deg,
        #142131 95.74%,
        rgba(20, 33, 49, 0) 174.47%
      );
      border-radius: 17px 17px 0 0;
    }
    .react-tabs__tab {
      position: relative;
      bottom: 0;
      padding: 11px 20px;
      .width_700 & {
        padding: 11px 40px;
      }

      p {
        padding: 0;
        margin: 0;
      }

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;
      line-height: 17px;
      text-align: center;
      letter-spacing: 0.18em;
      color: #697e9c;
      text-transform: uppercase;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;

      :focus:after {
        display: none !important;
      }

      b {
        position: absolute;
        left: 50%;
        top: -2px;
        transform: translateX(-50%);
        background-color: #4fff8b;
        width: 0%;
        height: 2px;
        border-radius: 2px;
        overflow: hidden;
        transition: all 0.3s ease-in;
      }

      :hover b {
        width: 60%;
      }
    }

    .react-tabs__tab--selected {
      background: ${({ isgameend }) => (isgameend ? '#0b141e' : '#050a0e')};
      border: 2px solid #7389a9;
      border-radius: 17px 17px 0 0;
      border-bottom: 0;
      color: #4fff8b;

      ${TimeLine} {
        display: none;
      }

      ${TotalChip} {
        display: none;
      }

      ::before {
        position: absolute;
        background: ${({ isgameend }) => (isgameend ? '#0b141e' : '#050a0e')};
        content: '';
        display: block;
        height: 2px;
        width: 100%;
        left: 0;
        bottom: -2px;
        pointer-events: none;
        z-index: 10;
      }

      b {
        width: 60%;
      }
    }
  }
`;
