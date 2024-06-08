import React from "react";
import { ComponentMeta } from "@storybook/react";

import CoinflipList from "./List";

export default {
  title: "Coinflip/CoinflipList",
  component: CoinflipList,
  argTypes: {},
} as ComponentMeta<typeof CoinflipList>;

export const Primary = () => {
  return (
    <>
      <CoinflipList />
    </>
  );
};
