import React from "react";
import styled, { css } from "styled-components";

interface TabHeaderButtonProps {
  label: string;
  activeTab: string;
  left?: boolean;
  right?: boolean;
  onClick: () => void;
}

const TabHeaderTopBorder = styled.div`
  position: absolute;
  left: 50%;
  top: -2px;
  transform: translateX(-50%);
  background-color: ${({ theme }) => theme.colors.success};
  width: 0%;
  height: 2px;
  border-radius: 2px;
  overflow: hidden;
  transition: all 0.3s ease-in;
`;

interface TabHeaderButtonWrapperProps {
  active: boolean;
  left: boolean;
  right: boolean;
}

const TabHeaderButtonWrapper = styled.li<TabHeaderButtonWrapperProps>`
  position: relative;
  display: inline-block;
  list-style: none;
  padding: 0.5rem 2rem;
  margin-right: -30px;
  cursor: pointer;
  box-sizing: border-box;
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
  font-family: "Inter";
  font-weight: 600;
  font-size: 14px;
  color: #697e9c;
  border-color: #142131;
  border-width: 2px;
  border-style: solid;
  border-bottom: 0;
  ${({ active, theme }) =>
    active
      ? css`
          background-color: #05090d;
          color: ${theme.colors.success};
          border-color: #7389a9;
          border-left-style: solid;
          border-right-style: solid;
          border-top-style: solid;
          border-width: 2px;
          z-index: 1;
        `
      : css`
          background-color: #142131;
        `}

  &::before,
  &::after {
    position: absolute;
    bottom: 0px;
    width: 10px;
    height: 10px;
    border: 2px solid #7389a9;
  }

  &::before {
    ${({ active, left }) => (active && !left ? `content: "";` : "")}
    left: -12px;
    border-bottom-right-radius: 6px;
    border-width: 0px 2px 2px 0px;
    box-shadow: 5px 0px 0 #05090d;
  }

  &:after {
    ${({ active, right }) => (active && !right ? `content: "";` : "")}
    right: -12px;
    border-bottom-left-radius: 6px;
    border-width: 0px 0px 2px 2px;
    box-shadow: -5px 0px 0 #05090d;
  }

  &:hover {
    ${TabHeaderTopBorder} {
      width: 60%;
    }
  }
`;

const StyledTabHeaderButton = styled.div`
  padding: 0.5rem 2rem;
  letter-spacing: 0.18em;
`;

export default function TabHeaderButton({
  label,
  activeTab,
  left = false,
  right = false,
  onClick,
}: TabHeaderButtonProps) {
  const active = label === activeTab;
  return (
    <TabHeaderButtonWrapper
      active={active}
      onClick={onClick}
      left={left}
      right={right}
    >
      <TabHeaderTopBorder />
      <StyledTabHeaderButton>{label}</StyledTabHeaderButton>
    </TabHeaderButtonWrapper>
  );
}
