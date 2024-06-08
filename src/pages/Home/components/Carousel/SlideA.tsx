import React from "react";
import styled from "styled-components";
import { Flex, Text } from "components";
import { StyledSlideA } from "./styles";
import mascotImg from "assets/imgs/home/carousel/mascot.png";
import { Floating } from "utils/animationToolkit";

export default function SlideA() {
  return (
    <StyledSlideA>
      <Intro>
        <Welcome>WELCOME TO</Welcome>
        <Arena>THE ARENA</Arena>
        <Dueler>DUELER</Dueler>
      </Intro>
      <Description>
        Challenge your friends and other duelers to a game of Coin Flip or
        Jackpot. Win some CHIPS or rare NFTs to add to your collection.
      </Description>

      <Mascot src={mascotImg} alt="mascot" />
    </StyledSlideA>
  );
}

const Intro = styled(Flex)`
  flex-direction: column;
  align-items: center;
  z-index: 10;
  .width_700 & {
    align-items: start;
    width: 350px;
  }
`;

const Welcome = styled(Text)`
  color: white;
  font-family: "Termina";
  font-weight: 800;
  font-size: 20px;
  letter-spacing: -0.02em;
  text-align: center;
  .width_700 & {
    font-size: 24.2821px;
  }
`;

const Arena = styled(Text)`
  color: white;
  font-family: "Termina";
  font-weight: 800;
  font-size: 38px;
  line-height: 46px;
  letter-spacing: -0.02em;
  text-align: center;
  .width_700 & {
    font-size: 42.9587px;
  }
`;

const Dueler = styled(Arena)`
  color: #4fff8b;
`;

const Description = styled(Text)`
  color: #d6d6d6;
  font-weight: 500;
  letter-spacing: -0.02em;
  font-size: 16.71px;
  line-height: 26px;
  text-align: center;
  z-index: 10;
  .width_700 & {
    width: 400px;
    font-size: 16px;
    text-align: left;
  }
`;

const Mascot = styled.img`
  position: absolute;
  width: 200px;
  bottom: 10px;
  /* right: 50%; */
  /* transform: translate(-50%, -50%); */

  animation: ${Floating} 1.5s ease-in-out infinite alternate;
  animation-delay: 0.5s;

  .width_700 & {
    right: 10%;
    width: 274px;
  }
`;
