import styled from 'styled-components';
import { layout, grid, flexbox } from 'styled-system';

import Box from './Box';
import { GridProps } from './types';

const Grid = styled(Box)<GridProps>`
  display: grid;
  ${layout}
  ${flexbox}
  ${grid}
`;

export default Grid;
