import { JackpotFairData, Candidate } from 'api/types/jackpot';
import { CoinflipRoundData } from 'api/types/coinflip';
import { defaultTower, DreamtowerFairData } from 'api/types/dreamtower';
import { convertObjectToText } from 'utils/convert';
import {
  verifyFairness,
  generateTower,
  verifyCrashFairness
} from 'utils/verifyFairness';

import {
  coinflipSampleInput,
  jackpotSampleInput,
  dreamtowerSampleInput,
  crashSampleInput
} from './constant';

export const getCoinflipInput = (gameType?: string, gameData?: any) => {
  if (gameType !== 'coinflip') return coinflipSampleInput;

  const {
    signedString: randomString,
    headsUser,
    tailsUser,
    winnerId,
    amount: betAmount
  } = gameData as CoinflipRoundData;

  const players = [
    {
      id: headsUser!.id,
      name: headsUser!.name,
      color: 'Green'
    },
    {
      id: tailsUser!.id,
      name: tailsUser!.name,
      color: 'Purple'
    }
  ];

  return convertObjectToText(
    {
      randomString,
      betAmount,
      result: winnerId === headsUser!.id ? 'Green' : 'Purple',
      players
    },
    1
  );
};

export const getJackpotInput = (gameType?: string, gameData?: any) => {
  if (gameType !== 'jackpot') return jackpotSampleInput;

  const { signedString: randomString, players: candidates } =
    gameData as JackpotFairData;
  const players = candidates.map(player => {
    return {
      id: player.id,
      name: player.name,
      // avatar: player.avatar,
      usdAmount: player.usdAmount,
      nftAmount: player.nftAmount
    };
  });

  return convertObjectToText({ randomString, players }, 1);
};

export const getDreamTowerInput = (gameType?: string, gameData?: any) => {
  if (gameType !== 'dreamtower') return dreamtowerSampleInput;

  const { clientSeed, serverSeed, nonce, difficulty } =
    gameData as DreamtowerFairData;
  const input = {
    clientSeed,
    serverSeed: serverSeed ?? '',
    nonce,
    difficulty: difficulty.level
  };

  return convertObjectToText(input, 1);
};

export const getCrashInput = (gameType?: string, gameData?: any) => {
  if (gameType !== 'crash') return crashSampleInput;

  const { clientSeed, serverSeed } = gameData;
  const input = { serverSeed, clientSeed };

  return convertObjectToText(input, 1);
};

export const getCoinflipResult = async (data: any) => {
  const candidates: Candidate[] = data.players.map((player: any) => {
    return {
      id: 1,
      name: player.name,
      avatar: player.avatar || '',
      percent: 1,
      color: player.color
    };
  });
  const result = await verifyFairness(candidates, data.randomString);
  let index = data.players.findIndex(
    //@ts-ignore
    (player: any) => player.color === result.winner.color
  );
  return data.players[index].color;
};

export const getJackpotResult = async (data: any) => {
  // let totalBet = 0;
  const players = data.players.map((player: any) => {
    const usdAmount = player.usdAmount || 0;
    const nftAmount = player.nftAmount || 0;
    const betAmount = player.betAmount || usdAmount + nftAmount;
    // totalBet += betAmount;
    return { ...player, usdAmount, nftAmount, betAmount };
  });

  const candidates: Candidate[] = players.map((player: any) => {
    return {
      id: player.id,
      name: player.name,
      avatar: player.avatar || '',
      percent: player.betAmount
    };
  });

  const result = await verifyFairness(candidates, data.randomString);
  const winner = data.players.find(
    (player: any) => player.id === result.winner.id
  );
  return convertObjectToText(winner, 1);
};

const difficulties = {
  easy: {
    blocksInRow: 4,
    starsInRow: 3
  },
  medium: {
    blocksInRow: 3,
    starsInRow: 2
  },
  hard: {
    blocksInRow: 2,
    starsInRow: 1
  },
  expert: {
    blocksInRow: 3,
    starsInRow: 1
  },
  master: {
    blocksInRow: 4,
    starsInRow: 1
  }
};

export const getDreamTowerResult = async (data: any) => {
  const { clientSeed, serverSeed, nonce, difficulty } = data;
  //@ts-ignore
  if (!difficulties[difficulty.toLowerCase()])
    return {
      towerData: { tower: defaultTower, blocksInRow: 4 },
      result: 'Invalid Input'
    };
  const { blocksInRow, starsInRow } =
    //@ts-ignore
    difficulties[difficulty.toLowerCase()];
  const tresult = await generateTower(
    serverSeed,
    clientSeed,
    nonce,
    blocksInRow,
    starsInRow,
    9
  );
  let tower = new Array(9).fill(-1).map(() => new Array(blocksInRow).fill(0));
  let result = '[\n';
  for (let i = tresult.length - 1; i >= 0; i--) {
    let str = ' '.repeat(2) + '[';
    for (let j = 0; j < tresult[i].length; j++) {
      tower[tresult.length - i - 1][tresult[i][j]] = 1;
      str += tresult[i][j];
      if (j !== tresult[i].length - 1) str += ', ';
    }
    str += ']';
    if (i !== 0) str += ',';
    str += '\n';
    result += str;
  }
  result += ']';
  return { towerData: { tower, blocksInRow }, result };
};

export const getCrashResult = async (data: any) => {
  const { serverSeed, clientSeed } = data;

  const result = await verifyCrashFairness(serverSeed, clientSeed);
  return result.toFixed(2);
};
