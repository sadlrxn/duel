import { ComponentMeta, Story } from "@storybook/react";
import Content, { ContentProps } from "./Content";
import random from "lodash/random";

export default {
  title: "Jackpot/ContentPage",
  component: Content,
  argTypes: {},
} as ComponentMeta<typeof Content>;

const Template: Story<ContentProps> = (args) => <Content {...args} />;

export const Default = Template.bind({});

Default.args = {
  nfts: [
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
    {
      mintAddress: "4w2VY8V877KbuRE9wa9F9T4JT5eLruRqjuAZxTxPzXGr",
      image:
        "https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://nftstorage.link/ipfs/bafybeibj2xl7aorp2kpas6wlixsai5pvm3dsfwwo5ah4jyjfdqvqmwoyby/140.png",
      price: random(5000, 20000),
    },
  ],
};
