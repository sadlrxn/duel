import { useMemo, useCallback, useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { TabList, Tab, TabPanel } from 'react-tabs';
import gsap from 'gsap';

import {
  Box,
  Flex,
  Button,
  useModal,
  Snow4,
  Snow5,
  Snow6,
  SnowMan,
  P,
  CoinIcon,
  ChipIcon as PokerChipIcon
} from 'components';
import { toast } from 'utils/toast';
import state, { useAppSelector } from 'state';
import { initialGameData } from 'state/jackpot/reducer';
import { setRoom, setRequest, setAutoBet } from 'state/jackpot/actions';
import { updateBalance } from 'state/user/actions';
import { sendMessage } from 'state/socket';

import { ChipIcon } from 'components/Chip';
import {
  JackpotRoom,
  Players,
  UserStatus,
  BetCashModal,
  BetNFTModal,
  ShowNFTModal,
  AutoBetModal
} from './components';
import { useFetchRoundInfo, useMatchBreakpoints, useQuery } from 'hooks';
import {
  JackpotFairData,
  TGameStatus,
  JackpotRoom as JackpotRoomType,
  jackpotRooms
} from 'api/types/jackpot';
import { createOldRoundInfo, getJackpotProgress } from './utils';
import {
  StyledTabs,
  Container,
  ChipIconContainer,
  TotalChip,
  TimeLine,
  HistoryButtonContainer,
  HistoryButton,
  BetButtonContainer
} from './styles';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

export interface JackpotGameProps {}

export default function JackpotGame(_: JackpotGameProps) {
  const prices = [1, 2, 5, 10, 25, 50, 100, 250, 500];
  const query = useQuery();
  const navigate = useNavigate();
  const user = useAppSelector(state => state.user);
  const room = useAppSelector(state => state.jackpot.room);
  const { game: currentGame, autoBet } = useAppSelector(
    state => state.jackpot[state.jackpot.room]
  );
  const jackpot = useAppSelector(state => state.jackpot);
  const meta = useAppSelector(state => state.meta.jackpot);
  const [roundId, setRoundId] = useState<number | undefined>(undefined);

  useEffect(() => {
    const roundId = query.get('roundId');
    if (!roundId) setRoundId(undefined);
    else setRoundId(+roundId);
  }, [query]);

  const { isHoliday } = user;
  const { data: oldGame } = useFetchRoundInfo('jackpot', roundId);
  const { isMobile } = useMatchBreakpoints();

  const [nftValue, setNftValue] = useState({ nfts: [], name: '', level: 0 });
  const [isHistory, setIsHistory] = useState(false);
  const [status, setStatus] = useState<TGameStatus>('rollend');
  const [time, setTime] = useState(Date.now());

  const historyIndex = useMemo(() => {
    if (!roundId) return -1;
    const index = jackpot.history.games.findIndex(
      history => history.roundId === roundId
    );
    if (index === -1) return -2;
    return index;
  }, [roundId, jackpot.history.games]);

  const game = useMemo(() => {
    return roundId
      ? oldGame
        ? createOldRoundInfo(oldGame as JackpotFairData)
        : initialGameData
      : currentGame;
  }, [currentGame, oldGame, roundId]);

  useEffect(() => {
    if (roundId) {
      setIsHistory(true);
      setStatus('rollend');
      setTime(Date.now() - 15 * 1000);
    } else {
      setIsHistory(false);
    }
  }, [roundId]);

  const players = useMemo(() => game?.players || [], [game]);

  const [win, winnerId] = useMemo(() => {
    if (!game) return [false, 0];
    const win =
      user.name !== '' &&
      game.winner.name === user.name &&
      (isHistory
        ? game.status === 'rolling' || game.status === 'rollend'
        : game.status === 'rollend');
    const winnerId = isHistory
      ? game.status === 'rolling' || game.status === 'rollend'
        ? game.winner.id
        : 0
      : game.status === 'rollend'
      ? game.winner.id
      : 0;
    return [win, winnerId];
  }, [game, isHistory, user.name]);

  const [userData, minBetAmount, maxBetAmount, betCount] = useMemo(() => {
    const data = players.find(player => player.id === user.id);
    const usd = data?.usdAmount ? data.usdAmount : 0;
    const nft = data?.nftAmount ? data.nftAmount : 0;
    const total = usd + nft;

    const betAmount = convertBalanceToChip(total);
    const minBetAmount = Math.max(
      convertBalanceToChip(meta[room].minBetAmount) - betAmount,
      0
    );
    const maxBetAmount = Math.min(
      convertBalanceToChip(meta[room].maxBetAmount) - betAmount,
      convertBalanceToChip(user.balance)
    );
    const betCount = data?.count ?? 0;
    const userData = {
      user: {
        id: user.id,
        name: user.name,
        level: 2,
        percent:
          ((usd + nft) / (game?.totalBetAmount ? game!.totalBetAmount : 100)) *
          100,
        avatar: user.avatar
      },
      nfts: data?.nfts ?? [],
      nftsToShow: 5,
      amount: {
        usd: convertBalanceToChip(usd),
        nft: convertBalanceToChip(nft),
        total: convertBalanceToChip(total)
      }
    };
    return [userData, minBetAmount, maxBetAmount, betCount];
  }, [game, meta, players, room, user]);

  const [onBetCash] = useModal(
    <BetCashModal
      userData={userData}
      balance={user.balance}
      status={game?.status || 'available'}
    />,
    true,
    true,
    true,
    'JackpotBetCashModal'
  );
  const [onBetNFT] = useModal(
    <BetNFTModal userData={userData} status={game?.status || 'available'} />,
    true,
    true,
    true,
    'JackpotBetNftModal'
  );
  const [onShowNFT] = useModal(
    <ShowNFTModal {...nftValue} />,
    true,
    true,
    true,
    'JackpotShowNftModal'
  );

  const [onAutoBet] = useModal(
    <AutoBetModal />,
    true,
    true,
    true,
    'JackpotAutoBetModal'
  );

  const handleBet = useCallback(
    (price: number) => {
      const amount = convertChipToBalance(price);
      state.dispatch(updateBalance({ type: -1, usdAmount: amount }));
      const content = JSON.stringify({
        amount
      });
      state.dispatch(
        sendMessage({
          type: 'event',
          room: 'jackpot',
          level: room,
          content
        })
      );
      state.dispatch(setRequest({ room, request: true }));
    },
    [room]
  );

  const handleBetCash = useCallback(() => {
    if (user.name === '') toast.info('Please sign');
    else if (
      game?.status &&
      game?.status !== 'rolling' &&
      game?.status !== 'rollend'
    )
      onBetCash();
    else toast.warning('Please wait for the next round.');
  }, [onBetCash, user, game?.status]);

  const handleBetNFT = useCallback(() => {
    if (user.name === '') toast.info('Please sign');
    else if (
      game?.status &&
      game?.status !== 'rolling' &&
      game?.status !== 'rollend'
    )
      onBetNFT();
    else toast.warning('Please wait for the next round.');
  }, [onBetNFT, user, game?.status]);

  const handleShowNFT = useCallback(
    (props: any) => {
      setNftValue(props);
      onShowNFT();
    },
    [onShowNFT]
  );

  const handleAutoBet = useCallback(() => {
    if (user.name === '') toast.info('Please sign');
    else if (autoBet) state.dispatch(setAutoBet({ room }));
    else onAutoBet();
  }, [onAutoBet, autoBet, room, user.name]);

  const handleGoToCurrentGame = useCallback(() => {
    navigate('/jackpot');
  }, [navigate]);

  const handleReplay = useCallback(() => {
    if (!game) return;
    setStatus('rolling');
    setTime(Date.now());
  }, [game]);

  useEffect(() => {
    let timeout: NodeJS.Timeout;
    if (status === 'rolling') {
      timeout = setTimeout(() => {
        setStatus('rollend');
        setTime(Date.now() - 15 * 1000);
      }, game.rolltime * 1000);
    }
    return () => {
      clearTimeout(timeout);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  const tabRef = useRef<any>(null);

  const setProgress = useCallback(() => {
    const result = Object.keys(jackpotRooms).map(room => {
      return getJackpotProgress({
        //@ts-ignore
        countingTime: jackpot[room].game.countingTime ?? 40,
        //@ts-ignore
        updatedTime: jackpot[room].game.time,
        //@ts-ignore
        status: jackpot[room].game.status,
        //@ts-ignore
        rollingTime: meta[room].rollingTime,
        //@ts-ignore
        winnerTime: meta[room].winnerTime
      });
    });
    if (!tabRef) return;
    const q = gsap.utils.selector(tabRef);
    gsap.fromTo(
      q('.jackpot_timeline'),
      {
        background: i => (result[i].roll ? '#FFE87F' : '#4FFF8B'),
        width: i =>
          `${((result[i].max - result[i].count) / result[i].max) * 100}%`
      },
      {
        width: i =>
          `${((result[i].max - result[i].count + 1) / result[i].max) * 100}%`,
        ease: 'none',
        duration: 1
      }
    );
  }, [jackpot, meta]);

  useEffect(() => {
    setProgress();
    const interval = setInterval(() => {
      setProgress();
    }, 1000);
    return () => clearInterval(interval);
  }, [setProgress]);

  return (
    <Container>
      <Players
        players={players}
        handleShowNFT={handleShowNFT}
        winnerId={winnerId}
      />
      <Box width="100%" position="relative">
        {game ? (
          <>
            <Flex justifyContent="space-between" ref={tabRef}>
              <StyledTabs
                selectedIndex={Object.keys(jackpotRooms).findIndex(
                  key => key === room
                )}
                onSelect={tabIndex => {
                  state.dispatch(
                    setRoom(
                      Object.keys(jackpotRooms)[
                        tabIndex
                      ] as unknown as JackpotRoomType
                    )
                  );
                }}
                isgameend={roundId !== undefined || winnerId !== 0}
              >
                <TabList>
                  <Flex>
                    {Object.keys(jackpotRooms).map(room => {
                      return (
                        <Tab key={room}>
                          {room.toUpperCase()}
                          <P
                            as={Flex}
                            mt="3px"
                            alignItems="center"
                            fontSize="10px"
                            fontWeight={500}
                            color="white"
                            lineHeight="12px"
                            gap={3}
                            opacity={0.35}
                          >
                            <CoinIcon size={9} />
                            {
                              //@ts-ignore
                              convertBalanceToChip(meta[room].minBetAmount) +
                                `-${
                                  room === 'wild'
                                    ? '∞'
                                    : convertBalanceToChip(
                                        //@ts-ignore
                                        meta[room].maxBetAmount
                                      )
                                }`
                            }
                          </P>
                          <b />
                          {jackpot[room as unknown as JackpotRoomType].game
                            .totalBetAmount > 0 && (
                            <TotalChip>
                              <ChipIcon $size={8.6} />
                              {convertBalanceToChip(
                                //@ts-ignore
                                jackpot[room].game.totalBetAmount
                              ).toFixed(2)}
                            </TotalChip>
                          )}

                          <TimeLine>
                            <Box
                              className="jackpot_timeline"
                              height={
                                jackpot[room as unknown as JackpotRoomType].game
                                  .totalBetAmount > 0
                                  ? '100%'
                                  : '0px'
                              }
                            />
                          </TimeLine>
                        </Tab>
                      );
                    })}
                  </Flex>
                </TabList>
                {Object.keys(jackpotRooms).map(room => {
                  return <TabPanel key={room} />;
                })}
              </StyledTabs>
              <HistoryButtonContainer
                gap={10}
                fontSize="12px"
                fontWeight="600"
                lineHeight="18px"
              >
                <HistoryButton
                  prev
                  disabled={
                    historyIndex < -1 ||
                    jackpot.history.games.length <= historyIndex + 1
                  }
                  onClick={() => {
                    navigate(
                      `/jackpot?roundId=${
                        jackpot.history.games[historyIndex + 1].roundId
                      }`
                    );
                  }}
                >
                  {'<-'}
                </HistoryButton>
                <HistoryButton
                  disabled={historyIndex <= 0}
                  onClick={() => {
                    navigate(
                      `/jackpot?roundId=${
                        jackpot.history.games[historyIndex - 1].roundId
                      }`
                    );
                  }}
                >
                  {'->'}
                </HistoryButton>
              </HistoryButtonContainer>
            </Flex>
            {isHoliday && (
              <>
                <Snow4 position="absolute" zIndex={2} top={-25} left={-13} />
                <Flex position="absolute" zIndex={2} top={34} right={-22}>
                  <SnowMan position="absolute" top={-93} right={136} />
                  {!isMobile && <Snow5 position="relative" top={-5} mr={-23} />}
                  <Snow6 />
                </Flex>
              </>
            )}
            <JackpotRoom
              roundData={{
                ...game,
                status: isHistory ? status : game.status,
                time: isHistory ? time : game.time
              }}
              isHistory={isHistory}
              time={game.time}
              betCount={betCount}
            />
          </>
        ) : (
          <>Loading...</>
        )}
        {roundId === undefined && (
          <BetButtonContainer roundId={roundId}>
            <ChipIconContainer
              flexDirection="column"
              gap={5}
              alignItems="center"
            >
              <Flex
                gap={13}
                justifyContent="center"
                alignItems="center"
                my="10px"
                width="100%"
                flexWrap="wrap"
              >
                <Flex
                  gap={13}
                  justifyContent="center"
                  alignItems="center"
                  width="max-content"
                >
                  {prices.slice(0, 5).map(price => {
                    return (
                      <PokerChipIcon
                        price={price}
                        key={price}
                        disabled={
                          price < minBetAmount ||
                          price > maxBetAmount ||
                          betCount >= meta[room].betCountLimit
                        }
                        onClick={
                          currentGame.request
                            ? undefined
                            : () => {
                                handleBet(price);
                              }
                        }
                      />
                    );
                  })}
                </Flex>
                <Flex
                  gap={13}
                  justifyContent="center"
                  alignItems="center"
                  width="max-content"
                >
                  {prices.slice(5).map(price => {
                    return (
                      <PokerChipIcon
                        price={price}
                        key={price}
                        disabled={
                          price < minBetAmount ||
                          price > maxBetAmount ||
                          betCount >= meta[room].betCountLimit
                        }
                        onClick={
                          currentGame.request
                            ? undefined
                            : () => {
                                handleBet(price);
                              }
                        }
                      />
                    );
                  })}
                </Flex>
              </Flex>
            </ChipIconContainer>
            <Flex gap={22} width="100%">
              <Button
                variant="secondary"
                outlined
                scale="sm"
                width="100%"
                background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)"
                color="#FFF"
                borderColor="chipSecondary"
                letterSpacing="0.16em"
                onClick={handleBetCash}
                disabled={isHistory && status === 'rolling'}
              >
                CUSTOM BET
              </Button>
              <Button
                variant="secondary"
                outlined
                scale="sm"
                width="100%"
                background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
                color="#FFFFFF"
                borderColor="success"
                letterSpacing="0.16em"
                onClick={handleBetNFT}
              >
                BET NFT
              </Button>
            </Flex>
            <Button
              variant="secondary"
              outlined
              scale="sm"
              width="100%"
              background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)"
              color="#FFF"
              borderColor="chipSecondary"
              letterSpacing="0.16em"
              onClick={handleAutoBet}
            >
              {autoBet ? 'STOP AUTOBET' : 'AUTO BET'}
            </Button>
          </BetButtonContainer>
        )}
        {roundId !== undefined && (
          <Flex
            gap={22}
            width="100%"
            mt="10px"
            mb="25px"
            flexDirection="column"
            alignItems="center"
          >
            <Button
              variant="secondary"
              outlined
              scale="sm"
              width="259px"
              background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)"
              color="#FFF"
              borderColor="chipSecondary"
              letterSpacing="0.16em"
              onClick={handleReplay}
              disabled={isHistory && status === 'rolling'}
            >
              {status === 'rolling' ? 'REPLAYING...' : 'REPLAY SPIN'}
            </Button>
            <Button
              variant="secondary"
              outlined
              scale="sm"
              width="259px"
              background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
              color="#FFFFFF"
              borderColor="success"
              letterSpacing="0.16em"
              onClick={handleGoToCurrentGame}
            >
              BACK TO CURRENT GAME
            </Button>
          </Flex>
        )}
        <UserStatus
          {...userData}
          handleShowNFT={handleShowNFT}
          background={
            win
              ? 'linear-gradient(90deg, rgba(255, 226, 75, 0.2) 0%, rgba(255, 226, 75, 0) 100%)'
              : userData.amount.usd > 0
              ? 'linear-gradient(90deg, #123329 0%, #0f1a26 51.08%)'
              : 'linear-gradient(90deg, #0f1a26 0%, #0f1a26 51.08%)'
          }
        />
      </Box>
    </Container>
  );
}
