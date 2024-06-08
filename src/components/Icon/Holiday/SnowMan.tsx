import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import snowManImg from 'assets/imgs/holiday/SnowMan.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function SnowMan({ size = 142, ...props }: Props) {
  return (
    <Box style={{ pointerEvents: 'none' }} {...props}>
      <Image width={size} height={size} src={snowManImg} />
    </Box>
  );
}
