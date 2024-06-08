import { ButtonHTMLAttributes } from 'react';
import {
  LayoutProps,
  SpaceProps,
  BorderProps,
  TypographyProps,
  ColorProps,
  BackgroundProps
} from 'styled-system';

export interface BaseButtonProps
  extends LayoutProps,
    SpaceProps,
    BorderProps,
    TypographyProps,
    BackgroundProps,
    ColorProps {}

export interface ButtonProps
  extends BaseButtonProps,
    ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  scale?: Scale;
  outlined?: boolean;
  color?: string;
  nonClickable?: boolean;
  onClick?: any;
  sound?: boolean;
}

export const scales = {
  DEFAULT: 'default',
  LG: 'lg',
  MD: 'md',
  SM: 'sm',
  XS: 'xs'
} as const;

export const variants = {
  PRIMARY: 'primary',
  SECONDARY: 'secondary',
  DEPOSIT: 'deposit',
  ICON: 'icon',
  PRIVATE: 'private',
  CONNECT: 'connect'
} as const;

export type Scale = typeof scales[keyof typeof scales];
export type Variant = typeof variants[keyof typeof variants];
