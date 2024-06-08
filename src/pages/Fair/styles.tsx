import styled from "styled-components";
import { Tabs, TabPanel } from "react-tabs";

export const StyledTabs = styled(Tabs)`
  .react-tabs__tab-list {
    display: flex;
    width: 800px;
    margin: 0;
    border: none;
    div {
      background: linear-gradient(
        5.33deg,
        #142131 95.74%,
        rgba(20, 33, 49, 0) 174.47%
      );
      border-radius: 17px 17px 0 0;
    }
    .react-tabs__tab {
      bottom: 0;
      padding: 25px 30px;

      .width_700 & {
        padding: 25px 40px;
      }

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
      background: #06080f;
      border: 2px solid #4f617b;
      border-radius: 17px 17px 0 0;
      border-bottom: 0;
      color: #4fff8b;

      ::before {
        position: absolute;
        background: #06080f;
        content: "";
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
  /* position: relative; */

  & > div {
    box-sizing: border-box;
    background: linear-gradient(180deg, #000000 -13.46%, #0b141e 16.48%);
    border: 2px solid transparent;
    border-image: linear-gradient(
        180deg,
        #4f617b 1.45%,
        rgba(26, 41, 61, 0) 49.08%
      )
      1;
    border-bottom: 0px;
    min-height: 300px;
    padding: 24px;

    .width_700 & {
      padding: 40px;
    }
  }
`;
