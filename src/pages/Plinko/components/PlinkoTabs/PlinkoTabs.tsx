import React from 'react';
import styled from 'styled-components';
import { Tabs, TabPanel, TabList, Tab } from 'react-tabs';

import { Flex } from 'components';
import { AutoBet, Manual } from '../../Tabs';

export default function PlinkoTabs() {
  return (
    <StyledTabs>
      <TabList>
        <Flex>
          <Tab>
            MANUAL
            <b />
          </Tab>
          <Tab>
            AUTOMATIC
            <b />
          </Tab>
        </Flex>
      </TabList>
      <StyledTabPanel>
        <div className="container">
          <Manual />
        </div>
      </StyledTabPanel>
      <StyledTabPanel>
        <div className="container">
          <AutoBet />
        </div>
      </StyledTabPanel>
    </StyledTabs>
  );
}

const StyledTabs = styled(Tabs)`
  font-size: 12px;

  .react-tabs__tab-list {
    display: flex;
    border: none;
    width: 100%;
    margin: 0px;

    div {
      width: 100%;
      background: #202e3f;
      border-radius: 10px 10px 0 0;
    }
    .react-tabs__tab {
      bottom: 0;
      padding: 12px 15px;
      width: 100%;

      border-top-right-radius: 10px;

      font-family: 'Inter';
      font-style: normal;
      font-weight: 600;
      font-size: 14px;

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

      border: 0px solid #4f617b;
      border-bottom-width: 2px;
    }

    .react-tabs__tab--selected {
      background: #111e2e;
      backdrop-filter: blur(5px);

      border-width: 2px;
      border-radius: 10px 10px 0 0;
      border-bottom: 0;
      color: #4fff8b;

      b {
        width: 60%;
      }
    }
  }
`;

const StyledTabPanel = styled(TabPanel)`
  .container {
    position: relative;
    width: 100%;

    padding: 20px 23px 20px 20px;

    background: linear-gradient(
      180deg,
      rgba(19, 32, 49, 0.75) -0.73%,
      rgba(26, 41, 60, 0.75) 72.81%
    );
    border: 2px solid #4f617b;
    border-top-width: 0px;
    border-radius: 0px 0px 10px 10px;

    backdrop-filter: blur(5px);
  }
`;
