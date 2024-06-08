import styled from 'styled-components';
import { color, space, typography, layout, border } from 'styled-system';

import { TextProps } from './types';

const P = styled.p<TextProps>`
  margin: 0;
  ${color}
  ${space}
  ${typography}
  ${layout}
  ${border}

  ${({ textTransform }) => textTransform && `text-transform: ${textTransform};`}
`;

export default P;
