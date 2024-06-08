import styled from "styled-components";

import { BaseButton } from "components/Button";

export const Text = styled.span`
  font: 500 14px "Inter";
`;

export const Image = styled.img`
  width: 16px;
  height: 16px;
`;

export const StyledButton = styled(BaseButton)`
  position: relative;

  display: flex;
  align-items: center;
  justify-content: center;
  /* gap: 0.3rem; */
  gap: 7px;

  background: #070c12;
  color: #4f617b;
  font: 400 14px / 16px "Inter";
  padding: 7px 12px;
  border-left: 2px solid #4f617b;
  opacity: 0.8;
  cursor: default;

  transition: 0.5s;

  &:hover {
    background: #0f1a27 !important;
  }
`;
