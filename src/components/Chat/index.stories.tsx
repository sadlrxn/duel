import React from "react";
import { ComponentStory, ComponentMeta } from "@storybook/react";
import { withRouter } from "storybook-addon-react-router-v6";
import Provider from "../../Providers";
import Chat from ".";

const withProvider = (story: any) => <Provider>{story()}</Provider>;

export default {
  title: "Components/Chat",
  component: Chat,
  argTypes: {},
  decorators: [withRouter, withProvider],
} as ComponentMeta<typeof Chat>;

export const Default: ComponentStory<typeof Chat> = () => <Chat />;
