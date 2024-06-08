import {
  JackpotFairData,
  JackpotRoundData,
  BetPlayer as Player,
  Candidate,
  TGameStatus
} from 'api/types/jackpot';
import { NFT } from 'api/types/nft';
import { convertJackpotBetServerToClient } from 'state/middleware/grandJackpotListener';
import state from 'state';

export const getTotalNfts = (players: Player[]): NFT[] => {
  let nfts: NFT[] = [];
  players.forEach(p => (nfts = nfts.concat(p.nfts ?? [])));
  return nfts;
};

export const calcTotalUsdAmount = (players: Player[]) => {
  return players.map(p => p.usdAmount).reduce((partial, a) => partial + a, 0);
};

export const calcTotalNftAmount = (players: Player[]) => {
  return players
    .map(p => p.nftAmount)
    .reduce((partial: number, a) => partial + (a ?? 0), 0);
};

export const convertPlayersToCandidate = (players: Player[]): Candidate[] => {
  const usdBetAmount = calcTotalUsdAmount(players);
  const nftBetAmount = calcTotalNftAmount(players);
  const totalBetAmount = usdBetAmount + nftBetAmount || 100;
  return players.map(p => ({
    id: p.id,
    name: p.name,
    avatar: p.avatar,
    count: Math.ceil(
      ((p.usdAmount + (p.nftAmount ?? 0)) / totalBetAmount) * 50
    ),
    percent: p.percent
  }));
};

interface RoundData extends JackpotFairData {
  endedAt?: number;
}

export const createOldRoundInfo = (
  round: RoundData,
  isGrand: boolean = false
): JackpotRoundData => {
  //@ts-ignore
  const { players: bets, endedAt } = round;
  const usdBetAmount = calcTotalUsdAmount(round.players);
  const nftBetAmount = calcTotalNftAmount(round.players);
  const rolltime = 15;

  const differ = state.getState().socket.differ;
  const time = endedAt ? new Date(endedAt).getTime() + differ : Date.now();

  let totalAmount = 0;
  let nfts: NFT[] = [];

  const players = bets
    .map(bet => {
      const player = convertJackpotBetServerToClient({ ...bet });
      if (!isGrand) player.role = undefined;
      nfts = [...nfts, ...player.nfts!];
      if (player.role !== 'admin')
        totalAmount += player.usdAmount + player.nftAmount!;
      return player;
    })
    .sort(
      (p1, p2) => p2.usdAmount + p2.nftAmount! - (p1.usdAmount + p1.nftAmount!)
    );

  players.forEach(player => {
    if (player.role === 'admin') player.percent = 0;
    else
      player.percent =
        ((player.nftAmount! + player.usdAmount) / totalAmount) * 100;
  });

  return {
    ...round,
    status: 'rollend',
    nfts: getTotalNfts(round.players),
    usdBetAmount,
    nftBetAmount,
    players,
    totalBetAmount: usdBetAmount + nftBetAmount,
    time,
    rolltime,
    animationTime: rolltime,
    candidates: convertPlayersToCandidate(
      players.filter(p => p.role !== 'admin')
    ),
    request: false
  };
};

export const getJackpotProgress = ({
  countingTime,
  updatedTime,
  status,
  rollingTime,
  winnerTime
}: {
  countingTime: number;
  updatedTime: number;
  status: TGameStatus;
  rollingTime: number;
  winnerTime: number;
}) => {
  let max = countingTime,
    count = countingTime;
  let time = (Date.now() - updatedTime) / 1000;
  let roll = false;
  if (time < 0) time = 0;
  switch (status) {
    case 'started':
      max = countingTime;
      count = Math.ceil(time);
      if (count > max) {
        count -= max;
        max = rollingTime - winnerTime;
        roll = true;
      }
      break;
    case 'rolling':
      max = rollingTime - winnerTime;
      count = Math.ceil(time);
      roll = true;
      break;
    case 'rollend':
      max = winnerTime;
      count = Math.ceil(time);
      roll = true;
      break;
  }

  count = max - count;
  if (count < 0) count = 0;
  return { max, count, roll };
};
