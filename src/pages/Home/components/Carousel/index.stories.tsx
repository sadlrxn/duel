import React from "react";
import { ComponentMeta } from "@storybook/react";

import Carousel from ".";
import { Flex } from "components/Box";

export default {
  title: "home/Carousel",
  component: Carousel,
  argTypes: {},
} as ComponentMeta<typeof Carousel>;

export const Default = () => (
  <Flex mt={"50px"}>
    <Carousel />
  </Flex>
);
