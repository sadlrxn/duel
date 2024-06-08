import { UserRole } from 'config';
import { NFT } from './nft';
// import { User } from './user';

export const jackpotRooms = {
  wild: 'wild',
  medium: 'medium',
  low: 'low'
};

export interface User {
  id: number;
  name: string;
  avatar: string;
  role?: UserRole;
}

export interface Candidate extends User {
  count?: number;
  percent?: number;
}

export interface BetPlayer extends Candidate {
  usdAmount: number;
  nftAmount?: number;
  nfts?: NFT[];
}

export type TGameStatus =
  | 'available'
  | 'created'
  | 'started'
  | 'finished'
  | 'counting'
  | 'rolling'
  | 'rollend';

export type JackpotRoom = 'low' | 'medium' | 'wild';

export interface JackpotAutoBet {
  chip?: number;
  nfts?: NFT[];
}

export interface JackpotGame {
  low: JackpotRoomGame;
  medium: JackpotRoomGame;
  wild: JackpotRoomGame;
  room: JackpotRoom;
  history: JackpotHistory;
}

export interface JackpotRoomGame {
  game: JackpotRoundData;
  // history: JackpotHistory;
  fetch: boolean;
  countText?: string;
  autoBet?: JackpotAutoBet;
}

export interface GrandJackpotGame {
  game: JackpotRoundData;
  history: JackpotHistory;
  fetch: boolean;
}

export interface JackpotGameMeta {
  minBetAmount: number;
  maxBetAmount: number;
  betCountLimit: number;
  playerLimit: number;
  countingTime: number;
  rollingTime: number;
  winnerTime: number;
  fee: number;
}

export interface GrandJackpotGameMeta {
  minBetAmount: number;
  bettingTime: number;
  countingTime: number;
  rollingTime: number;
  winnerTime: number;
  fee: number;
}

export interface JackpotHistory {
  winner: string;
  games: JackpotHistoryData[];
}

export interface JackpotRoundData {
  status: TGameStatus;
  roundId: number;
  ticketId?: string;
  signedString?: string;
  players: BetPlayer[];
  nfts: NFT[];
  winner: User;
  usdBetAmount: number;
  nftBetAmount: number;
  totalBetAmount: number;
  time: number;
  rolltime: number;
  animationTime: number;
  candidates: Candidate[];
  request: boolean;
  usdProfit?: number;
  nftProfit?: NFT[];
  usdFee?: number;
  nftFee?: NFT[];
  countingTime?: number;
}

export interface JackpotHistoryData {
  roundId: number;
  ticketId: string;
  signedString: string;
  players: User[];
  winner: User;
  chance: number;
  prize: number;
  time: number;
}

export interface JackpotBet {
  amount: number;
  nfts: NFT[];
}

export interface JackpotFairData {
  roundId: number;
  ticketId: string;
  signedString: string;
  players: BetPlayer[];
  winner: User;
}

export const initialFair: JackpotFairData = {
  roundId: 0,
  ticketId: '',
  signedString: '',
  players: [],
  winner: { id: 0, name: '', avatar: '' }
};
