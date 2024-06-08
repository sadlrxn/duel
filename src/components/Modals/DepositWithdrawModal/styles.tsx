import styled from 'styled-components';
import { Tabs, TabPanel } from 'react-tabs';
import { Button } from 'components/Button';

import { Floating, SlideUp } from 'utils/animationToolkit';

export const StyledTabs = styled(Tabs)`
  display: flex;
  flex-direction: column;
  flex: 1;

  .react-tabs__tab-list {
    display: flex;
    margin: 0;
    z-index: 30;
    border: none;
    position: relative;
    overflow: auto hidden;
    margin-bottom: -2px;
    padding-bottom: 2px;
    max-width: 100vw;
    scrollbar-width: none;
    &::-webkit-scrollbar {
      display: none;
    }
    div {
      background: #202e3f;
      border-radius: 0px 0px 0 0;
      min-width: 900px;
      ${({ theme }) => theme.mediaQueries.lg} {
        border-radius: 17px 17px 0 0;
      }
    }
    .react-tabs__tab {
      bottom: 0;
      padding: 12px 20px;

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;
      line-height: 17px;
      text-align: center;
      letter-spacing: 0.18em;
      color: #697e9c;
      text-transform: uppercase;

      :focus:after {
        display: none !important;
      }

      b {
        z-index: 2;
        position: absolute;
        left: 50%;
        bottom: 0;
        ${({ theme }) => theme.mediaQueries.lg} {
          top: -2px;
        }
        transform: translateX(-50%);
        background-color: #4fff8b;
        width: 0%;
        height: 2px;
        border-radius: 2px;
        overflow: hidden;
        transition: all 0.3s ease-in;
      }

      :hover b {
        width: 75%;
      }
    }

    .react-tabs__tab--selected {
      background: #202e3f;
      border: 2px solid transparent;
      border-bottom: 0;
      color: #4fff8b;

      border-radius: 0px 0px 0 0;

      ${({ theme }) => theme.mediaQueries.md} {
        background: #132031;
      }
      ${({ theme }) => theme.mediaQueries.lg} {
        border-radius: 17px 17px 0 0;
        background: #132031;
        border-color: #4f617b;
      }

      ::before {
        position: absolute;
        background: transparent;
        ${({ theme }) => theme.mediaQueries.lg} {
          background: #132031;
        }
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
        width: 75%;
      }
    }
  }
`;

export const StyledTabPanel = styled(TabPanel)`
  &.react-tabs__tab-panel--selected {
    display: flex;
    flex-direction: column;
    flex: 1;
  }

  .container {
    position: relative;

    display: flex;
    flex-direction: column;
    flex: 1;
    background: linear-gradient(to bottom, #4f617b 1%, #1a293d00 95%);
    border-radius: 0 0 10px 10px;
    padding: 0px;
    ${({ theme }) => theme.mediaQueries.md} {
      padding: 2px;
    }

    .box {
      display: flex;
      flex-direction: column;
      flex: 1;
      padding: 30px;
      max-height: calc(100vh - 170px);
      overflow: auto;
      background: linear-gradient(180deg, #132031 0%, #1a293d 100%);
      border-radius: 0 0 10px 10px;
    }

    .chip-left {
      position: absolute;
      top: 45%;
      left: -4%;
      z-index: 10;

      ${({ theme }) => theme.mediaQueries.md} {
        top: 45%;
        left: 10%;
      }

      animation: ${SlideUp} 1s ease,
        ${Floating} 2s ease-in-out infinite alternate;
      animation-delay: 0s, 1s;
    }

    .chip-right {
      position: absolute;
      top: -5%;
      right: -10%;
      z-index: 10;

      ${({ theme }) => theme.mediaQueries.md} {
        top: 25%;
        right: 10%;
      }

      animation: ${SlideUp} 1s ease,
        ${Floating} 2s ease-in-out infinite alternate;
      animation-delay: 0s, 1s;
    }
  }
`;
