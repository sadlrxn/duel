import React, { useCallback } from 'react';
import styled from 'styled-components';
import { Button, Flex, Text } from 'components';
import { StyledSlideB } from './styles';
import duel1Img from 'assets/imgs/home/carousel/duel_3.png';
import duel2Img from 'assets/imgs/home/carousel/claim.png';
import { Floating } from 'utils/animationToolkit';
import { Link } from 'react-router-dom';

export default function SlideB() {
  return (
    <StyledSlideB>
      {/* <MintLink href="https://mint.duel.win/" target="_blank" rel="noreferrer"> */}
      <Intro>
        <Welcome>2223</Welcome>
        <Arena>DUELBOTS</Arena>
        {/* <Dueler>SOLD OUT</Dueler> */}
      </Intro>
      <Description>
        Attention Duelbot hodlers. <br /> Don&apos;t let your CHIPs run out -
        stake your Duelbots and watch as your CHIPs start rolling in. Get back
        80% of fees earned from Duel&apos;s V1 games.
      </Description>
      <Link to={'/staking'}>
        <StyledBtn
          color="black"
          fontWeight={700}
          fontSize="14px"
          p="10px 15px"
          borderRadius={'5px'}
          background={'sucess'}
        >
          Stake Duelbot
        </StyledBtn>
      </Link>

      <DuelbotX src={duel1Img} alt="duel-1" />
      <DuelbotY src={duel2Img} alt="duel-2" />
      {/* </MintLink> */}
    </StyledSlideB>
  );
}

const Intro = styled(Flex)`
  position: relative;
  flex-direction: column;
  align-items: center;
  z-index: 10;
`;

const Welcome = styled(Text)`
  color: white;
  font-family: 'Termina';
  font-weight: 800;
  font-size: 20px;
  line-height: 25px;
  letter-spacing: -0.02em;
  text-align: center;
  .width_700 & {
    font-size: 24.2821px;
  }
`;

const Arena = styled(Text)`
  color: white;
  font-family: 'Termina';
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
  position: relative;
  color: #d6d6d6;
  z-index: 10;
  font-weight: 500;
  letter-spacing: -0.02em;
  font-size: 16.71px;
  line-height: 26px;
  text-align: center;
  .width_700 & {
    width: 370px;
    font-size: 16px;
    text-align: center;
  }
`;

const DuelbotX = styled.img`
  position: absolute;

  width: 250px;
  z-index: 5;
  left: 50%;
  transform: translate(-50%, 0%);
  bottom: 0;
  .width_700 & {
    width: 300px;
    left: 5%;
    transform: translate(0%, 0%);
  }
`;

const DuelbotY = styled.img`
  position: absolute;
  display: none !important;
  /* width: 180px; */

  right: -5%;
  .width_700 & {
    display: block !important;
    /* width: 240px; */
    right: 2%;
    top: -8%;
  }
`;

const StyledBtn = styled(Button)`
  position: relative;
  z-index: 50;
`;
const MintLink = styled.a`
  z-index: 20;
`;
