import React from "react";
import { ComponentMeta } from "@storybook/react";

import TooltipButton from "./Tooltip";
import History from "./History";
import Fee from "./Fee";
import Fairness from "./Fairness";

export default {
  title: "Coinflip/Buttons/TooltipButton",
  component: TooltipButton,
  argTypes: {},
} as ComponentMeta<typeof TooltipButton>;

export const Primary = () => {
  return (
    <>
      <History />
      <Fee percentage={2} />
      <Fairness />
    </>
  );
};
