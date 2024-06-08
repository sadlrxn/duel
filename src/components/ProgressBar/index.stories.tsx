import React, { useState } from "react";
import { ComponentMeta, Story } from "@storybook/react";
import random from "lodash/random";
import { Button } from "../Button";
import Box from "../Box/Box";
import Progress from "./Progress";
import { variants, ProgressProps } from "./types";

export default {
  title: "Components/Progress",
  component: Progress,
  argTypes: {},
} as ComponentMeta<typeof Progress>;

const DefaultTemplate: Story<ProgressProps> = (args) => {
  const [progress, setProgress] = useState(random(1, 100));

  const handleClick = () => setProgress(random(1, 100));

  return (
    <div style={{ padding: "32px", width: "400px" }}>
      {Object.values(variants).map((variant) => {
        return (
          <Box key={variant} mb="16px">
            <Progress {...args} variant={variant} step={progress} />
          </Box>
        );
      })}
      Small
      <Progress {...args} scale="sm" step={progress} />
      <div style={{ marginTop: "32px" }}>
        <Button type="button" onClick={handleClick}>
          Random Progress
        </Button>
      </div>
    </div>
  );
};

export const Default = DefaultTemplate.bind({});

Default.args = {
  useDark: false,
};
