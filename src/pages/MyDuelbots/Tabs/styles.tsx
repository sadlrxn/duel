import { Grid } from 'components';
import styled from 'styled-components';

export const StyledDuelbotsContainer = styled(Grid)`
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  justify-items: center;
  margin-top: 20px;

  gap: 22px;
  min-height: 350px;
`;
