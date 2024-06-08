import React from "react";
import { ComponentMeta, Story } from "@storybook/react";

import CreateGame from "./CreateGame";

export default {
  title: "Coinflip/CreateGame",
  component: CreateGame,
  argTypes: {},
} as ComponentMeta<typeof CreateGame>;

const Template: Story<any> = (args) => <CreateGame {...args} />;

export const Default = Template.bind({});

Default.args = {};
