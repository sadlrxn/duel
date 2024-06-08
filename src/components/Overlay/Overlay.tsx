import { FC } from 'react';
import styled, { css, keyframes } from 'styled-components';
import { Box, BoxProps } from '../Box';

const unmountAnimation = keyframes`
    0% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  `;

const mountAnimation = keyframes`
    0% {
     opacity: 0;
    }
    100% {
     opacity: 1;
    }
  `;

const StyledOverlay = styled(Box)<{ isUnmounting?: boolean }>`
  position: fixed;
  top: 0px;
  left: 0px;
  width: 100%;
  height: 100%;
  background-color: ${({ theme }) => `${theme.colors.overlay}`};
  z-index: 150;
  will-change: opacity;
  animation: ${mountAnimation} 350ms ease forwards;
  ${({ isUnmounting }) =>
    isUnmounting &&
    css`
      animation: ${unmountAnimation} 350ms ease forwards;
    `}
`;

interface OverlayProps extends BoxProps {
  isUnmounting?: boolean;
}

export const Overlay: FC<OverlayProps> = props => {
  return (
    <>
      {/* <BodyLock /> */}
      <StyledOverlay role="presentation" {...props} />
    </>
  );
};

export default Overlay;
