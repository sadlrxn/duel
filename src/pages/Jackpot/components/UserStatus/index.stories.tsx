import { ComponentMeta, Story } from "@storybook/react";
import random from "lodash/random";
import UserStatus, { UserStatusProps } from "./index";

export default {
  title: "Jackpot/UserStatus2",
  component: UserStatus,
  argTypes: {},
} as ComponentMeta<typeof UserStatus>;

const Template: Story<UserStatusProps> = (args) => <UserStatus {...args} />;

export const Default = Template.bind({});

Default.args = {
  user: {
    id: 10,
    name: "Ninja",
    avatar:
      "https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png",
    level: 3,
    percent: 40,
  },
  nfts: [
    {
      mintAddress: "123123123123",
      image:
        "https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png",
      price: random(0, 10000),
    },
  ],
  amount: {
    usd: random(0, 1000),
    nft: random(0, 10000),
    total: random(0, 100000),
  },
  nftsToShow: 2,
};
