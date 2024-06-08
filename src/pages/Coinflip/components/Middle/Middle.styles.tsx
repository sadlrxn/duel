import styled from "styled-components";

import { Box, Flex, Button, Span } from "components";

const SIDE_COUNT = 30;
const CHIP_WIDTH = 3;

export const Div = styled(Flex)`
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;

  border-radius: 50%;
  /* margin: 10px; */

  position: absolute;
  overflow: hidden;
  backface-visibility: hidden;
`;

export const Duel = styled(Div)`
  background: ${({ theme }) => theme.colors.gradients.duel};
  transform: translateZ(${CHIP_WIDTH}px);
`;

export const Ana = styled(Div)`
  background: ${({ theme }) => theme.colors.gradients.ana};
  transform: rotateY(-180deg);
  transform: translateZ(-${CHIP_WIDTH}px) rotateY(-180deg);
`;

export const Side = styled(Box)<{ index: number }>`
  position: absolute;
  height: ${CHIP_WIDTH * 2}px;
  width: 12px;
  top: calc(50% - ${CHIP_WIDTH}px);
  left: calc(50% - 6px);
  background-color: ${({ index }) =>
    index % 2 === 0 ? "#85ffe288" : "#5d24ff88"};
  transform: rotateX(90deg)
    rotateY(${({ index }) => (index * 360) / SIDE_COUNT}deg) translateZ(24px);
`;

export const Coin = styled(Box)`
  position: relative;
  width: 50px;
  height: 50px;
  transform-style: preserve-3d;
  transform: perspective(800px);
`;

export const CounterText = styled(Span)`
  font-weight: bold;
  font-size: 20px;
  color: #ffffff;
  width: 30px;
  height: 30px;
  text-align: center;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
`;

export const StyledButton = styled(Button)`
  gap: 7px;
  padding: 7px 17px;
  border-radius: 5px;
  background: transparent;
  border-radius: 5px;
`;

export const Container = styled(Flex)`
  /* width: 0px; */
  position: relative;
  height: 100%;
  gap: 10px;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;
