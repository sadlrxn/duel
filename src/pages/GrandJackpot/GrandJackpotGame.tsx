import { useMemo, useCallback, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';

import { Box, Flex, Button, useModal } from 'components';
import { useAppSelector } from 'state';
import { toast } from 'utils/toast';
import { UserStatus, ShowNFTModal } from '../Jackpot/components';
import { BetCashModal, BetNFTModal, JackpotRoom, Players } from './components';
import { TGameStatus, JackpotFairData } from 'api/types/jackpot';
import { useFetchRoundInfo, useQuery } from 'hooks';
import { createOldRoundInfo } from 'pages/Jackpot/utils';
import { initialGameData } from 'state/grandJackpot/reducer';
import { convertBalanceToChip } from 'utils/balance';

const Container = styled(Box)``;

export interface GrandJackpotGameProps {}

export default function GrandJackpotGame(_: GrandJackpotGameProps) {
  const query = useQuery();
  const navigate = useNavigate();
  const user = useAppSelector(state => state.user);
  const currentGame = useAppSelector(state => state.grandJackpot.game);
  const [roundId, setRoundId] = useState<number | undefined>(undefined);

  useEffect(() => {
    const roundId = query.get('roundId');
    if (!roundId) setRoundId(undefined);
    else setRoundId(+roundId);
  }, [query]);

  const { data: oldGame } = useFetchRoundInfo('jackpot', roundId);

  const [nftValue, setNftValue] = useState({ nfts: [], name: '', level: 0 });
  const [isHistory, setIsHistory] = useState(false);
  const [status, setStatus] = useState<TGameStatus>('rolling');
  const [time, setTime] = useState(Date.now());

  const game = useMemo(() => {
    return roundId
      ? oldGame
        ? createOldRoundInfo(oldGame as JackpotFairData, true)
        : initialGameData
      : currentGame;
  }, [currentGame, oldGame, roundId]);

  useEffect(() => {
    if (roundId) {
      setIsHistory(true);
      setStatus('rollend');
      setTime(Date.now() - (60 * 60 + 15) * 1000);
    } else {
      setIsHistory(false);
    }
  }, [roundId]);

  const players = useMemo(() => game?.players || [], [game]);

  const [win, winnerId] = useMemo(() => {
    if (!game) return [false, 0];
    const win =
      user.name &&
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

  const userData = useMemo(() => {
    const data = players.find(player => player.id === user.id);
    const usd = data?.usdAmount ? data.usdAmount : 0;
    const nft = data?.nftAmount ? data.nftAmount : 0;
    const total = usd + nft;
    return {
      user: {
        id: user.id,
        name: user.name,
        level: 2,
        percent: data ? data.percent ?? 0 : 0,
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
  }, [players, user.avatar, user.id, user.name]);

  const [onBetCash] = useModal(
    <BetCashModal
      userData={userData}
      balance={user.balance}
      status={game?.status || 'available'}
    />,
    true,
    true,
    true,
    'GrandJackpotBetCashModal'
  );
  const [onBetNFT] = useModal(
    <BetNFTModal userData={userData} status={game?.status || 'available'} />,
    true,
    true,
    true,
    'GrandJackpotBetNftModal'
  );
  const [onShowNFT] = useModal(
    <ShowNFTModal {...nftValue} />,
    true,
    true,
    true,
    'GrandJackpotShowNftModal'
  );

  const handleBetCash = useCallback(() => {
    if (user.name === '') toast.info('Please sign');
    else if (game?.status === 'started') onBetCash();
    else toast.warning('Please wait for the next round.');
  }, [onBetCash, user, game?.status]);

  const handleBetNFT = useCallback(() => {
    if (user.name === '') toast.info('Please sign');
    else if (game?.status === 'started') onBetNFT();
    else toast.warning('Please wait for the next round.');
  }, [onBetNFT, user, game?.status]);

  const handleShowNFT = useCallback(
    (props: any) => {
      setNftValue(props);
      onShowNFT();
    },
    [onShowNFT]
  );

  const handleGoToCurrentGame = useCallback(() => {
    navigate('/grandjackpot');
  }, [navigate]);

  const handleReplay = useCallback(() => {
    setStatus('rolling');
    setTime(Date.now());
  }, []);

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

  return (
    <Flex flexDirection="column" gap={36}>
      <Container>
        <Box width="100%">
          {game ? (
            <JackpotRoom
              roundData={{
                ...game,
                status: isHistory ? status : game.status,
                time: isHistory ? time : game.time
              }}
              handleShowNFT={handleShowNFT}
              isHistory={isHistory}
            />
          ) : (
            <>Loading...</>
          )}
          <Flex
            flexDirection={roundId ? 'column' : 'row'}
            alignItems="center"
            justifyContent="center"
            gap={12}
            mb={60}
          >
            {user.role !== 'admin' || roundId !== undefined ? (
              <>
                <Button
                  variant="secondary"
                  outlined
                  scale="sm"
                  width={roundId ? 259 : 153}
                  background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)"
                  color="chip"
                  borderColor="chipSecondary"
                  letterSpacing="0.16em"
                  onClick={roundId ? handleReplay : handleBetCash}
                  disabled={isHistory && status === 'rolling'}
                >
                  {roundId
                    ? status === 'rolling'
                      ? 'REPLAYING...'
                      : 'REPLAY SPIN'
                    : 'BET CHIPS'}
                </Button>
                <Button
                  variant="secondary"
                  outlined
                  scale="sm"
                  width={roundId ? 259 : 153}
                  background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
                  color="#FFFFFF"
                  borderColor="success"
                  letterSpacing="0.16em"
                  onClick={roundId ? handleGoToCurrentGame : handleBetNFT}
                >
                  {roundId ? 'BACK TO CURRENT GAME' : 'BET NFT'}
                </Button>
              </>
            ) : (
              <>
                <Button
                  variant="secondary"
                  outlined
                  scale="sm"
                  background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0) 162.5%)"
                  color="chip"
                  borderColor="#FE7EAC"
                  letterSpacing="0.16em"
                  onClick={handleBetCash}
                  px="30px"
                >
                  ADD CHIPS AS ADMIN
                </Button>
                <Button
                  variant="secondary"
                  outlined
                  scale="sm"
                  background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
                  color="#FFFFFF"
                  borderColor="#9C4FFF"
                  letterSpacing="0.16em"
                  px="30px"
                  onClick={handleBetNFT}
                >
                  ADD NFT AS ADMIN
                </Button>
              </>
            )}
          </Flex>
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
      <Players
        players={players}
        winnerId={winnerId}
        handleShowNFT={handleShowNFT}
      />
    </Flex>
  );
}
