import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import santaImg from 'assets/imgs/holiday/Santa_Horse.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function SantaHorse({ size = 250, ...props }: Props) {
  return (
    <Box {...props}>
      <Image width={size} height={size * 0.29} src={santaImg} />
    </Box>
  );
}
