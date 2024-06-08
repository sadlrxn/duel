import styled from "styled-components";
import {
  background,
  color,
  borders,
  layout,
  space,
  boxShadow,
  BackgroundProps,
  ColorProps,
  BorderProps,
  LayoutProps,
  SpaceProps,
  BoxShadowProps,
} from "styled-system";

import { Flex } from "components/Box";

interface CircleProps
  extends BackgroundProps,
    ColorProps,
    BorderProps,
    LayoutProps,
    SpaceProps,
    BoxShadowProps {}

export const Circle = styled(Flex)<CircleProps>`
  justify-content: center;
  align-items: center;
  border-radius: 100%;

  ${background}
  ${borders}
  ${color}
  ${layout}
  ${space}
  ${boxShadow}
`;

Circle.defaultProps = {};
