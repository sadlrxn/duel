import React, { useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { Tab, TabList } from 'react-tabs';

import { Box, Flex } from 'components';

import { JackpotFairData } from 'api/types/jackpot';
import { CoinflipRoundData } from 'api/types/coinflip';
import { useMatchBreakpoints } from 'hooks';

import { Desc, Verification } from './Tabs';
import { StyledTabs, StyledTabPanel } from './styles';
import 'react-tabs/style/react-tabs.css';
import Seed from './Tabs/Seed';

const Fair: React.FC = () => {
  const { state } = useLocation();
  const { isMobile } = useMatchBreakpoints();

  const [gameType, gameData, room] = useMemo(() => {
    if (!state) return [undefined, undefined, undefined];
    const { gameType, gameData, room } = state as {
      gameType?: string;
      gameData?: JackpotFairData | CoinflipRoundData;
      room?: string;
    };
    return [gameType, gameData, room];
  }, [state]);

  const [fairContent, verifyContent, unhashContent] = useMemo(() => {
    return isMobile
      ? ['FAIR', 'VERIFY', 'SEED']
      : ['PROVABLY FAIR', 'VERIFICATION', 'CLIENT / SERVER SEED'];
  }, [isMobile]);

  return (
    <Box padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}>
      <StyledTabs defaultIndex={room === 'seed' ? 2 : gameType ? 1 : 0}>
        <TabList>
          <Flex>
            <Tab>
              {fairContent}
              <b />
            </Tab>
            <Tab>
              {verifyContent}
              <b />
            </Tab>
            <Tab>
              {unhashContent}
              <b />
            </Tab>
          </Flex>
        </TabList>
        <StyledTabPanel>
          <Desc />
        </StyledTabPanel>
        <StyledTabPanel>
          <Verification gameType={gameType} gameData={gameData} />
        </StyledTabPanel>
        <StyledTabPanel>
          <Seed />
        </StyledTabPanel>
      </StyledTabs>
    </Box>
  );
};

export default React.memo(Fair);
