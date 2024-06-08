import React, { useState } from "react";
import { ComponentMeta } from "@storybook/react";
import Tab, { TabContent } from ".";

export default {
  title: "Components/Tab",
  component: Tab,
  argTypes: {},
} as ComponentMeta<typeof Tab>;

export const Default = () => {
  const [activeTab, setActiveTab] = useState("LOW");

  return (
    <Tab
      labels={["LOW", "HIGH"]}
      activeTab={activeTab}
      setActiveTab={setActiveTab}
    >
      <TabContent label="LOW">
        a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />
      </TabContent>
      <TabContent label="HIGH">b</TabContent>
    </Tab>
  );
};
