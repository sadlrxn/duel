import React from 'react';
import styled from 'styled-components';

import { Flex, Grid, Topbar } from 'components';
// import { useAppDispatch } from 'state';

import { PlinkoProvider } from 'contexts/Plinko/Provider';

import { Background, History, PlinkoGame, PlinkoTabs } from './components';

export default function Plinko() {
  return (
    <PlinkoProvider>
      <GameContainer>
        <Background />
        <Container>
          <Topbar title={'PLINKO'} />
          <History />
          <PlinkoTabs />
        </Container>
        <Flex
          width="100%"
          height="100%"
          justifyContent="center"
          alignItems="center"
        >
          <PlinkoGame />
        </Flex>
      </GameContainer>
    </PlinkoProvider>
  );
}

const GameContainer = styled(Grid)`
  gap: 20px;
  width: 100%;
  height: calc(100vh - 65px);

  .width_700 & {
    grid-template-columns: 324px 1fr;
    max-height: calc(100% - 110px);
  }
`;

const Container = styled(Flex)`
  flex-direction: column;
  gap: 18px;
  height: 100%;
`;
