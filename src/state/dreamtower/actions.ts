import { createAction } from '@reduxjs/toolkit';
import {
  DreamtowerDifficulty,
  DreamtowerHistoryRound
} from 'api/types/dreamtower';

export const setHistory = createAction<DreamtowerHistoryRound[]>(
  'dreamtower/setHistory'
);
export const review = createAction<DreamtowerHistoryRound>('dreamtower/review');
export const setHistoryWinner = createAction<string>(
  'dreamtower/requestHistory'
);

export const setDifficulty = createAction<DreamtowerDifficulty>(
  'dreamtower/difficulty'
);
export const setBetAmount = createAction<number>('dreamtower/betAmount');
export const reset = createAction<void>('dreamtower/reset');
export const updateStatus = createAction<any>('dreamtower/update');
export const raise = createAction<any>('dreamtower/raise');
export const setAutoStatus = createAction<any>('dreamtower/autoStatus');
export const setRound = createAction<any>('dreamtower/setRound');
export const clear = createAction<void>('dreamtower/clear');

export const setAutoPath = createAction<{ height: number; index: number }>(
  'dreamtower/setPath'
);
export const setAutoBetAmount = createAction<number>(
  'dreamtower/autoBetAmount'
);
export const setAutoBetCount = createAction<number | undefined>(
  'dreamtower/autoBetCount'
);
export const setChangeBetOnWin = createAction<number | undefined>(
  'dreamtower/changeBetOnWin'
);
export const setChangeBetOnLoss = createAction<number | undefined>(
  'dreamtower/changeBetOnLoss'
);
export const setStopProfit = createAction<number | undefined>(
  'dreamtower/stopProfit'
);
export const setStopLoss = createAction<number | undefined>(
  'dreamtower/stopLoss'
);
export const setAccumulated = createAction<number>('dreamtower/accumulated');

// export const endGame = createAction<void>("dreamtower/end");
