import styled from 'styled-components';
import { color, space, typography, layout, border } from 'styled-system';

import { TextProps } from './types';

const Span = styled.span<TextProps>`
  line-height: 1.5;
  ${color}
  ${space}
  ${typography}
  ${layout}
  ${border}

  ${({ textTransform }) => textTransform && `text-transform: ${textTransform};`}
`;

export default Span;
