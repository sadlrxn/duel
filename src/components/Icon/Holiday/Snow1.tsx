import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import snowImg from 'assets/imgs/holiday/Snow_1.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function Snow1({ size = 100, ...props }: Props) {
  return (
    <Box style={{ pointerEvents: 'none' }} {...props}>
      <Image width={size} height={size * 0.75} src={snowImg} />
    </Box>
  );
}
