import { ComponentMeta, Story } from "@storybook/react";
import Badge from "./Badge";
import { BadgeProps } from "./types";
export default {
  title: "Components/Badge",
  component: Badge,
  argTypes: {},
} as ComponentMeta<typeof Badge>;

const Template: Story<BadgeProps> = (args) => <Badge {...args} />;

export const Default = Template.bind({});

Default.args = {
  children: <>0.95%</>,
  variant: "primary",
};
