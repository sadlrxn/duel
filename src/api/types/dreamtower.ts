import { User } from './user';

export type DreamtowerGameStatus = 'playing' | 'loss' | 'win' | 'cashout' | '';
export type DreamtowerAutoStatus = 'running' | '';
export type DreamtowerMode = 'manual' | 'auto';

export interface DreamtowerDifficulty {
  level: string;
  blocksInRow: number;
  starsInRow: number;
}

export interface DreamtowerGame {
  tower: number[][];
  roundId: number;
  difficulty: DreamtowerDifficulty;
  bets: number[];
  betAmount: number;
  multiplier: number;
  status: DreamtowerGameStatus;
  height: number;
  nextMultiplier: number | null;
  profit: number;
}

export interface DreamtowerAuto {
  status: DreamtowerAutoStatus;
  betAmount: number;
  betCount: number | undefined;
  changeBetOnWin: number | undefined;
  changeBetOnLoss: number | undefined;
  stopProfit: number | undefined;
  stopLoss: number | undefined;
  accumulated: number;
}

export interface Dreamtower {
  game: DreamtowerGame;
  auto: DreamtowerAuto;
  history: DreamtowerHistory;
}

export interface DreamtowerGameMeta {
  difficulties: DreamtowerDifficulty[];
  fee: number;
  towerHeight: number;
  minAmount: number;
  maxAmount: number;
}

export interface DreamtowerBetData {
  roundId: number | null;
  betAmount: number;
  bets: number[];
  difficulty: string;
}

export interface DreamtowerHistoryRound {
  roundId: number;
  user: User;
  difficulty: DreamtowerDifficulty;
  status: DreamtowerGameStatus;
  time: number;
  tower: number[][];
  multiplier: number;
  betAmount: number;
  bets: number[];
  profit: number;
  paidBalanceType: string;
}

export interface DreamtowerHistory {
  rounds: DreamtowerHistoryRound[];
  review: DreamtowerHistoryRound;
  winner: string;
}

export interface DreamtowerFairData extends DreamtowerHistoryRound {
  expired: boolean;
  clientSeed: string;
  serverSeed?: string;
  serverSeedHash: string;
  nonce: number;
  seedNonce: number;
}

export const defaultTower = new Array(9)
  .fill(-1)
  .map(() => new Array(4).fill(-1));

export const initialFair: DreamtowerFairData = {
  roundId: 0,
  user: { id: 0, name: '', avatar: '', count: 0 },
  difficulty: { level: 'Easy', blocksInRow: 4, starsInRow: 3 },
  status: '',
  time: Date.now(),
  tower: defaultTower,
  multiplier: 0,
  betAmount: 0,
  bets: [],
  expired: true,
  clientSeed: '',
  serverSeed: '',
  serverSeedHash: '',
  nonce: 0,
  seedNonce: 0,
  profit: 0,
  paidBalanceType: 'chip'
};
