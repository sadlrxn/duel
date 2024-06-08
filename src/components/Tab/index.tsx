import React from "react";
import styled from "styled-components";
import upperCase from "lodash/upperCase";
import TabHeaderButton from "./TabHeaderButton";

interface TabProps {
  children: React.ReactElement[];
  labels: string[];
  activeTab: string;
  setActiveTab: React.Dispatch<React.SetStateAction<string>>;
}

const TabHeader = styled.ul`
  margin-top: 0px;
  margin-bottom: 0px;
  padding-left: 0px;
  &::after {
    clear: both;
    content: "";
    display: table;
  }
`;

const TabContentWrapper = styled.div`
  background: linear-gradient(180deg, #05090d 0%, #0b141e 100%);

  margin-top: -2px;
  border: solid #7389a9;
  border-width: 2px 2px 0 2px;
  border-top-right-radius: 5px;
  color: #7389a9;
`;

export const TabContent = styled.div<{ label: string }>``;

export default function Tab({
  children,
  labels,
  activeTab,
  setActiveTab,
  ...props
}: TabProps) {
  return (
    <div className="tabs">
      <TabHeader>
        {children.map((_, index) => {
          const label = upperCase(labels[index]);
          return (
            <TabHeaderButton
              activeTab={upperCase(activeTab)}
              key={label}
              label={label}
              left={index === 0}
              onClick={() => setActiveTab(labels[index])}
            />
          );
        })}
      </TabHeader>
      <TabContentWrapper {...props}>
        {children.map((child) => {
          if (child.props.label !== activeTab) return undefined;
          return child.props.children;
        })}
      </TabContentWrapper>
    </div>
  );
}
