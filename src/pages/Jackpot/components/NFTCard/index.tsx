import React from 'react';
import styled from 'styled-components';
import {
  LazyLoadImage,
  LazyLoadImageProps
} from 'react-lazy-load-image-component';

import coin from 'assets/imgs/coins/win.png';

import { Box, Chip, Text } from 'components';
import { imageProxy } from 'config';

export interface StyledImageProps extends LazyLoadImageProps {
  size: string | number;
}

const StyledImage = styled(LazyLoadImage)<StyledImageProps>`
  width: ${({ size }) => size}px;
  height: ${({ size }) => size}px;
  border-radius: 6px;
`;

const CoinImage = styled.img`
  border-radius: 6px;
  border: 2px solid ${({ theme }) => theme.colors.chip};
  padding: 10px;
  background: linear-gradient(
    180deg,
    rgba(255, 162, 39, 0.1) 0%,
    rgba(255, 162, 39, 0.2) 100%,
    rgba(255, 162, 39, 0.2) 100%
  );
`;

const StyledText = styled(Text)`
  font-weight: 600;
  font-size: 16px;
  line-height: 19px;

  border-radius: 17px;
  background: #6d81a2;
  text-align: center;
  color: #2e3d53;
  padding: 1px 7px;
`;

const More = styled(Box)`
  display: flex;
  flex-direction: column;
  gap: 4px;
  justify-content: center;
  align-items: center;
  border-radius: 6px;
  background: #2e3d53;
  font-weight: 500;
  font-size: 16px;
  line-height: 19px;
  color: #7688ad;
  cursor: pointer;
`;

const Container = styled(Box)`
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
`;

interface NFTCardProps {
  type?: 'chip' | 'nft' | 'more';
  size?: number;
  price?: number;
  image?: string;
  clickable?: boolean;
  onClick?: any;
}

export default function NFTCard({
  type = 'chip',
  size = 86,
  price = 0,
  image = '',
  clickable = true,
  onClick
}: NFTCardProps) {
  return (
    <Container onClick={onClick} cursor={clickable ? 'pointer' : 'default'}>
      {type === 'nft' ? (
        <StyledImage src={imageProxy(300) + image} alt="" size={size} />
      ) : type === 'chip' ? (
        <CoinImage src={coin} alt="" width={size} height={size} />
      ) : (
        <More size={size}>
          <StyledText>{price}</StyledText>
          More
        </More>
      )}
      {type !== 'more' && (
        <Chip price={price} color={type === 'chip' ? 'chip' : 'success'} />
      )}
    </Container>
  );
}
