import styled from 'styled-components';
import { Tabs, TabPanel } from 'react-tabs';
import { Box, Flex } from 'components';
import towerBg from 'assets/imgs/dreamtower/tower.png';

export const TabContainer = styled(Box)`
  width: 100%;
  .width_700 & {
    width: auto;
  }
`;

export const LevelButtonContainer = styled(Flex)`
  gap: 10px;
  margin-top: 10px;
  flex-direction: column;
  .width_700 & {
    flex-direction: row;
  }
`;

export const StyledTowerBox = styled(Box)<{
  isHoliday?: boolean;
  scale: number;
}>`
  width: 320px;
  height: 840px;
  background-image: url(${towerBg});
  background-size: cover;
  position: relative;
  ${({ isHoliday }) => !isHoliday && 'margin-bottom: 30px;'}

  transform-origin: center top;
  scale: ${({ scale }) => scale};

  .width_900 & {
    margin-top: ${({ isHoliday }) => (isHoliday ? '-30px' : '00px')};
  }
`;

export const StyledTabs = styled(Tabs)`
  .react-tabs__tab-list {
    display: flex;
    .width_700 & {
      width: 580px;
    }
    margin: 0;
    border: none;
    div {
      background: #202e3f;
      border-radius: 10px 10px 0 0;
    }
    .react-tabs__tab {
      bottom: 0;
      padding: 12px 15px;

      border-top-right-radius: 10px;

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 10px;

      ${({ theme }) => theme.mediaQueries.md} {
        font-size: 14px;
      }

      line-height: 10px;
      text-align: center;
      letter-spacing: 0.18em;
      color: #697e9c;
      text-transform: uppercase;

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
      background: #101b2c;

      border: 2px solid #4f617b;
      border-radius: 10px 10px 0 0;
      border-bottom: 0;
      color: #4fff8b;

      ::before {
        position: absolute;
        background: #101b2c;
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

export const StyledTabPanel = styled(TabPanel)`
  .container {
    position: relative;
    max-width: 580px;

    border: 2px solid #4f617b;
    border-radius: 0px 10px 10px 10px;

    padding: 16px;
    ${({ theme }) => theme.mediaQueries.md} {
      padding: 30px;
    }
    background: #132031bf;

    backdrop-filter: blur(5px);
  }
`;

export const StyledFlex = styled(Flex)`
  flex-direction: column;
  align-items: center;
  gap: 30px;
  margin-bottom: -60px;
  .width_900 & {
    align-items: start;
    flex-direction: row;
  }
`;
