import { createAction } from '@reduxjs/toolkit';

import { CoinflipRoundData as Round } from 'api/types/coinflip';

export const setGameData = createAction<Round[]>('coinflip/game_data');

export const setHistory = createAction<Round[]>('coinflip/history');

export const setHistoryWinner = createAction<string>('coinflip/history_winner');

export const setRound = createAction<Round>('coinflip/round');

export const endRound = createAction<number>('coinfilp/end');

export const cancelRound = createAction<number>('coinflip/cancel');

export const reset = createAction('coinflip/reset');

export const setCreateRequest = createAction<number>('coinflip/create');

export const setRequest = createAction<{ roundId: number; status: boolean }>(
  'coinflip/set_request'
);

export const setFetch = createAction<boolean>('coinflip/fetchdata');
