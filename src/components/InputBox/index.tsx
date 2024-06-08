import styled from 'styled-components';
import { Flex } from 'components/Box';

export const InputBox = styled(Flex)`
  align-items: center;
  background: #03060999;
  border-radius: 11px;

  &:focus-within {
    box-shadow: 0 0 0 1px #7389a9;
  }

  input {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 400;
    font-size: 20px;
    line-height: 17px;
    color: #8192aa;
    flex-grow: 1;
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
