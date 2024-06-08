import React from 'react';
import styled from 'styled-components';
import { Flex, Text, Button } from 'components';
import { ReactComponent as TwitterIcon } from 'assets/imgs/icons/twitter.svg';
import { StyledSlideC } from './styles';
import chip1Img from 'assets/imgs/home/carousel/chip_1.png';
import chip2Img from 'assets/imgs/home/carousel/chip_2.png';
import chip3Img from 'assets/imgs/home/carousel/chip_3.png';
import chip4Img from 'assets/imgs/home/carousel/chip_4.png';
import chip5Img from 'assets/imgs/home/carousel/chip_5.png';
import chip6Img from 'assets/imgs/home/carousel/chip_6.png';
import chip7Img from 'assets/imgs/home/carousel/chip_7.png';
import chip8Img from 'assets/imgs/home/carousel/chip_8.png';
import { Floating } from 'utils/animationToolkit';

export default function SlideC() {
  return (
    <StyledSlideC>
      <Intro>
        <DailyGems>DAILY CHIPS</DailyGems>
        <Giveaway>GIVEAWAY</Giveaway>
      </Intro>

      <Description>
        Like, retweet and tag two friends for a chance to win $10 CHIPS credit
        to your Duel account.
      </Description>
      <a
        href="https://twitter.com/DuelCasino/status/1608704536218341376"
        target="_blank"
        rel="noreferrer"
      >
        <Button
          color="white"
          fontWeight={700}
          fontSize="12px"
          p="12px 15px"
          background={'#1DA1F2'}
        >
          <TwitterIcon color="white" />
          GIMME MY CHIPS
        </Button>
      </a>

      <Chip1 src={chip1Img} alt="chip-1" />
      <Chip2 src={chip2Img} alt="chip-2" />
      <Chip3 src={chip3Img} alt="chip-3" />
      <Chip4 src={chip4Img} alt="chip-4" />
      <Chip5 src={chip5Img} alt="chip-5" />
      <Chip6 src={chip6Img} alt="chip-6" />
      <Chip7 src={chip7Img} alt="chip-7" />
      <Chip8 src={chip8Img} alt="chip-8" />
    </StyledSlideC>
  );
}

const Intro = styled(Flex)`
  margin-top: 30px;
  flex-direction: column;
  align-items: center;
  z-index: 10;
`;

const DailyGems = styled(Text)`
  color: #4fff8b;
  font-family: 'Termina';
  font-weight: 800;
  font-size: 20px;
  letter-spacing: -0.02em;
  text-align: center;
  .width_700 & {
    font-size: 24px;
  }
`;

const Giveaway = styled(Text)`
  color: white;
  font-family: 'Termina';
  font-weight: 800;
  font-size: 38px;
  line-height: 46px;
  letter-spacing: -0.02em;
  text-align: center;
  .width_700 & {
    font-size: 42.9587px;
    line-height: 52px;
  }
`;

const Description = styled(Text)`
  color: #d6d6d6;
  z-index: 10;
  font-weight: 500;
  letter-spacing: -0.02em;
  font-size: 16.71px;
  line-height: 26px;
  text-align: center;
  .width_700 & {
    width: 400px;
    font-size: 16px;
    text-align: center;
  }
`;

const Chip = styled.img`
  position: absolute;
  /* width: calc(100% * 274 / 1000); */

  animation: ${Floating} 1.5s ease-in-out infinite alternate;
`;

const Chip1 = styled(Chip)`
  top: -10px;
  left: -15%;
  animation-delay: 1.5s;

  .width_700 & {
    top: -50px;
    left: 10%;
  }
`;

const Chip2 = styled(Chip)`
  bottom: 30px;
  right: 3%;
  width: 100px;
  animation-delay: 0.5s;

  .width_700 & {
    width: 110px;
    bottom: 30px;
    right: 3%;
  }
`;

const Chip3 = styled(Chip)`
  top: -30px;
  left: 35%;
  display: none !important;
  animation-delay: 1.5s;

  .width_700 & {
    display: block !important;
    top: -30px;
    left: 35%;
  }
`;

const Chip4 = styled(Chip)`
  top: 20px;
  right: 5%;
  width: 50px;
  animation-delay: 1s;

  .width_700 & {
    width: 70px;
    top: -30px;
    right: 10%;
  }
`;

const Chip5 = styled(Chip)`
  top: 55%;
  left: 5%;
  width: 50px;
  animation-delay: 1s;

  .width_700 & {
    width: 70px;
    top: 35%;
    left: 5%;
  }
`;

const Chip6 = styled(Chip)`
  top: -40px;
  right: 40%;
  animation-delay: 1s;

  .width_700 & {
    top: -20px;
    right: 30%;
  }
`;

const Chip7 = styled(Chip)`
  top: 55%;
  right: 0%;
  width: 40px;
  animation-delay: 0.5s;

  .width_700 & {
    width: 95px;
    top: 35%;
    right: 5%;
  }
`;

const Chip8 = styled(Chip)`
  bottom: 20px;
  left: 3%;
  width: 90px;

  animation-delay: 0.5s;

  .width_700 & {
    width: 117px;
    bottom: 20px;
    left: 10%;
  }
`;
