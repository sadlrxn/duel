import { createReducer } from '@reduxjs/toolkit';
import {
  reset,
  setRound,
  endRound,
  cancelRound,
  setGameData,
  setHistory,
  setHistoryWinner,
  setCreateRequest,
  setRequest,
  setFetch
} from './actions';
import { CoinflipGame } from 'api/types/coinflip';

export const initialState: CoinflipGame = {
  games: [],
  history: {
    games: [],
    winner: 'All Games'
  },
  createRequest: 0,
  fetch: false
};

export default createReducer<CoinflipGame>(initialState, builder =>
  builder
    .addCase(setGameData, (state, { payload }) => {
      state.games = payload;
    })
    .addCase(setHistory, (state, { payload }) => {
      state.history.games = payload;
    })
    .addCase(setHistoryWinner, (state, { payload }) => {
      state.history.winner = payload;
    })
    .addCase(setRound, (state, { payload }) => {
      const index = state.games.findIndex(
        round => round.roundId === payload.roundId
      );

      if (index === -1) state.games.push(payload);
      else state.games.splice(index, 1, { ...state.games[index], ...payload });
    })
    .addCase(endRound, (state, { payload }) => {
      const index = state.games.findIndex(round => round.roundId === payload);

      if (index !== -1) {
        const round = state.games.splice(index, 1);
        round[0].status = 'ended';
        round[0].time = Date.now();
        const names: string[] = [
          round[0].headsUser!.name,
          round[0].tailsUser!.name
        ];
        if (
          state.history.winner === 'All Games' ||
          names.indexOf(state.history.winner) !== -1
        )
          state.history.games.unshift(round[0]);
      }

      const length = state.history.games.length;
      if (length > Math.floor(length / 6) * 6)
        state.history.games = state.history.games.slice(
          0,
          Math.floor(length / 6) * 6
        );
    })
    .addCase(cancelRound, (state, { payload }) => {
      const index = state.games.findIndex(round => round.roundId === payload);

      if (index !== -1) state.games.splice(index, 1);
    })
    .addCase(reset, () => initialState)
    .addCase(setCreateRequest, (state, { payload }) => ({
      ...state,
      createRequest: payload
    }))
    .addCase(setRequest, (state, { payload }) => {
      const index = state.games.findIndex(
        game => game.roundId === payload.roundId
      );
      if (index !== -1) state.games[index].request = payload.status;
    })
    .addCase(setFetch, (state, { payload }) => {
      state.fetch = payload;
    })
);
