import styled from 'styled-components';
import Grid from './Grid';
import { GridProps } from './types';

const NftBox = styled(Grid)<GridProps>`
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  justify-items: center;
  margin-top: 30px;
  padding: 20px 0px;
  /* overflow: hidden; */
  overflow-y: auto;

  grid-row-gap: 15px;
  grid-column-gap: 20px;

  -ms-overflow-style: none;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }

  ${({ theme }) => theme.mediaQueries.md} {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    grid-row-gap: 25px;
    grid-column-gap: 40px;
  }
`;

NftBox.defaultProps = {
  minHeight: '25vh',
  maxHeight: '350px'
};

export default NftBox;
