import { User as Candidate } from './user';

export type CoinflipGameStatus = 'created' | 'joined' | 'ended';

export interface CoinflipGame {
  games: CoinflipRoundData[];
  history: {
    games: CoinflipRoundData[];
    winner: string;
  };
  createRequest: number;
  fetch: boolean;
}

export interface CoinflipGameMeta {
  flipTime: number;
  fee: number;
  createRoundLimit: number;
  minBetAmount: number;
  maxBetAmount: number;
}

export interface CoinflipBasicData {
  roundId: number;
  ticketId: string;
  signedString: string;
  headsUser: Candidate | null;
  tailsUser: Candidate | null;
  amount: number;
  winnerId?: number;
  creatorId: number;
  paidBalanceType: 'chip' | 'coupon';
}

export interface CoinflipRoundData extends CoinflipBasicData {
  status: CoinflipGameStatus;
  prize: number;
  time: number;
  request: boolean;
}

export interface CoinflipBetData {
  roundId: number | null;
  amount: number;
  side: 'heads' | 'tails';
  opponent: 'dueler' | 'bot';
}

export interface CoinflipSideData {
  avatar: string;
  name: string;
  amount: number;
  creator: boolean;
}

export const initialRound: CoinflipRoundData = {
  status: 'created',
  roundId: 0,
  headsUser: null,
  tailsUser: null,
  amount: 0,
  prize: 0,
  ticketId: '',
  signedString: '',
  winnerId: 0,
  creatorId: 0,
  paidBalanceType: 'chip',
  time: 0,
  request: false
};

export interface CoinflipFairData extends CoinflipBasicData {}

export const initialFair: CoinflipFairData = {
  roundId: 0,
  ticketId: '',
  signedString: '',
  amount: 0,
  winnerId: 0,
  creatorId: 0,
  paidBalanceType: 'chip',
  headsUser: {
    id: 0,
    name: '',
    avatar: '',
    count: 0
  },
  tailsUser: {
    id: 0,
    name: '',
    avatar: '',
    count: 0
  }
};
