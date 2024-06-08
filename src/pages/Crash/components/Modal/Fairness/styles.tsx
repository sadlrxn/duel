import styled from 'styled-components';

import { Box, Flex, Grid, Button, Text, Modal } from 'components';
import { ReactComponent as BaseChipIcon } from 'assets/imgs/coins/coin.svg';

export const StyledModal = styled(Modal)`
  max-width: 1000px;

  font-size: 20px;
  border-radius: 20px;
  background: linear-gradient(180deg, #6a7f9e 0%, rgba(106, 127, 158, 0) 100%);
  padding: 2px;

  width: calc(100vw - 15px);
  max-height: calc(100vh - 65px);

  ${({ theme }) => theme.mediaQueries.sm} {
    max-height: 85vh;
    width: 95vw;
  }
`;

export const Container = styled(Box)`
  max-height: 70vh;

  display: flex;
  flex-direction: column;

  overflow: auto;
  padding-right: 5px;
`;

export const TitleContainer = styled(Grid)`
  display: flex;
  color: white;
  font-size: 1em;
  font-weight: 600;
  line-height: 1.2;
  letter-spacing: 0.18em;
  gap: 14px;

  margin-bottom: 30px;
`;

export const CopyButton = styled(Button)`
  width: 28px;
  height: 28px;

  svg {
    width: 12px;
    height: 12px;
  }
`;

CopyButton.defaultProps = { variant: 'secondary' };

export const InputBox = styled(Flex)<{ readOnly?: boolean }>`
  align-items: center;
  gap: 20px;
  padding: 10px 20px;
  background: ${({ readOnly }) => (readOnly ? '#121C2A' : '#03060999')};
  border-radius: 11px;
  flex-grow: 1;

  &:focus-within {
    box-shadow: 0 0 0 1px #7389a9;
  }

  input {
    font-weight: 400;
    font-size: 14px;
    line-height: 17px;
    color: white;
    flex-grow: 1;
    border: none;
    outline: none;
    background: transparent;

    &::placeholder {
      color: #ffffff80;
    }

    ::-webkit-outer-spin-button,
    ::-webkit-inner-spin-button {
      -webkit-appearance: none;
      margin: 0;
    }
    -moz-appearance: textfield;
  }
`;

export const VerifyButton = styled(Button)`
  border: 2px solid ${({ theme }) => theme.colors.success};
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0.3) 100%);
  border-radius: 7px;

  font-size: 14px;
  font-weight: 600;
  line-height: 17px;
  letter-spacing: 16%;
  text-transform: uppercase;

  height: min-content;

  width: 100%;

  ${({ theme }) => theme.mediaQueries.md} {
    width: max-content;
  }

  a {
    display: flex;
    justify-content: center;
    align-items: center;

    color: white;
    width: 100%;
    height: 100%;

    padding: 7px 14px;
  }

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
  }
`;

VerifyButton.defaultProps = { variant: 'secondary' };

export const Title = styled(Text)`
  color: #d0daeb;
  font-size: 18px;
  line-height: 22px;
  font-weight: 600;
  letter-spacing: 0.1em;

  text-transform: uppercase;
`;

export const Description = styled(Text)`
  color: #b9d2fd;
  font-size: 14px;
  line-height: 20px;
  font-weight: 400;
`;

export const Divider = styled(Box)`
  height: 1px;
  background: #ffffff1a;

  margin-top: 33px;
  margin-bottom: 25px;
`;

export const GrayButton = styled(Button)`
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(180deg, #2a3d57 0%, #2a3d57 100%);
  height: 52px;

  a {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    color: #8192aa;
  }
`;

GrayButton.defaultProps = { variant: 'secondary' };

export const StyledTd = styled.td``;

export const Table = styled.table`
  position: relative;
  width: 100%;
  min-width: 650px;

  font-size: 16px;
  font-weight: 400;
  line-height: 19px;
  color: #ffffff;

  border-collapse: separate;
  border-spacing: 0;

  /* ${({ theme }) => theme.mediaQueries.xs} {
    min-width: auto;
  } */

  th,
  td {
    text-align: left;
    vertical-align: middle;
    padding: 9px 20px;
  }

  ${StyledTd} {
    padding-left: 20px;
  }

  th {
    position: sticky;
    top: 0;
    color: #768bad;
    background: #162536;

    &:last-child {
      text-align: right;
    }
  }

  td {
    background: #0306094c;
  }

  tr {
    td:last-child {
      text-align: right;
    }
  }

  tr:first-child {
    td {
      padding-top: 20px;

      &:first-child {
        border-top-left-radius: 10px;
      }

      &:last-child {
        border-top-right-radius: 10px;
      }
    }
  }

  tr:last-child {
    td {
      padding-bottom: 20px;

      &:first-child {
        border-bottom-left-radius: 10px;
      }

      &:last-child {
        border-bottom-right-radius: 10px;
      }
    }
  }
`;

export const ChipIcon = styled(BaseChipIcon)`
  min-width: 14px;
  max-width: 14px;
  min-height: 14px;
  max-height: 14px;
`;
