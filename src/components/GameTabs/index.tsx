import React from 'react';
import styled, { css } from 'styled-components';
import { TabList, Tab, Tabs, TabPanel } from 'react-tabs';

import { Box, BoxProps, Flex } from 'components/Box';

export interface GameTabsProps extends BoxProps {
  onSelect: any;
  tabIndex: number;
  disabled?: boolean;
  labels?: string[];
  tabs?: React.ReactNode[];
  tabwidth?: string;
  tabSelectedBackground?: string;
  tabPanelWidth?: string;
  tabPanelBackground?: string;
  tabPanelMaxHeight?: string;
}

export default function GameTabs({
  onSelect,
  tabIndex,
  disabled = false,
  labels = [],
  tabs = [],
  tabwidth,
  tabSelectedBackground,
  tabPanelWidth,
  tabPanelBackground,
  tabPanelMaxHeight,
  ...props
}: GameTabsProps) {
  return (
    <Box {...props}>
      <StyledTabs
        tabwidth={tabwidth}
        onSelect={onSelect}
        selectedIndex={tabIndex}
        tabselectedbackground={tabSelectedBackground}
      >
        <TabList>
          <Flex>
            {labels.map(label => {
              return (
                <Tab disabled={disabled} key={label}>
                  {label} <b />
                </Tab>
              );
            })}
          </Flex>
        </TabList>
        {tabs.map((tab, index) => {
          return (
            <StyledTabPanel
              key={labels[index]}
              width={tabPanelWidth}
              background={tabPanelBackground}
              maxheight={tabPanelMaxHeight}
            >
              <div className="container">{tab}</div>
            </StyledTabPanel>
          );
        })}
      </StyledTabs>
    </Box>
  );
}

export const StyledTabs = styled(Tabs)<{
  tabwidth?: string;
  tabselectedbackground?: string;
}>`
  .react-tabs__tab-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, 1fr);
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

      ${({ tabwidth }) =>
        tabwidth &&
        css`
          width: ${tabwidth};
        `}

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;

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

      ${({ tabselectedbackground }) =>
        tabselectedbackground &&
        css`
          background: ${tabselectedbackground};
        `}

      border: 2px solid #4f617b;
      border-radius: 10px 10px 0 0;
      border-bottom: 0;
      color: #4fff8b;

      ::before {
        position: absolute;
        background: #101b2c;
        ${({ tabselectedbackground }) =>
          tabselectedbackground &&
          css`
            background: ${tabselectedbackground};
          `}
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

export const StyledTabPanel = styled(TabPanel)<{
  width?: string;
  background?: string;
  maxheight?: string;
}>`
  .container {
    position: relative;
    /* max-width: 580px; */

    border: 2px solid #4f617b;
    border-radius: 0px 0px 10px 10px;

    padding: 10px 14px;
    ${({ theme }) => theme.mediaQueries.md} {
      padding: 18px 20px;
    }
    background: #132031bf;

    backdrop-filter: blur(5px);

    ${({ width }) =>
      width &&
      css`
        width: ${width};
      `}

    ${({ maxheight }) =>
      maxheight &&
      css`
        max-height: ${maxheight};
      `}

    ${({ background }) =>
      background &&
      css`
        background: ${background};
      `}
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
