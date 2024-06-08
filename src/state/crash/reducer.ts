import { createReducer } from '@reduxjs/toolkit';

import {
  setGameData,
  endRound,
  reset,
  setFetch,
  clearRoundData,
  setPending,
  setGameDataByProperties
} from './actions';
import { CrashGame } from 'api/types/crash';

const initialState: CrashGame = {
  roundId: 0,
  history: [],
  status: 'bet',
  from: 1,
  to: 1,
  serverTimeElapsed: 0,
  startedTime: Date.now(),
  time: Date.now(),
  bets: [],
  cashIns: [],
  cashOuts: [],
  fetch: false
};

export default createReducer(initialState, builder =>
  builder
    .addCase(setGameData, (state, { payload }) => {
      Object.keys(payload).forEach(key => {
        //@ts-ignore
        state[key] = payload[key];
      });
    })
    .addCase(endRound, (state, { payload }) => {
      const index = state.history.findIndex(
        item => item.roundId === payload.roundId
      );
      if (index !== -1) return;
      state.history.unshift(payload);
      state.history = state.history.slice(0, 10);
    })
    .addCase(setFetch, (state, { payload }) => {
      state.fetch = payload;
    })
    .addCase(clearRoundData, state => ({
      ...initialState,
      fetch: state.fetch
    }))
    .addCase(setPending, (state, { payload }) => {
      const index = state.bets.findIndex(bet => bet.betId === payload.betId);
      if (index === -1) return;
      state.bets[index].pending = payload.status;
    })
    .addCase(setGameDataByProperties, (state, { payload }) => {
      payload.properties.forEach((property, index) => {
        //@ts-ignore
        state[property] = payload.data[index];
        console.info(property, payload.data[index]);
      });
    })
    .addCase(reset, () => initialState)
    .addDefaultCase(state => state)
);
