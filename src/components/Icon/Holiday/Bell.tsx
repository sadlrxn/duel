import styled, { keyframes, css } from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import bellImg from 'assets/imgs/holiday/Bell.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
  isAuto?: boolean;
  isSelected?: boolean;
  roundId?: number;
  opacity?: number;
}

const Image = styled(LazyLoadImage)<Props>`
  ${({ isAuto, isSelected }) => {
    const isZoomIn = isAuto || !isSelected;
    const delay = !isSelected;
    const duration = delay ? 1 : 0.3;
    return css`
      animation: ${isZoomIn ? zoomIn(delay) : zoomOut} ${duration}s ease-out;
    `;
  }}
`;

const zoomIn = (delay?: boolean) => keyframes`
  0% { scale: 0; opacity: 0 }
  ${delay ? '60%' : '0%'} { scale: 0; opacity: 0 }
  100% { scale: 1 }
`;

const zoomOut = keyframes`
  from { scale: 2; opacity: 0 }
  to { scale: 1 }
`;

export default function Bell({
  size = 78,
  roundId = 0,
  opacity = 1,
  isSelected = true,
  isAuto = false,
  ...props
}: Props) {
  return (
    <Box {...props}>
      <Image
        key={roundId}
        width={size}
        height={size}
        src={bellImg}
        style={{ opacity: opacity, marginLeft: `${-size / 7}px` }}
        isAuto={isAuto}
        isSelected={isSelected}
      />
    </Box>
  );
}
