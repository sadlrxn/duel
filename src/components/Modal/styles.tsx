import styled from 'styled-components';
import { Box } from 'components/Box';

export const ModalContainer = styled(Box)`
  position: relative;
  z-index: 200;

  .close {
    position: absolute;
    z-index: 20;
    right: 30px;
    top: 40px;
    cursor: pointer;

    &:hover {
      color: white;
    }
  }
`;
