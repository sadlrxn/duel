import { ComponentMeta, Story } from "@storybook/react";
import NFTCardWithPrice, { NFTCardWithPriceProps } from ".";

export default {
  title: "Components/NFTCards/SimpleNFTCard",
  component: NFTCardWithPrice,
  argTypes: {},
} as ComponentMeta<typeof NFTCardWithPrice>;

const Template: Story<NFTCardWithPriceProps> = (args) => (
  <NFTCardWithPrice {...args} />
);

export const Default = Template.bind({});

Default.args = {
  price: 1313,
  image:
    "https://beta.api.solanalysis.com/images/200x200/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/0.png",
};
