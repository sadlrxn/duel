import React from 'react';
import styled from 'styled-components';
import {
  variant as StyledSystemVariant,
  margin,
  typography,
  space,
  background,
  color
} from 'styled-system';
import { Box } from '../Box';
import { styleVariants, notificationStyleVariants } from './theme';
import { BadgeProps, Variant } from './types';

export const StyledBadge = styled(Box)<BadgeProps>`
  position: relative;
  padding: 0.3rem 0.5rem;
  border-radius: 9px;
  font-weight: 600;
  font-size: 12px;
  min-width: max-content;
  max-width: max-content;
  ${StyledSystemVariant({
    variants: styleVariants
  })}
  ${margin}
  ${typography}
  ${space}
  ${background}
  ${color}
`;

StyledBadge.defaultProps = {
  variant: 'primary'
};

export const Notification = styled(Box)<{ variant?: Variant }>`
  position: absolute;
  top: 0px;
  transform: translate(100%, -50%);
  border-radius: 26px;
  padding: 1px 0.5rem;
  text-align: center;
  font-weight: 500;
  font-size: 12px;
  color: #182738;
  ${StyledSystemVariant({
    variants: notificationStyleVariants
  })}
`;

Notification.defaultProps = {
  variant: 'primary',
  right: '12px'
};
export default function Badge({
  children,
  notification,
  ...props
}: BadgeProps) {
  return (
    <StyledBadge {...props}>
      {children}
      {notification ? (
        <Notification {...props}>{notification}</Notification>
      ) : null}
    </StyledBadge>
  );
}
