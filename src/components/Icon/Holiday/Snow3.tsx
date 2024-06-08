import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import snowImg from 'assets/imgs/holiday/Snow_3.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function Snow3({ size = 80, ...props }: Props) {
  return (
    <Box style={{ pointerEvents: 'none' }} {...props}>
      <Image width={size} height={size * 0.625} src={snowImg} />
    </Box>
  );
}
