import React from 'react';
import styled from 'styled-components';
import { Flex, Text, Button } from 'components';
// import { ReactComponent as TwitterIcon } from 'assets/imgs/icons/twitter.svg';
import { StyledSlideD } from './styles';
import { Link } from 'react-router-dom';
import rocketImg from 'assets/imgs/home/carousel/rocket.png';
// import chip1Img from 'assets/imgs/home/carousel/chips/gold-chip-1.png';
// import chip2Img from 'assets/imgs/home/carousel/chips/gold-chip-2.png';
// import chip3Img from 'assets/imgs/home/carousel/chips/gold-chip-3.png';
// import chip4Img from 'assets/imgs/home/carousel/chips/gold-chip-4.png';
// import chip5Img from 'assets/imgs/home/carousel/chips/gold-chip-5.png';
// import chip6Img from 'assets/imgs/home/carousel/chips/gold-chip-6.png';
// import chip7Img from 'assets/imgs/home/carousel/chips/gold-chip-7.png';
// import chip8Img from 'assets/imgs/home/carousel/chips/gold-chip-8.png';
// import chip9Img from 'assets/imgs/home/carousel/chips/gold-chip-9.png';

// import { Floating } from 'utils/animationToolkit';

export default function SlideD() {
  return (
    <StyledSlideD>
      <Intro>
        <Hero1>WE HAVE LIFT off!</Hero1>
        <Hero2>CRASH IS LIVE</Hero2>
      </Intro>

      <Description>
        Take yourself to new heights with the Duel Rocket. Prove you have what
        it takes to reach the stars. How high can you fly?
        <br />
        <br />
        One small step for man, one giant leap for degens...
      </Description>

      <Link to={'/crash'}>
        <StyledBtn
          p="10px 30px"
          borderRadius={'5px'}
          fontSize="14px"
          fontWeight={700}
          color="black"
          mt="30px"
        >
          PLAY NOW
        </StyledBtn>
      </Link>
      <Rocket src={rocketImg} alt="rocket" />
      {/* <Chip1 src={chip1Img} alt="chip-1" />
      <Chip2 src={chip2Img} alt="chip-2" />
      <Chip3 src={chip3Img} alt="chip-3" />
      <Chip4 src={chip4Img} alt="chip-4" />
      <Chip5 src={chip5Img} alt="chip-5" />
      <Chip6 src={chip6Img} alt="chip-6" />
      <Chip7 src={chip7Img} alt="chip-7" />
      <Chip8 src={chip8Img} alt="chip-8" />
      <Chip9 src={chip9Img} alt="chip-9" /> */}
    </StyledSlideD>
  );
}

const Intro = styled(Flex)`
  flex-direction: column;
  align-items: center;
  z-index: 10;

  .width_700 & {
    align-items: start;
  }
`;

const Hero1 = styled(Text)`
  color: #fff;
  font-family: 'Termina';
  font-weight: 800;
  font-size: 14px;
  letter-spacing: -0.02em;
  text-transform: uppercase;
  line-height: 17px;
  .width_700 & {
    font-size: 24px;
    line-height: 29px;
  }
`;

const Hero2 = styled(Text)`
  color: #4fff8b;
  font-family: 'Termina';
  font-weight: 800;
  font-size: 25px;
  line-height: 30px;
  letter-spacing: -0.02em;

  .width_700 & {
    font-size: 35px;
    line-height: 30px;
  }
`;

const Description = styled(Text)`
  color: #d6d6d6;
  z-index: 10;
  font-weight: 500;
  letter-spacing: -0.02em;
  font-size: 15px;
  line-height: 22px;
  text-align: center;
  .width_700 & {
    text-align: start;
    width: 430px;
    font-size: 15px;
  }
`;

const StyledBtn = styled(Button)`
  position: relative;
  z-index: 50;
`;

const Rocket = styled.img`
  position: absolute;
  display: none !important;
  .width_700 & {
    display: block !important;
    width: auto;
    bottom: 0px;
    right: 5%;
  }
`;
