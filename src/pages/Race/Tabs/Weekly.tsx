import React, { useEffect, useState } from 'react';
import useSWR, { mutate } from 'swr';
import styled from 'styled-components';
import FlipClockCountdown from '@leenguyen/react-flip-clock-countdown';
import {
  Box,
  Text,
  Span,
  Flex,
  CoinIcon,
  Avatar,
  Button,
  Grid
} from 'components';
import { InputBox } from 'components/InputBox';
import Logo from 'components/Icon/Logo';
import '@leenguyen/react-flip-clock-countdown/dist/index.css';
import weeklyImg from 'assets/imgs/race/weekly.jpg';
import { useMatchBreakpoints } from 'hooks';
import { api } from 'services';
import { StyledTable } from '../components/styles';
import RaffleTicket from '../components/RaffleTicket';
import { formatNumber, formatUserName } from 'utils/format';
import { ClipLoader } from 'react-spinners';
import { convertBalanceToChip } from 'utils/balance';

export interface WeeklyInfo {
  me: {
    avatar: string;
    id: number;
    name: string;
    rank: number;
    ticketCount: number;
    tickets: {
      date: string;
      ticketId: number;
    }[];
  };
  players: {
    avatar: string;
    id: number;
    name: string;
    rank: number;
    ticketCount: number;
    tickets: {
      date: string;
      ticketId: number;
    }[];
  }[];
  chipsPerTicket: number;
  totalPrize: number;
  totalTickets: number;
  remaining: number;
  status: 'running' | 'pending';
}

