import { Link, NavLink } from 'react-router-dom';
import styled, { css } from 'styled-components';
import { Sidebar, Menu, MenuItem, SubMenu } from 'react-pro-sidebar';
import { Flex } from 'components/Box';
import { Text } from 'components/Text';
import { Button } from 'components/Button';

export const StyledSidebar = styled(Sidebar)`
  position: fixed !important;
  top: 0px;
  padding-top: 65px;
  height: 100%;
  border: none !important;

  color: #96a8c2;
  z-index: 100;

  &.ps-collapsed .ps-submenu-expand-icon {
    display: none;
  }

  .ps-sidebar-container {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .ps-sidebar-container {
    background: linear-gradient(0deg, #101c2b, #101c2b),
      linear-gradient(180.6deg, #111923 0.46%, #132133 99.43%),
      linear-gradient(180deg, #000306 0%, #010912 60.19%);
  }

  .ps-menuitem-root .ps-menu-button {
    height: 54px;
    background: transparent;
    margin: 3.5px 0px 3.5px 6px;
    border-radius: 10px 0 0 10px;
    font-family: 'Inter';
    font-weight: 600;
    font-size: 14px;
    line-height: 17px;
    padding: 12px 8px 12px 16px;

    &:hover {
      color: #bdcde3;
      background: linear-gradient(
        90deg,
        rgba(150, 168, 194, 0.2) 0%,
        rgba(150, 168, 194, 0.05) 100%
      );
      background-color: inherit;
    }
    &.active,
    &.ps-active {
      background: linear-gradient(
        90deg,
        rgba(79, 255, 139, 0.2) 0%,
        rgba(79, 255, 139, 0.05) 100%
      );

      border-right: 3px solid #4fff8b;
      color: #4fff8b;
    }

    .ps-menu-icon {
      width: 30px;
      min-width: 30px;
      height: 30px;
      margin-right: 14px;
    }
  }

  .ps-menuitem-root.submenu-grand .ps-menu-button {
    color: #c8ab11;
    &:hover {
      color: #ffd912;
      background: linear-gradient(
        90deg,
        rgba(200, 171, 17, 0.2) 0%,
        rgba(200, 171, 17, 0.05) 100%
      );
      background-color: inherit;
    }
    &.active {
      background: linear-gradient(
        90deg,
        rgba(200, 171, 17, 0.2) 0%,
        rgba(200, 171, 17, 0.05) 100%
      );

      border-right: 3px solid #ffd912;
      color: #ffd912;
    }
  }

  .ps-menuitem-root .ps-submenu-content {
    background-color: inherit;
  }

  .ps-submenu-content .ps-menu-button {
    font-size: 13px;
    font-weight: 500;
    height: 30px;

    margin-left: 50px;
    &:hover {
      color: #bdcde3;
    }
    &.active {
      background: transparent;
      border-right: none;
      color: #4fff8b;
      font-weight: 600;
    }

    .ps-menu-label {
    }
  }
`;
