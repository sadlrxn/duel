import React, { ReactNode } from "react";

import { StyledButton, Text } from "./Tooltip.styles";

interface IToolTipProps {
  icon: ReactNode;
  text: string;
  children?: React.ReactNode;
  className?: string;
  onClick?: () => void;
  tooltipPosition?: "top" | "bottom";
}

export default function Tooltip({ icon, text }: IToolTipProps) {
  return (
    <StyledButton>
      {icon}
      <Text>{text}</Text>
    </StyledButton>
  );
}
