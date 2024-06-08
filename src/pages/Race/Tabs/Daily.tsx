import React, { useEffect, useMemo, useState } from 'react';
import { ClipLoader } from 'react-spinners';
import styled from 'styled-components';
import useSWR, { mutate } from 'swr';
import { api } from 'services';
import Logo from 'components/Icon/Logo';
import {
  FirstMedal,
  First,
  SecondMedal,
  Second,
  ThirdMedal,
  Third
} from '../Icons';
import {
  StyledTable,
  StyledLeaderBox,
  StyledCrownIcon
} from '../components/styles';
import dailyImg from 'assets/imgs/race/daily.jpg';

import FlipClockCountdown from '@leenguyen/react-flip-clock-countdown';
import '@leenguyen/react-flip-clock-countdown/dist/index.css';
import {
  Avatar,
  Badge,
  Box,
  Button,
  Chip,
  Flex,
  Grid,
  Span,
  Text
} from 'components';
import { convertBalanceToChip } from 'utils/balance';
import { useMatchBreakpoints } from 'hooks';
import { formatUserName } from 'utils/format';

const medals = [
  <FirstMedal key={0} />,
  <SecondMedal key={1} />,
  <ThirdMedal key={2} />
];

const prizes = [<First key={0} />, <Second key={1} />, <Third key={2} />];

export interface DailyInfo {
  me: {
    avatar: string;
    id: number;
    name: string;
    rank: number;
    wagered: number;
  };
  players: {
    avatar: string;
    id: number;
    name: string;
    rank: number;
    wagered: number;
  }[];
  prizes: number[];
  remaining: number;
  status: 'running' | 'pending';
}

