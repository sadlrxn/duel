import { ComponentMeta, Story } from "@storybook/react";
import Players, { PlayersProps } from "./Players";

export default {
  title: "Jackpot/Players",
  component: Players,
  argTypes: {},
} as ComponentMeta<typeof Players>;

const Template: Story<PlayersProps> = (args) => <Players {...args} />;

export const Default = Template.bind({});

Default.args = {
  players: [],
};
