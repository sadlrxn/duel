import React from 'react';
import { FairnessIcon } from 'components/Icon';
import styled from 'styled-components';
import { MarginProps, margin } from 'styled-system';

const StyledButton = styled.button`
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px;
  color: ${({ theme }) => theme.colors.text};
  border: 1px solid;
  border-color: #768bad;
  border-radius: 50%;
  background: transparent;
  font-size: 14px;
  &:hover {
    color: ${({ theme }) => theme.colors.success};
    border-color: ${({ theme }) => theme.colors.success};
    cursor: pointer;

    svg {
      path {
        fill: ${({ theme }) => theme.colors.success};
      }
    }
  }
  ${margin}
`;

interface VerifyButtonProps extends MarginProps {
  onClick?: any;
}

export default function VerifyButton({ onClick, ...props }: VerifyButtonProps) {
  return (
    <StyledButton
      onClick={(e: any) => {
        if (e && e.stopPropagation) e.stopPropagation();
        onClick && onClick();
      }}
      {...props}
    >
      <FairnessIcon size={15} />
    </StyledButton>
  );
}
