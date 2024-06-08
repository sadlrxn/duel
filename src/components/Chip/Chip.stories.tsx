import React from "react";
import { ComponentMeta, Story } from "@storybook/react";

import Chip from "./Chip";

export default {
  title: "Coinflip/Chip",
  component: Chip,
  argTypes: {},
} as ComponentMeta<typeof Chip>;

const Template: Story<any> = (args) => <Chip {...args} />;

export const Default = Template.bind({});

Default.args = {
  price: 182731983,
  fontWeight: 500,
  fontSize: "14px",
  color: "chip",
};