export default function Daily() {
  const [count, setCount] = useState(5);
  const { isMobile } = useMatchBreakpoints();
  const {
    data: status,
    error,
    isValidating
  } = useSWR(`/daily-race/status`, async arg =>
    api.get<DailyInfo>(arg + `?count=${count}`).then(res => res.data)
  );

  const handleShowMore = () => {
    const pageCount = 5;

    setCount(count + pageCount);
  };

  useEffect(() => {
    mutate('/daily-race/status');
  }, [count]);

  const winners = useMemo(() => {
    if (!status) return undefined;

    const w = [status.me, ...status.players]
      .filter(player => player.rank < 4 && player.rank > 0)
      .sort((a, b) => {
        return a.rank - b.rank;
      });

    if (w.length > 1) {
      const first = w[0];
      w.splice(0, 1);
      w.splice(1, 0, first);
    }

    return w;
  }, [status]);

  return (
    <div className="container">
      <div className="box">
        <StyledContainer>
          <FirstBox>
            <Logo width={77} height={18} />

            <ResponsiveDailyText
              fontFamily="Termina"
              fontWeight={800}
              fontSize="43px"
              color="white"
              mt={'-5px'}
            >
              DAILY
              <Span fontFamily="Termina" fontWeight={800} color="#4FFF8B">
                RACE
              </Span>
            </ResponsiveDailyText>

            <ResponsiveText
              fontWeight={500}
              color="#D6D6D6"
              maxWidth={'370px'}
              mt="15px"
            >
              Every 24 hours Duel rewards players with the highest total wagers
              that day with CHIPS. Wager more for next days race or withdraw
              your CHIPS at anytime.
            </ResponsiveText>
          </FirstBox>

          <Box background={'#0D151E99'} borderRadius="12px" padding="12px 30px">
            <Text
              fontFamily="Termina"
              fontWeight={600}
              fontSize={14}
              textAlign="center"
              color="#D7D7D7"
              textTransform="uppercase"
              lineHeight={'25px'}
              mb="20px"
            >
              {status && status.status === 'pending'
                ? 'Daily Race Starts In'
                : 'Daily Race Ends In'}
            </Text>
            <FlipClockCountdown
              to={
                new Date(
                  new Date().getTime() + (status ? status.remaining * 1000 : 0)
                )
              }
              labels={['Days', 'Hours', 'Minutes', 'Seconds']}
              labelStyle={{
                fontSize: 16,
                fontWeight: 500,
                color: '#D6D6D6',
                paddingTop: '5px'
              }}
              digitBlockStyle={{
                width: isMobile ? 30 : 40,
                height: isMobile ? 50 : 55,
                fontSize: isMobile ? 40 : 50,
                fontWeight: 700,
                background: '#D9D9D9',
                color: '#333333'
              }}
              dividerStyle={{ color: 'black', height: 1 }}
              separatorStyle={{ color: 'white', size: '8px' }}
              renderMap={[false, true, true, true]}
            />
          </Box>
        </StyledContainer>

        <Box mt="50px">
          <Text fontSize="25px" fontWeight={500} color="white" mb="25px">
            Daily Race Leaderboard
          </Text>

          <StyledLeaderBox
            background="#121A25"
            p="12px 15px"
            borderRadius="13px"
          >
            {winners && (
              <Grid
                mt={['-40px', '-40px', '-40px', '-60px']}
                mb="20px"
                alignItems={'end'}
                gridTemplateColumns={'1fr 1fr 1fr'}
              >
                {winners.length === 1 && <Box></Box>}
                {winners.map(winner => (
                  <Box key={winner.id}>
                    {winner.rank === 1 && (
                      <Flex justifyContent={'center'} mb="2px">
                        <StyledCrownIcon />
                      </Flex>
                    )}

                    <Avatar
                      userId={winner.id}
                      image={winner.avatar}
                      border={
                        winner.rank === 1
                          ? '2px solid #F6C20A'
                          : winner.rank === 2
                          ? '2px solid #C1CFD2'
                          : '2px solid #ED6939'
                      }
                      padding="0px"
                      size={
                        isMobile
                          ? winner.rank === 1
                            ? '95px'
                            : `70px`
                          : winner.rank === 1
                          ? '154px'
                          : `112px`
                      }
                      filter={
                        isMobile
                          ? winner.rank === 1
                            ? 'drop-shadow(0px 0px 30px rgba(246, 194, 10, 0.5))'
                            : winner.rank === 2
                            ? `drop-shadow(0px 0px 20px rgba(193, 207, 210, 0.5))`
                            : `drop-shadow(0px 0px 20px rgba(237, 105, 57, 0.5))`
                          : winner.rank === 1
                          ? 'drop-shadow(0px 0px 18px rgba(246, 194, 10, 0.5))'
                          : winner.rank === 2
                          ? `drop-shadow(0px 0px 12px rgba(193, 207, 210, 0.5))`
                          : `drop-shadow(0px 0px 12px rgba(237, 105, 57, 0.5))`
                      }
                    />
                    <Flex
                      position={'relative'}
                      justifyContent={'center'}
                      mt={['-10px', '-10px', '-10px', '-15px']}
                      zIndex={10}
                    >
                      {prizes[winner.rank - 1]}
                    </Flex>

                    <Text
                      mt="10px"
                      textAlign={'center'}
                      color="white"
                      fontSize={['10px', '10px', '10px', '16px']}
                      fontWeight={700}
                    >
                      {formatUserName(winner.name)}
                    </Text>

                    <Flex justifyContent={'center'}>
                      <Chip
                        fontSize={isMobile ? '10px' : '16px'}
                        price={convertBalanceToChip(winner.wagered)}
                      />
                    </Flex>
                  </Box>
                ))}
              </Grid>
            )}

            <Box overflowX={'auto'}>
              <StyledTable>
                <thead>
                  <tr>
                    <th align="left">Rank</th>
                    <th align="left">DUELER</th>

                    <th align="left">Daily Wagered</th>
                    <th align="right">Daily Prize</th>
                  </tr>
                </thead>
                <tbody>
                  {status && status.me.rank !== 0 && (
                    <tr className="me">
                      <td>
                        {status.me.rank < 4
                          ? medals[status.me.rank - 1]
                          : `#${status.me.rank}`}
                      </td>
                      <td>
                        <Flex alignItems="center">
                          <Avatar
                            userId={status.me.id}
                            image={status.me.avatar}
                            border="none"
                            borderRadius="5px"
                            padding="0px"
                            size="30px"
                          />
                          <Span ml="16px">
                            {formatUserName(status.me.name)}
                          </Span>
                        </Flex>
                      </td>
                      <td>
                        <Chip price={convertBalanceToChip(status.me.wagered)} />
                      </td>
                      <td align="right">
                        {status.me.rank <= status.prizes.length ? (
                          <Badge fontSize="14px">
                            <Chip
                              price={convertBalanceToChip(
                                status.prizes[status.me.rank - 1]
                              )}
                              color="#4FFF8B"
                            />
                          </Badge>
                        ) : (
                          <></>
                        )}
                      </td>
                    </tr>
                  )}

                  {status &&
                    status.players.map(player => (
                      <tr key={player.id}>
                        <td>
                          {player.rank < 4
                            ? medals[player.rank - 1]
                            : `#${player.rank}`}
                        </td>
                        <td>
                          <Flex alignItems="center">
                            <Avatar
                              userId={player.id}
                              image={player.avatar}
                              border="none"
                              borderRadius="5px"
                              padding="0px"
                              size="30px"
                            />
                            <Span ml="16px">{formatUserName(player.name)}</Span>
                          </Flex>
                        </td>
                        <td>
                          <Chip price={convertBalanceToChip(player.wagered)} />
                        </td>
                        <td align="right">
                          {player.rank <= status.prizes.length ? (
                            <Badge fontSize="14px">
                              <Chip
                                price={convertBalanceToChip(
                                  status.prizes[player.rank - 1]
                                )}
                                color="#4FFF8B"
                              />
                            </Badge>
                          ) : (
                            <></>
                          )}
                        </td>
                      </tr>
                    ))}
                </tbody>
              </StyledTable>
            </Box>

            {(isValidating || (status && status.players.length >= count)) && (
              <Button
                variant="secondary"
                outlined
                scale="sm"
                width={153}
                background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
                color="#FFFFFF"
                borderColor="chipSecondary"
                marginX="auto"
                marginTop={10}
                onClick={handleShowMore}
              >
                {isValidating ? (
                  <ClipLoader color="#fff" size={20} />
                ) : (
                  'SHOW MORE'
                )}
              </Button>
            )}
          </StyledLeaderBox>
        </Box>
      </div>
    </div>
  );
}

const StyledContainer = styled(Box)`
  height: 480px;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  gap: 20px;
  align-items: center;
  position: relative;
  background-image: url(${dailyImg});
  background-size: auto 100%;
  background-repeat: no-repeat;
  background-position: center top;

  padding: 20px 15px;

  border: 1px solid #45566e;
  border-radius: 12px;

  .width_1100 & {
    height: 370px;
    display: flex;
    flex-direction: row;
    text-align: inherit;
    background-size: 100% 100%;
    justify-content: space-between;
    padding: 30px 40px;
  }
`;

const FirstBox = styled(Box)`
  text-align: center;

  .width_1100 & {
    text-align: inherit;
  }
`;

const ResponsiveDailyText = styled(Text)`
  font-size: 28px;
  .width_1100 & {
    font-size: 43px;
  }
`;

const ResponsiveText = styled(Text)`
  margin: 15px auto 0px auto;
  .width_1100 & {
    margin: 15px auto 0px 0px;
  }
`;
