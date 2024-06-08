import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import santaImg from 'assets/imgs/holiday/Santa.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function Santa({ size = 74, ...props }: Props) {
  return (
    <Box style={{ pointerEvents: 'none' }} {...props}>
      <Image width={size} height={size} src={santaImg} />
    </Box>
  );
}
