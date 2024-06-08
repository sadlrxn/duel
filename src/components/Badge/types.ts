import { HTMLAttributes } from "react";
import { MarginProps, TypographyProps } from "styled-system";
import { BoxProps } from "components/Box";

export const variants = {
  PRIMARY: "primary",
  SECONDARY: "secondary",
} as const;

export type Variant = typeof variants[keyof typeof variants];

export interface BadgeProps
  extends MarginProps,
    TypographyProps,
    BoxProps,
    HTMLAttributes<HTMLDivElement> {
  variant?: Variant;
  notification?: string;
  children?: React.ReactNode | React.ReactNode[];
}
