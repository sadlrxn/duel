import styled from "styled-components";

import { Box, Flex, Chip, Text, Avatar as BaseAvatar } from "components";

export const Image = styled(Flex)<{ side: string }>`
  justify-content: center;
  align-items: center;
  width: 50px;
  height: 50px;
  background: ${({ side, theme }) =>
    side === "duel"
      ? theme.coinflip.gradients.duel
      : theme.coinflip.gradients.ana};
  border-radius: 100%;

  display: none;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    display: flex;
  }
`;

export const Divider = styled(Flex)<{ side: string }>`
  height: 50px;
  border-left: 2px solid
    ${({ side, theme }) =>
      side === "duel"
        ? theme.coinflip.duelSecondary
        : theme.coinflip.anaSecondary};

  display: none;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    display: block;
  }
`;

export const Avatar = styled(BaseAvatar)<{ side: string }>`
  background: ${({ theme, side }) =>
    side === "duel"
      ? theme.coinflip.gradients.duel
      : theme.coinflip.gradients.ana};

  padding: 2px;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    padding: 0px;
    background: none;
  }
`;

export const StyledText = styled(Text)`
  font-family: Inter;
  font-size: 14px;
  font-weight: 500;
  line-height: 17px;
  letter-spacing: 0em;
  text-align: left;

  display: none;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    display: block;
  }
`;

export const CrownImage = styled.img<{ show: boolean }>`
  width: 15px;
  height: 15px;
  opacity: ${({ show }) => (show ? 1 : 0)};
  transition: opacity 0.5s;
`;

export const StyledChip = styled(Chip)`
  font-size: 10px;
  font-weight: 500;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    font-size: 14px;
    font-weight: 700;
  }
`;

export const DetailContainer = styled(Flex)<{ side: string }>`
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  align-items: ${({ side }) => (side === "duel" ? "start" : "end")};

  .width_1100 & {
    flex-direction: row;
    align-items: center;
    gap: 19px;
  }
`;

export const DataContainer = styled(Box)<{ side: string }>`
  display: flex;
  flex-direction: column;
  gap: 7px;
  justify-content: space-around;
  align-items: center;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    flex-direction: row-reverse;
    justify-content: space-between;
    /* ${CrownImage} {
      display: none;
    } */
  }
`;

export const Container = styled(Box)<{ side: string; active: boolean }>`
  display: grid;
  grid-template-columns: max-content max-content max-content auto;
  direction: ${({ side }) => (side === "duel" ? "ltr" : "rtl")};
  ${({ active }) => (active ? "" : "opacity: 0.4;")}
  align-items: center;
  padding: 16px 19px;
  gap: 17px;
  transition: all 0.4s;
`;
