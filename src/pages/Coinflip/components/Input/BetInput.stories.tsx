import React from "react";
import { ComponentMeta } from "@storybook/react";

import BetInput from "./BetInput";

export default {
  title: "Coinflip/BetInput",
  component: BetInput,
  argTypes: {},
} as ComponentMeta<typeof BetInput>;

export const Primary = () => {
  return <BetInput value="5" />;
};
