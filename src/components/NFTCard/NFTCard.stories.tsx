import React from "react";
import { ComponentMeta, Story } from "@storybook/react";

import NFTCard, { NFTCardProps } from ".";

export default {
  title: "Components/NFTCards/Default",
  component: NFTCard,
  argTypes: {},
} as ComponentMeta<typeof NFTCard>;

const Template: Story<NFTCardProps> = (args) => <NFTCard {...args} />;

export const Default = Template.bind({});

Default.args = {
  price: 1827,
  collectionName: "Degods",
  name: "DeGods #2401",
  image:
    "https://beta.api.solanalysis.com/images/200x200/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/0.png",
};
