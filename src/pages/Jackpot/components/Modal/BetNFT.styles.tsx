import styled from "styled-components";

import { Flex, Button } from "components";

export const InputContainer = styled(Flex)`
  align-items: center;
  gap: 20px;
  padding: 10px 20px;
  background: #03060999;
  border-radius: 11px;

  &:focus-within {
    box-shadow: 0 0 0 1px #7389a9;
  }

  input {
    font-family: "Inter";
    font-style: normal;
    font-weight: 400;
    font-size: 14px;
    line-height: 17px;
    color: #8192aa;

    border: none;
    outline: none;
    background: transparent;

    &::placeholder {
      color: #8192aa;
    }

    ::-webkit-outer-spin-button,
    ::-webkit-inner-spin-button {
      -webkit-appearance: none;
      margin: 0;
    }
    -moz-appearance: textfield;
  }
`;

export const WagerButton = styled(Button)`
  display: flex;
  gap: 8px;
  padding-left: 10px;
  padding-right: 10px;
  background: transparent;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  color: ${({ theme }) => theme.coinflip.title};
  border: 1px solid ${({ theme }) => theme.coinflip.title};
  border-radius: 6px;
  height: 52px;
`;
WagerButton.defaultProps = { variant: "secondary" };
