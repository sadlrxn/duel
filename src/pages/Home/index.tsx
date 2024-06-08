import React, { useEffect } from 'react';
import styled from 'styled-components';

import { Box, Flex, Grid } from 'components/Box';
import { Text } from 'components';
import Carousel from './components/Carousel';
import GameBox from './components/GameBox';
// import { Options } from "react-select";

import CoinflipImg from 'assets/imgs/home/coinflip.png';
import JackpotImg from 'assets/imgs/home/jackpot.png';
import DreamTowerImg from 'assets/imgs/home/dreamtower.png';
import CrashImg from 'assets/imgs/home/crash.png';
import Coin from 'components/Icon/Coin';
import { api } from 'services';
import { useAppDispatch, useAppSelector } from 'state';
import { loadStatistics } from 'state/main/actions';
import { formatNumber } from 'utils/format';
import { convertBalanceToChip } from 'utils/balance';
// import TopItem, { UserInfo } from "./components/TopItem";

// const TOP_WIN_OPTIONS: Options<{ label: string; value: string }> = [
//   { label: "All Games", value: "all" },
//   { label: "Blackjack", value: "blackjack" },
//   { label: "Coinflip", value: "coinflip" },
//   { label: "Jackpot", value: "jackpot" },
// ];

// const TOP_PLAYERS_OPTIONS: Options<{ label: string; value: string }> = [
//   { label: "Last 24 hours", value: "last24hours" },
//   { label: "Last 7 days", value: "last7days" },
//   { label: "All Time", value: "all" },
// ];

// const win_data: UserInfo[] = [
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       earned: "145,232",
//     },
//   },
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       earned: "145,232",
//     },
//   },
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       earned: "145,232",
//     },
//   },
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       earned: "145,232",
//     },
//   },
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       earned: "145,232",
//     },
//   },
// ];

// const player_data: UserInfo[] = [
//   {
//     user: {
//       avatar: "asdf",
//       name: "eroist",
//     },
//     info: {
//       exp: "145,232",
//       tier: "Dueler V",
//     },
//   },
// ];

export default function Home() {
  const statistics = useAppSelector(state => state.main);
  const dispatch = useAppDispatch();
  useEffect(() => {
    async function loadHandler() {
      try {
        const { data } = await api.get('/statistics');
        dispatch(loadStatistics(data));
      } catch (error) {}
    }
    loadHandler();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <Container padding={['45px 12px', '45px 12px', '45px 12px', '45px 25px']}>
      <Carousel />
      <Box>
        <Text
          textTransform="uppercase"
          color={'#B4BCCA'}
          fontSize="17px"
          fontWeight={700}
          mt="40px"
        >
          Duel Originals
        </Text>
        <GameBoxContainer>
          <GameBox img={CoinflipImg} title="coin flip" to="/coinflip" />
          <GameBox img={JackpotImg} title="jackpot" to="/jackpot" />
          <GameBox img={DreamTowerImg} title="Dream Tower" to="/dream-tower" />
          <GameBox img={CrashImg} title="Crash" to="/crash" />
        </GameBoxContainer>

        <Text
          textTransform="uppercase"
          color={'#B4BCCA'}
          fontSize="17px"
          fontWeight={700}
          mt="40px"
        >
          Statistics
        </Text>

        <GameBoxContainer>
          <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
            <Text
              textTransform="uppercase"
              color={'#96A8C2'}
              fontSize="12px"
              fontWeight={600}
            >
              total bets placed
            </Text>

            <Text
              textTransform="uppercase"
              color={'#fff'}
              fontSize="20px"
              fontWeight={700}
              mt="10px"
            >
              {formatNumber(statistics.totalBets)}
            </Text>
          </Box>

          <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
            <Text
              textTransform="uppercase"
              color={'#96A8C2'}
              fontSize="12px"
              fontWeight={600}
            >
              totaL amount wagered
            </Text>
            <Flex alignItems={'center'} gap={5} mt="10px">
              <Coin />
              <Text
                textTransform="uppercase"
                color={'#fff'}
                fontSize="20px"
                fontWeight={700}
              >
                {formatNumber(convertBalanceToChip(statistics.totalWagered))}
              </Text>
            </Flex>
          </Box>

          <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
            <Text
              textTransform="uppercase"
              color={'#96A8C2'}
              fontSize="12px"
              fontWeight={600}
            >
              total amount won
            </Text>

            <Flex alignItems={'center'} gap={5} mt="10px">
              <Coin />
              <Text
                textTransform="uppercase"
                color={'#fff'}
                fontSize="20px"
                fontWeight={700}
              >
                {formatNumber(convertBalanceToChip(statistics.totalProfit))}
              </Text>
            </Flex>
          </Box>
        </GameBoxContainer>

        {/* <TopItemContainer>
          <TopItem
            title="Top Win"
            options={TOP_WIN_OPTIONS}
            isWinComp={true}
            data={win_data}
          />

          <TopItem
            title="Top Players"
            options={TOP_PLAYERS_OPTIONS}
            isWinComp={false}
            data={player_data}
          />
        </TopItemContainer> */}
      </Box>
    </Container>
  );
}

const Container = styled(Box)`
  max-width: 1002px;
  /* padding-top: 15px; */
  margin: 0px auto;
`;

const GameBoxContainer = styled(Grid)`
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  align-items: end;
  margin-top: 22px;
  gap: 22px;
`;

// const TopItemContainer = styled(Grid)`
//   grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
//   align-items: end;
//   margin-top: 55px;
//   gap: 22px;
//   grid-row-gap: 30px;
// `;
