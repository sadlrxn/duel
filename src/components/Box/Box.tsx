import styled from 'styled-components';
import {
  background,
  border,
  layout,
  position,
  space,
  typography,
  shadow
} from 'styled-system';

import { BoxProps } from './types';

const Box = styled.div<BoxProps>`
  ${background}
  ${border}
  ${layout}
  ${position}
  ${space}
  ${typography}
  ${shadow}
  ${({ gap }) =>
    typeof gap === 'number'
      ? `gap: ${gap}px`
      : typeof gap === 'string'
      ? `gap: ${gap}`
      : ''};
  ${({ cursor }) => cursor && `cursor: ${cursor};`}
`;

Box.defaultProps = {
  gap: 0
};

export default Box;
