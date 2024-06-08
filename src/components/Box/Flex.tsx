import styled from 'styled-components';
import { layout, flexbox } from 'styled-system';

import Box from './Box';
import { FlexProps } from './types';

const Flex = styled(Box)<FlexProps>`
  display: flex;
  ${layout}
  ${flexbox};
`;
export default Flex;
