import React from 'react';

import solImg from 'assets/imgs/crypto/Solana.png';
import { Flex } from 'components/Box';
import { LazyLoadImage } from 'react-lazy-load-image-component';

export default function Sol({ size = 24 }: { size?: number }) {
  return (
    <Flex
      width={size}
      height={size}
      justifyContent="center"
      alignItems="center"
    >
      <LazyLoadImage src={solImg} width={size} height={size} />
    </Flex>
  );
}
