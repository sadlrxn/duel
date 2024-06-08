import { HTMLAttributes } from 'react';

import {
  BackgroundProps,
  BorderProps,
  FlexboxProps,
  LayoutProps,
  PositionProps,
  SpaceProps,
  TypographyProps,
  ShadowProps,
  GridProps as _GridProps
} from 'styled-system';

export interface BoxProps
  extends BackgroundProps,
    BorderProps,
    LayoutProps,
    PositionProps,
    SpaceProps,
    TypographyProps,
    ShadowProps,
    HTMLAttributes<HTMLDivElement> {
  gap?: number | string;
  cursor?: string;
}

export interface FlexProps extends BoxProps, FlexboxProps {}

export interface GridProps extends FlexProps, _GridProps {}
