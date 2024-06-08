import styled from 'styled-components';
import { color, space, typography, layout, border } from 'styled-system';

import { TextProps } from './types';

const Text = styled.div<TextProps>`
  line-height: 1.5;
  ${color}
  ${space}
  ${typography}
  ${layout}
  ${border}

  ${({ textTransform }) => textTransform && `text-transform: ${textTransform};`}
`;

export default Text;
