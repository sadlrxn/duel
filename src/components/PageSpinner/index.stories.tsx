import React from "react";
import { ComponentMeta } from "@storybook/react";

import PageSpinner from ".";

export default {
  title: "Components/PageSpinner",
  component: PageSpinner,
  argTypes: {},
} as ComponentMeta<typeof PageSpinner>;

export const Default = () => <PageSpinner />;
