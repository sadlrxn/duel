import React from 'react';

import usdcImg from 'assets/imgs/crypto/USDC.png';
import { Flex } from 'components/Box';
import { LazyLoadImage } from 'react-lazy-load-image-component';

export default function Sol({ size = 64 }: { size?: number }) {
  return (
    <Flex
      width={size}
      height={size}
      justifyContent="center"
      alignItems="center"
    >
      <LazyLoadImage src={usdcImg} width={size} height={size} />
    </Flex>
  );
}
