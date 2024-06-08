import { createReducer } from '@reduxjs/toolkit';

import { PlinkoGame } from 'api/types/plinko';
import {
  addBall,
  removeBall,
  setAutoBet,
  setLines,
  setLevel,
  reset,
  setHistory
} from './actions';

export const initialState: PlinkoGame = {
  balls: [],
  lines: 8,
  level: 'low',
  autoBet: 0,
  history: []
};

export default createReducer<PlinkoGame>(initialState, builder =>
  builder
    .addCase(addBall, (state, { payload }) => {
      if (
        state.balls.findIndex(ball => ball.roundId === payload.roundId) === -1
      ) {
        state.balls.push(payload);
      }
    })
    .addCase(removeBall, (state, { payload }) => {
      const index = state.balls.findIndex(ball => ball.roundId === payload);
      if (index !== -1) {
        state.balls.splice(index, 1);
      }
    })
    .addCase(setAutoBet, (state, { payload }) => {
      state.autoBet = payload;
    })
    .addCase(setLines, (state, { payload }) => {
      state.lines = payload;
    })
    .addCase(setLevel, (state, { payload }) => {
      state.level = payload;
    })
    .addCase(setHistory, (state, { payload }) => {
      state.history = payload.slice().sort((h1, h2) => h2.roundId - h1.roundId);
    })
    .addCase(reset, () => initialState)
    .addDefaultCase(state => state)
);
