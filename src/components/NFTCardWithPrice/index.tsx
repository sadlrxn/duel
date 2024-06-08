import React from 'react';
import styled from 'styled-components';
import { Flex } from '../Box';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import { Chip } from '../Chip';
import { imageProxy } from 'config';

export interface NFTCardWithPriceProps {
  address: string;
  image: string;
  height?: number;
  price: number;
}

const StyledLazyLoadImage = styled(LazyLoadImage)`
  border-radius: 5px;
`;

export default function NFTCardWithPrice({
  image,
  price,
  height = 86.73
}: NFTCardWithPriceProps) {
  const chipColor = price > 10000 ? 'secondary' : undefined;
  return (
    <Flex flexDirection="column" alignItems="center" gap={2}>
      <StyledLazyLoadImage
        src={imageProxy(300) + image}
        height={height}
        alt={image}
      />
      <Chip price={price} color={chipColor} />
    </Flex>
  );
}
