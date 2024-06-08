import { User } from './user';
import { PaidBalanceType } from './chip';

export type CrashGameStatus = 'bet' | 'ready' | 'play' | 'explosion' | 'back';
export type CrashGameStatusServer =
  | 'crash-status-betting'
  | 'crash-status-pending'
  | 'crash-status-running'
  | 'crash-status-preparing';

export interface CrashGameHistory {
  roundId: number;
  multiplier: number;
}

export interface CrashGame {
  roundId: number;
  history: CrashGameHistory[];
  status: CrashGameStatus;
  from: number;
  to: number;
  serverTimeElapsed: number;
  startedTime: number;
  time: number;
  bets: CrashBet[];
  cashIns: CrashCashEvent[];
  cashOuts: CrashCashEvent[];
  fetch: boolean;
}

export interface CrashMeta {
  bettingDuration: number; //Betting time
  pendingDuration: number; //Prepare rocket
  preparingDuration: number; // explosion & update history
  betCountLimit: number;
  minBetAmount: number;
  maxBetAmount: number;
  houseEdge: number;
  maxPlayerLimit: number;
  minCashOutAt: number;
  maxCashOut: number;
  eventInterval: number;
}

export interface CrashServerEvent {
  roundId: number;
  roundStatus: CrashGameStatusServer;
  elapsed: number;
  multiplier: number;
  runStartedAt?: number;
  bets?: CrashCashEvent[];
  cashIns?: CrashCashEvent[];
  cashOuts?: CrashCashEvent[];
}

export interface CrashCashEvent {
  roundId: number;
  user: User;
  amount: number;
  balanceType: PaidBalanceType;
  betId: number;
  cashOutAt?: number;
  multiplier?: number;
}

export interface CrashBet {
  roundId: number;
  betId: number;
  user: User;
  betAmount: number;
  paidBalanceType: PaidBalanceType;
  profit?: number;
  payoutMultiplier?: number;
  cashOutAt?: number;
  pending: boolean;
}

export interface CrashAutoBet {
  betAmount: number;
  paidBalanceType: PaidBalanceType;
  pnl: number;
  cashOutAt: number;
  betId: number;
  roundId: number;
  profit: number;
  bettedRounds: number;
  isComplete: boolean;
  isBetted: boolean;
  rounds: number;
  onLoss: number;
  onWin: number;
  stopProfit: number;
  stopLoss: number;
}

export interface CrashFairData {
  roundId: number;
  seed: string;
  outcome: number;
  date: number;
  bets: CrashBet[];
}

export const initialFair: CrashFairData = {
  roundId: 0,
  seed: '',
  outcome: 1,
  date: Date.now(),
  bets: []
};

export const initialAutoBet: CrashAutoBet = {
  betAmount: 100000,
  paidBalanceType: 'chip',
  cashOutAt: 2.5,
  rounds: -1,
  pnl: 0,
  betId: -1,
  roundId: -1,
  profit: 0,
  bettedRounds: 0,
  isComplete: false,
  isBetted: false,
  onLoss: 0,
  onWin: 0,
  stopProfit: 0,
  stopLoss: 0
};
