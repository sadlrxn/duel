import { ComponentMeta, Story } from "@storybook/react";
import JackpotGame, { JackpotGameProps } from "./JackpotGame";

export default {
  title: "Jackpot/JackpotGame",
  component: JackpotGame,
  argTypes: {},
} as ComponentMeta<typeof JackpotGame>;

const Template: Story<JackpotGameProps> = (args) => <JackpotGame {...args} />;

export const Default = Template.bind({});

Default.args = {};
