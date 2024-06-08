import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import hatImg from 'assets/imgs/holiday/Hat.png';

import { Box, BoxProps } from 'components/Box';

interface Props extends BoxProps {
  size?: number;
}

const Image = styled(LazyLoadImage)<Props>``;

export default function Snow1({ size = 44, ...props }: Props) {
  return (
    <Box style={{ pointerEvents: 'none' }} {...props}>
      <Image width={size} height={size} src={hatImg} />
    </Box>
  );
}
