import React from "react";
import { ComponentMeta, Story } from "@storybook/react";

import Coin from "./Coin";

export default {
  title: "Coinflip/Coin",
  component: Coin,
  argTypes: {},
} as ComponentMeta<typeof Coin>;

const Template: Story<any> = (args) => <Coin {...args} />;

export const Default = Template.bind({});

Default.args = {
  side: "duel",
  size: 72,
  active: true,
};
