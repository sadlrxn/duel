import styled from "styled-components";

import { Box, Label } from "components/index";

export const Input = styled.input`
  flex-basis: 100%;
  background: rgba(3, 6, 9, 0.6);
  border-radius: 11px;
  padding: 17px 25px;
  border: 0;

  font-weight: 400;
  font-size: 16px;
  line-height: 19px;

  color: #ffffff;
`;

export const Image = styled.img`
  width: 21px;
  color: ${({ theme }) => theme.coinflip.private};
`;

export const StyledLabel = styled(Label)`
  display: block;
  font-size: 16px;
  color: ${({ theme }) => theme.coinflip.private};
`;

export const Container = styled(Box)`
  display: flex;
  flex-direction: column;
`;
