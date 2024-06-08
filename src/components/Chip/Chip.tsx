import React from 'react';

import { Span } from 'components/Text';

import { Container, ChipIcon } from './Chip.styles';
import { formatNumber } from 'utils/format';

const chipColors = {
  chip: {
    background: '#ffe24b',
    border: '#ffb31f'
  },
  coupon: {
    background: '#4BE9FF',
    border: '#1FAEFF'
  }
};

export default function Chip({
  price,
  fontSize = '14px',
  fontWeight = 500,
  letterSpacing = 1,
  color = 'chip',
  chipType = 'chip',
  size,
  background,
  border,
  prefix,
  decimal = 0,
  ...props
}: any) {
  return (
    <Container fontSize={fontSize} {...props}>
      <ChipIcon
        $size={size}
        //@ts-ignore
        $background={background ? background : chipColors[chipType].background}
        //@ts-ignore
        $border={border ? border : chipColors[chipType].border}
      />
      <Span
        fontWeight={fontWeight}
        letterSpacing={letterSpacing}
        fontSize={fontSize}
        color={color}
      >
        {prefix !== undefined && prefix}
        {formatNumber(price, decimal)}
      </Span>
    </Container>
  );
}
