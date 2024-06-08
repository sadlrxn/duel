import useSWR from 'swr';
import { JackpotFairData } from 'api/types/jackpot';
import { CoinflipFairData } from 'api/types/coinflip';
import { DreamtowerFairData } from 'api/types/dreamtower';
import api from 'utils/api';
import { CrashFairData } from 'api/types/crash';

type FairData =
  | CoinflipFairData
  | JackpotFairData
  | DreamtowerFairData
  | CrashFairData;

type DuelGame = 'coinflip' | 'jackpot' | 'dreamtower' | 'crash';
type RoundId = number | string;

const fetchRoundInfo = async (game: DuelGame, roundId?: RoundId) => {
  if (!roundId) return undefined;

  // TODO: Check roundId type

  const { data } = await api.get<FairData>(`/${game}/round-data`, {
    params: { roundId }
  });
  return data;
};

export default function useFetchRoundInfo(game: DuelGame, roundId?: RoundId) {
  return useSWR(
    `${
      game === 'coinflip'
        ? 'Coin Flip'
        : game === 'jackpot'
        ? 'Jackpot'
        : game === 'dreamtower'
        ? 'Dream Tower'
        : 'Crash'
    } Round Info: id:${roundId}`,
    () => fetchRoundInfo(game, roundId)
  );
}
