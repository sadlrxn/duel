import {
  LayoutProps,
  SpaceProps,
  TypographyProps,
  ColorProps,
  BorderProps
} from 'styled-system';

export interface TextProps
  extends LayoutProps,
    SpaceProps,
    TypographyProps,
    ColorProps,
    BorderProps {
  textTransform?: 'uppercase' | 'lowercase' | 'capitalize';
}
