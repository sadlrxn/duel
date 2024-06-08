import { forwardRef } from 'react';
import styled, { css } from 'styled-components';
import {
  layout,
  space,
  background,
  borders,
  variant,
  color
} from 'styled-system';
import { useSound } from 'hooks';

import { scaleVariants, styleVariants } from './theme';

import BaseButton from './BaseButton';
import { ButtonProps } from './types';

const InnerButton = styled(BaseButton)<ButtonProps>`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;

  &:hover {
    filter: brightness(110%);
  }

  &:active {
    filter: brightness(80%);
  }

  ${({ disabled }) => {
    if (disabled)
      return css`
        opacity: 0.75;
        cursor: not-allowed;
        pointer-events: none;
      `;
  }}

  ${({ nonClickable }) => {
    if (nonClickable)
      return css`
        cursor: not-allowed;
        pointer-events: none;
      `;
  }}

  ${variant({
    prop: 'scale',
    variants: scaleVariants
  })}
  ${variant({
    variants: styleVariants
  })}

  border-style: solid;

  ${({ outlined }) => (outlined ? 'border-width: 2px;' : '')}

  ${layout}
  ${space}
  ${background}
  ${borders}
  ${color}
`;

InnerButton.defaultProps = {
  variant: 'primary',
  scale: 'default',
  outlined: false,
  type: 'button'
};

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ children, onClick, sound = true, tabIndex, ...props }, ref) => {
    const { buttonPlay } = useSound();
    return (
      <InnerButton
        tabIndex={tabIndex}
        onClick={(e: any) => {
          if (e && e.stopPropagation) e.stopPropagation();
          sound && buttonPlay();
          onClick && onClick();
          return false;
        }}
        {...props}
        ref={ref}
      >
        {children}
      </InnerButton>
    );
  }
);

Button.displayName = 'Button';

export default Button;
