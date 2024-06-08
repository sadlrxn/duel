import React from 'react';

import bonkImg from 'assets/imgs/crypto/BONK.png';
import { Flex } from 'components/Box';
import { LazyLoadImage } from 'react-lazy-load-image-component';

export default function Bonk({ size = 64 }: { size?: number }) {
  return (
    <Flex
      width={size}
      height={size}
      justifyContent="center"
      alignItems="center"
      position="relative"
    >
      <LazyLoadImage
        style={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)'
        }}
        src={bonkImg}
        width={size * 1.5}
        height={size * 1.5}
      />
    </Flex>
  );
}
