import React from 'react';
import styled from 'styled-components';
import { Box, Flex, Text } from 'components';

import BlueCoin1 from 'assets/imgs/coins/chip-1.png';
import BlueCoin2 from 'assets/imgs/coins/chip-2.png';
import BlueCoin3 from 'assets/imgs/coins/chip-3.png';
import SolCoin1 from 'assets/imgs/coins/small.png';
import SolCoin2 from 'assets/imgs/coins/default.png';
import DiamondCard from 'assets/imgs/comingsoon/diamond_card.png';
import HeartCard from 'assets/imgs/comingsoon/heart_card.png';

const PageWrapper = styled(Flex)`
  position: relative;
  width: 100%;
  height: 100%;
  /* overflow-x: visible; */
  /* overflow-y: hidden; */
  overflow: hidden;
`;

const ComingSoonWrapper = styled(Flex)`
  min-height: calc(100vh - 165px);
`;

const StyledHeartCard = styled.img`
  position: absolute;
  top: -200px;
  left: 50px;
`;

const StyledDiamondCard = styled.img`
  position: absolute;
  bottom: -200px;
  right: 165px;
`;

const StyledBlueCoin1 = styled.img`
  position: absolute;
  left: 100px;
  bottom: -130px;
`;

const StyledBlueCoin2 = styled.img`
  position: absolute;
  right: 100px;
  top: -80px;
`;

const StyledBlueCoin3 = styled.img`
  position: absolute;
  top: -200px;
  left: 130px;
`;

const StyledSolCoin1 = styled.img`
  position: absolute;
  right: 5px;
  bottom: -100px;
`;

const StyledSolCoin2 = styled.img`
  position: absolute;
  bottom: -100px;
  left: -80px;
`;

export default function ComingSoon() {
  return (
    <PageWrapper
      justifyContent="center"
      alignItems="center"
      position="relative"
    >
      <ComingSoonWrapper
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        position="relative"
        pb="60px"
        px="60px"
      >
        <Box position={'relative'}>
          <Text
            fontFamily="Termina"
            fontSize="44px"
            fontWeight={600}
            lineHeight="53px"
            letterSpacing="37px"
            color="#fff"
            mt="32px"
            mb="14px"
            textAlign={'center'}
          >
            COMING SOON
          </Text>
          <Text
            fontSize="16px"
            fontWeight={400}
            color="#96A8C2"
            textAlign={'center'}
          >
            PVP Gambling Arena powered by Solana blockchain coming soon.
          </Text>
          <StyledHeartCard src={HeartCard} />
          <StyledDiamondCard src={DiamondCard} />
          <StyledBlueCoin1 src={BlueCoin1} />
          <StyledBlueCoin2 src={BlueCoin2} />
          <StyledBlueCoin3 src={BlueCoin3} />
          <StyledSolCoin1 src={SolCoin1} />
          <StyledSolCoin2 src={SolCoin2} />
        </Box>
      </ComingSoonWrapper>
    </PageWrapper>
  );
}
