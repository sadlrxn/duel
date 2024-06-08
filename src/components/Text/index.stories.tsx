import React from "react";

import Text from "./Text";

export default {
  title: "Components/Text",
  component: Text,
  argTypes: {
    fontSize: {
      name: "fontSize",
      table: {
        type: { summary: "string", detail: "Fontsize in px or em" },
        defaultValue: { summary: "16px" },
      },
      control: {
        type: null,
      },
    },
    color: {
      name: "color",
      table: {
        type: {
          summary: "string",
          detail: "Color from the theme, or CSS color",
        },
        defaultValue: { summary: "theme.colors.text" },
      },
      control: {
        type: null,
      },
    },
  },
};

export const Primary = () => {
  return (
    <div>
      <Text>Default</Text>
      <Text fontSize="30px">Custom fontsize</Text>
      <Text fontWeight={500}>FontWeight</Text>
      <Text fontWeight={700} fontSize="40px">
        fontweight 700
      </Text>
      <Text color="red">Red Text</Text>
      <Text textAlign="center">Center</Text>
    </div>
  );
};
