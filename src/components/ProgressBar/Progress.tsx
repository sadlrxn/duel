import React from "react";
import StyledProgress, { Bar } from "./styles";
import { ProgressProps, variants, scales } from "./types";

const stepGuard = (step: number) => {
  if (step < 0) {
    return 0;
  }

  if (step > 100) {
    return 100;
  }

  return step;
};

const Progress: React.FC<ProgressProps> = ({
  variant = variants.ROUND,
  scale = scales.MD,
  step = 0,
  useDark = true,
  color,
}) => {
  return (
    <StyledProgress $useDark={useDark} variant={variant} scale={scale}>
      <Bar
        $useDark={useDark}
        primary
        style={{ width: `${stepGuard(step)}%` }}
        $background={color}
      />
    </StyledProgress>
  );
};

export default Progress;
