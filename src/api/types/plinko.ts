export type LinesType = 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16;

export type PlinkoLevel = 'low' | 'medium' | 'high';

export interface PlinkoBall {
  roundId: number;
  path: string;
  betAmount: number;
  lines: number;
  level: PlinkoLevel;
  time: number;
  multiplier: number;
}

export interface PlinkoGame {
  balls: PlinkoBall[];
  lines: number;
  level: PlinkoLevel;
  autoBet: number;
  history: PlinkoHistory[];
  review?: PlinkoBall;
}

export interface PlinkoMeta {
  difficulties: string[];
  minBetAmount: number;
  maxBetAmount: number;
  betCountLimit: number;
}

export interface PlinkoHistory {
  roundId: number;
  multiplier: number;
}
