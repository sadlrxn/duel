import { createAction } from '@reduxjs/toolkit';

import { CrashGameHistory } from 'api/types/crash';

export const setGameData = createAction<any>('crash/set_game_data');

export const endRound = createAction<CrashGameHistory>('crash/end_round');

export const reset = createAction('crash/reset');

export const clearRoundData = createAction('crash/clear_round_data');

export const setFetch = createAction<boolean>('crash/set_fetch');

export const setPending = createAction<{ betId: number; status: boolean }>(
  'crash/set_pending'
);

export const setGameDataByProperties = createAction<{
  properties: string[];
  data: any[];
}>('crash/set_game_data_by_properties');
