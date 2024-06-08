import styled from "styled-components";
import Collapse from "rc-collapse";
import { Tabs, TabPanel } from "react-tabs";
import { Top, Floating, SlideUp } from "utils/animationToolkit";

export const StyledTabs = styled(Tabs)`
  .react-tabs__tab-list {
    display: flex;
    flex-direction: column;
    ${({ theme }) => theme.mediaQueries.md} {
      flex-direction: row;
    }
    gap: 35px;
    justify-content: space-evenly;
    background: #030609;
    border-radius: 30px;
    border: none;

    padding: 24px 31px;
    margin: 0 0 60px;
    .react-tabs__tab {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 100%;
      padding: 43px;
      background: linear-gradient(27.16deg, #121823 11.03%, #070b10 89.37%);
      border-radius: 22px;
      color: #8591a3;
      :focus:after {
        display: none !important;
      }
    }

    .react-tabs__tab--selected {
      position: relative;
      background: linear-gradient(
        0deg,
        rgba(79, 255, 139, 0.2) 0%,
        rgba(79, 255, 139, 0) 132.27%
      );
      border: none !important;

      animation: ${Floating} 0.8s ease-in-out infinite alternate;

      svg {
        color: #4fff8b;
      }

      ::before {
        content: "";
        position: absolute;
        inset: 0;
        border-radius: 22px;
        padding: 2px;
        background: linear-gradient(
          to bottom,
          rgba(79, 255, 139, 1) 10%,
          rgba(79, 255, 139, 0) 90%
        );
        mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        -webkit-mask: linear-gradient(#fff 0 0) content-box,
          linear-gradient(#fff 0 0);
        mask-composite: xor;
        -webkit-mask-composite: xor;
        /* mask-composite: exclude; */
      }
    }
  }
`;

export const StyledTabPanel = styled(TabPanel)`
  animation: ${Top} 0.6s;
`;

export const StyledCollapse = styled(Collapse)`
  background-color: transparent !important;
  border: none !important;
  animation: ${Top} 0.6s;

  .rc-collapse-content {
    background-color: transparent !important;
  }

  .rc-collapse-item {
    border: none !important;

    .rc-collapse-header {
      color: white !important;
      justify-content: space-between;
      padding: 21px 16px !important;

      .rc-collapse-expand-icon {
        order: 2;
      }

      .rc-collapse-header-text {
        font-size: 20px;
      }
    }

    .rc-collapse-content .rc-collapse-content-box {
      animation: ${SlideUp} 0.6s;
      font-size: 16px;
      color: #8591a3;
    }
  }

  .rc-collapse-item-active .rc-collapse-header {
    .rc-collapse-expand-icon {
      transition: 0.5s;
      transform: rotate(180deg);

      svg path {
        stroke: #4fff8b;
      }
    }

    .rc-collapse-header-text {
      color: #4fff8b !important;
    }
  }
`;

export const StyledLink = styled.a`
  color: #44aae6;
  &:link {
    text-decoration: none;
  }
  &:visited {
    text-decoration: none;
  }
  &:hover {
    color: #4fff8b;
    text-decoration: underline;
  }
  &:active {
    text-decoration: none;
  }
`;
