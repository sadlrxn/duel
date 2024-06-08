import styled from "styled-components";

import { Flex } from "components/Box";

export const Text = styled.span`
  color: #ffffff;
  font-family: "Inter";
  font-style: normal;
  font-weight: 500;
  font-size: 25px;
  line-height: 34px;
`;

export const Count = styled.span`
  color: #4f617b;
  font-family: "Inter";
  font-style: normal;
  font-weight: 400;
  font-size: 14px;
  line-height: 30px;
`;

export const DataContainer = styled(Flex)`
  flex-direction: row;
  align-items: end;
  gap: 14px;
`;

export const ButtonContainer = styled(Flex)`
  flex-direction: row;
  justify-content: end;
  gap: 12px;
`;

export const Heading = styled(Flex)`
  flex-direction: column;
  margin-bottom: 40px;
  gap: 5px;

  ${({ theme }) => theme.mediaQueries.md} {
    flex-direction: row;
    justify-content: space-between;
  }
`;

export const GameList = styled(Flex)`
  flex-direction: column;
  width: 100%;

  .enter,
  .appear {
    opacity: 0.01;
  }

  .enter.enter-active,
  .appear.appear-active {
    opacity: 1;
    transition: opacity 250ms ease-in;
  }

  .exit {
    opacity: 1;
  }

  .exit.exit-active {
    opacity: 0.01;
    transition: opacity 250ms ease-out;
  }
`;

export const Container = styled.section`
  position: relative;
  z-index: 0;
`;
