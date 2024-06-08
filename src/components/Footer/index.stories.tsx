import React from "react";
import { ComponentStory, ComponentMeta } from "@storybook/react";

import Footer from ".";

export default {
  title: "Components/Footer",
  component: Footer,
  argTypes: {},
} as ComponentMeta<typeof Footer>;

export const Default: ComponentStory<typeof Footer> = () => <Footer />;
