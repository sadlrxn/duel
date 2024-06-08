import styled, { css } from 'styled-components';
import {
  color,
  typography,
  layout,
  border,
  background,
  flexbox,
  TypographyProps,
  ColorProps,
  LayoutProps,
  BorderProps,
  BackgroundProps,
  FlexboxProps
} from 'styled-system';

interface InputProps
  extends ColorProps,
    TypographyProps,
    LayoutProps,
    BorderProps,
    BackgroundProps,
    FlexboxProps {
  placeholderColor?: string;
}

const Input = styled.input<InputProps>`
  -webkit-appearance: none;
  border: none;
  outline: none;
  /* font: v.$font; */
  background: transparent;

  &[type='number'] {
    -moz-appearance: textfield;

    &::-webkit-inner-spin-button,
    &::-webkit-outer-spin-button {
      -webkit-appearance: none;
      margin: 0;
    }
  }

  &::placeholder {
    ${({ placeholderColor }) =>
      placeholderColor &&
      css`
        color: ${placeholderColor};
      `}
  }

  ${color}
  ${typography}
  ${layout}
  ${border}
  ${background}
  ${flexbox}
`;

export default Input;
