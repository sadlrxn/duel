import React, { useState } from 'react';
import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import chip1 from 'assets/imgs/coins/Chip1.svg';
import chip2 from 'assets/imgs/coins/Chip2.svg';
import chip5 from 'assets/imgs/coins/Chip5.svg';
import chip10 from 'assets/imgs/coins/Chip10.svg';
import chip25 from 'assets/imgs/coins/Chip25.svg';
import chip50 from 'assets/imgs/coins/Chip50.svg';
import chip100 from 'assets/imgs/coins/Chip100.svg';
import chip250 from 'assets/imgs/coins/Chip250.svg';
import chip500 from 'assets/imgs/coins/Chip500.svg';

import chip1_hover from 'assets/imgs/coins/Chip1_hover.svg';
import chip2_hover from 'assets/imgs/coins/Chip2_hover.svg';
import chip5_hover from 'assets/imgs/coins/Chip5_hover.svg';
import chip10_hover from 'assets/imgs/coins/Chip10_hover.svg';
import chip25_hover from 'assets/imgs/coins/Chip25_hover.svg';
import chip50_hover from 'assets/imgs/coins/Chip50_hover.svg';
import chip100_hover from 'assets/imgs/coins/Chip100_hover.svg';
import chip250_hover from 'assets/imgs/coins/Chip250_hover.svg';
import chip500_hover from 'assets/imgs/coins/Chip500_hover.svg';

import { Flex } from 'components/Box';

const Container = styled(Flex)<{ isClick: boolean }>`
  ${({ isClick }) => isClick && 'transform: translateY(5px);'}
`;

interface ChipProps {
  price?: number;
  size?: number;
  onClick?: any;
  disabled?: boolean;
}

export default function Chip({
  price = 1,
  size = 44,
  disabled = false,
  onClick
}: ChipProps) {
  const chip = {
    chip1,
    chip2,
    chip5,
    chip10,
    chip25,
    chip50,
    chip100,
    chip250,
    chip500
  };
  const chip_hover = {
    chip1: chip1_hover,
    chip2: chip2_hover,
    chip5: chip5_hover,
    chip10: chip10_hover,
    chip25: chip25_hover,
    chip50: chip50_hover,
    chip100: chip100_hover,
    chip250: chip250_hover,
    chip500: chip500_hover
  };

  //@ts-ignore
  const [src, setSrc] = useState(chip['chip' + price]);
  const [isMouseOn, setIsMouseOn] = useState(false);
  const [clicked, setClicked] = useState(false);

  return (
    <Container
      width={size}
      height={size * 0.887}
      onMouseEnter={() => {
        //@ts-ignore
        setSrc(chip_hover['chip' + price]);
        setIsMouseOn(true);
      }}
      onMouseLeave={() => {
        //@ts-ignore
        setSrc(chip['chip' + price]);
        setIsMouseOn(false);
        setClicked(false);
      }}
      onMouseDown={() => {
        setClicked(true);
      }}
      onMouseUp={() => {
        setClicked(false);
      }}
      onClick={() => {
        onClick && onClick();
      }}
      cursor="pointer"
      position="relative"
      style={{
        pointerEvents: disabled ? 'none' : 'all',
        opacity: disabled ? '0.4' : '1'
      }}
      isClick={clicked && isMouseOn}
    >
      <LazyLoadImage
        style={{
          position: 'absolute',
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)'
        }}
        src={src}
        alt=""
      />
    </Container>
  );
}
