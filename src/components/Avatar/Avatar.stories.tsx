import { ComponentMeta, Story } from "@storybook/react";
import Avatar, { AvatarProps } from ".";

export default {
  title: "Components/Avatar",
  component: Avatar,
  argTypes: {},
} as ComponentMeta<typeof Avatar>;

const Template: Story<AvatarProps> = (args) => <Avatar {...args} />;

export const Default = Template.bind({});

Default.args = {
  image:
    "https://beta.api.solanalysis.com/images/200x200/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/0.png",
};