export default function Weekly() {
  const { isMobile } = useMatchBreakpoints();
  const [searchTxt, setSearchTxt] = useState('');
  const [count, setCount] = useState(5);
  const [raffleShowCount, setRaffleShowCount] = useState(isMobile ? 8 : 12);
  const {
    data: status,
    error,
    isValidating
  } = useSWR(`/weekly-raffle/status`, async arg =>
    api.get<WeeklyInfo>(arg + `?count=${count}`).then(res => res.data)
  );

  console.log(status);

  const handleChange = (e: any) => {
    setSearchTxt(e.target.value);
  };

  const handleShowMore = () => {
    const pageCount = 5;

    setCount(count + pageCount);
  };

  const handleShowMoreTicket = () => {
    let raffleCount = isMobile ? 4 : 6;
    raffleCount *= 3;
    setRaffleShowCount(raffleShowCount + raffleCount);
  };

  useEffect(() => {
    mutate('/weekly-raffle/status');
  }, [count]);

  return (
    <div className="container">
      <div className="box">
        <StyledContainer>
          <FirstBox>
            <Logo width={77} height={18} />

            <ResponsiveWeekText
              fontFamily="Termina"
              fontWeight={800}
              fontSize="43px"
              color="white"
              mt={'-5px'}
            >
              WEEKLY
              <Span fontFamily="Termina" fontWeight={800} color="#4FFF8B">
                RAFFLE
              </Span>
            </ResponsiveWeekText>

            <ResponsiveText
              fontSize={'15px'}
              fontWeight={500}
              color="#D6D6D6"
              maxWidth={'400px'}
              mt="10px"
            >
              Wager to collect raffle tickets, where you dont have to be a whale
              to win. With every &nbsp;
              <CoinIcon size={10} />
              &nbsp;{convertBalanceToChip(
                status ? status.chipsPerTicket : 0
              )}{' '}
              wagered you earn one raffle ticket. Winners drawn on Discord every
              Sunday at 12:00 am UTC.
            </ResponsiveText>

            {status && (
              <ResponsiveText
                fontWeight={500}
                color="white"
                maxWidth={'385px'}
                mt="15px"
              >
                You have earned {status.me.ticketCount} raffle tickets.
              </ResponsiveText>
            )}
          </FirstBox>

          <div>
            {status && status.status !== 'pending' && (
              <Flex
                justifyContent={'center'}
                alignItems="center"
                gap={20}
                mb="18px"
              >
                <CoinIcon size={30} />

                <StyledPrizeText
                  fontFamily={'Termina'}
                  color="white"
                  fontWeight={800}
                >
                  {convertBalanceToChip(status ? status.totalPrize : 0) / 1000}K
                  PRIZE
                </StyledPrizeText>
              </Flex>
            )}

            <Box
              background={'#0D151E99'}
              borderRadius="12px"
              padding={['6px 10px', '6px 10px', '6px 10px', '12px 20px']}
            >
              <Text
                fontFamily="Termina"
                fontWeight={600}
                fontSize={[10, 10, 10, 14]}
                textAlign="center"
                color="#D7D7D7"
                textTransform="uppercase"
                lineHeight={'25px'}
                mb="15px"
              >
                {status && status.status === 'pending'
                  ? 'WEEKLY Raffle Starts In'
                  : 'WEEKLY Raffle Ends In'}
              </Text>
              <FlipClockCountdown
                to={
                  new Date(
                    new Date().getTime() +
                      (status ? status.remaining * 1000 : 0)
                  )
                }
                labels={['Days', 'Hours', 'Minutes', 'Seconds']}
                labelStyle={{
                  fontStyle: 'Termina',
                  fontSize: 10,
                  fontWeight: 600,
                  color: '#D6D6D6',
                  paddingTop: '5px'
                }}
                digitBlockStyle={{
                  width: isMobile ? 25 : 35,
                  height: isMobile ? 45 : 55,
                  fontSize: isMobile ? 30 : 40,
                  fontWeight: 700,
                  background: '#D9D9D9',
                  color: '#333333'
                }}
                dividerStyle={{ color: 'black', height: 1 }}
                separatorStyle={{ color: 'white', size: '8px' }}
              />
            </Box>
          </div>
        </StyledContainer>

        {status && status.me.ticketCount !== 0 && (
          <Box mt="50px">
            <Flex alignItems="center" justifyContent="space-between" mb="25px">
              <Text fontSize="25px" fontWeight={500} color="white">
                My Raffle Tickets
              </Text>

              <InputBox gap={20} p="8px 10px">
                <input
                  type={'text'}
                  name="search"
                  value={searchTxt}
                  onChange={handleChange}
                  placeholder="Search Tickets"
                />
              </InputBox>
            </Flex>

            <Box
              maxHeight={isMobile ? '500px' : '300px'}
              overflowY={'auto'}
              background="#121A25"
              borderRadius="13px"
              p="12px 15px"
            >
              <Grid
                gridTemplateColumns={[
                  'repeat(2, minmax(130px, 1fr))',
                  'repeat(2, minmax(130px, 1fr))',
                  'repeat(2, minmax(130px, 1fr))',
                  'repeat(6, minmax(130px, 1fr))'
                ]}
                flexWrap="wrap"
                gap={10}
              >
                {status?.me.tickets
                  .filter(
                    ticket =>
                      ticket.ticketId
                        .toString()
                        .padStart(4, '0')
                        .indexOf(searchTxt) >= 0
                  )
                  .slice(0, raffleShowCount)
                  .map(ticket => (
                    <RaffleTicket
                      key={ticket.ticketId}
                      date={ticket.date}
                      ticketId={ticket.ticketId}
                    />
                  ))}
              </Grid>

              {raffleShowCount <= status.me.tickets.length && (
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
                  onClick={handleShowMoreTicket}
                >
                  SHOW MORE
                </Button>
              )}
            </Box>
          </Box>
        )}

        <Box mt="50px">
          <Text fontSize="25px" fontWeight={500} color="white" mb="25px">
            Weekly Raffle Leaderboard
          </Text>

          <Box
            background="#121A25"
            p="12px 15px"
            borderRadius="13px"
            overflowX={'auto'}
          >
            <StyledTable>
              <thead>
                <tr>
                  <th align="left">Rank</th>
                  <th align="left">DUELER</th>

                  <th>RAFFLE TICKETS</th>
                  <th align="right">Chance</th>
                </tr>
              </thead>
              <tbody>
                {status && status.me.rank !== 0 && (
                  <tr className="me">
                    <td>#{status.me.rank}</td>
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
                        <Span ml="16px">{formatUserName(status.me.name)}</Span>
                      </Flex>
                    </td>
                    <td align="center">{status.me.ticketCount}</td>
                    <td align="right">
                      {formatNumber(
                        (status.me.ticketCount / status.totalTickets) * 100
                      )}
                      %
                    </td>
                  </tr>
                )}
                {status &&
                  status.players.map(player => (
                    <tr key={player.id}>
                      <td>#{player.rank}</td>
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
                      <td align="center">{player.ticketCount}</td>
                      <td align="right">
                        {formatNumber(
                          (player.ticketCount / status.totalTickets) * 100
                        )}
                        %
                      </td>
                    </tr>
                  ))}
              </tbody>
            </StyledTable>

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
          </Box>
        </Box>
      </div>
    </div>
  );
}

const StyledContainer = styled(Box)`
  height: 480px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  justify-content: space-between;
  align-items: center;

  position: relative;
  background-image: url(${weeklyImg});
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
    background-size: 100% 100%;
    text-align: inherit;
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

const ResponsiveWeekText = styled(Text)`
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

const StyledPrizeText = styled(Text)`
  font-size: 35px;

  .width_1100 & {
    font-size: 48px;
  }
`;
