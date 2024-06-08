import { createReducer } from '@reduxjs/toolkit';
import {
  Dreamtower,
  DreamtowerAuto,
  DreamtowerGame,
  DreamtowerHistory,
  DreamtowerHistoryRound
} from 'api/types/dreamtower';
import {
  updateStatus,
  raise,
  setBetAmount,
  setDifficulty,
  reset,
  setAutoStatus,
  setRound,
  setAutoBetAmount,
  setAutoBetCount,
  setChangeBetOnWin,
  setChangeBetOnLoss,
  setStopProfit,
  setStopLoss,
  setAccumulated,
  setAutoPath,
  setHistory,
  review,
  setHistoryWinner,
  clear
} from './actions';

const defaultTower = new Array(9).fill(-1).map(() => new Array(4).fill(-1));

export const initialGameData: DreamtowerGame = {
  tower: defaultTower,
  roundId: 0,
  bets: [],
  betAmount: 0,
  difficulty: { level: 'Easy', blocksInRow: 4, starsInRow: 3 },
  multiplier: 0,
  status: '',
  height: 9,
  nextMultiplier: 1.29,
  profit: 0
};

export const initialAutoData: DreamtowerAuto = {
  status: '',
  betAmount: 0,
  betCount: undefined,
  changeBetOnWin: undefined,
  changeBetOnLoss: undefined,
  stopProfit: undefined,
  stopLoss: undefined,
  accumulated: 0
};

export const initialReview: DreamtowerHistoryRound = {
  betAmount: 0,
  bets: [],
  difficulty: { level: 'Easy', blocksInRow: 4, starsInRow: 3 },
  multiplier: 0,
  roundId: 0,
  status: '',
  time: 0,
  tower: defaultTower,
  user: { id: 0, name: '', avatar: '', count: 0 },
  profit: 0,
  paidBalanceType: 'chip'
};

export const initialHistory: DreamtowerHistory = {
  winner: 'All Games',
  rounds: [],
  review: initialReview
};

export const initialState: Dreamtower = {
  game: initialGameData,
  auto: initialAutoData,
  history: initialHistory
};

export default createReducer<Dreamtower>(initialState, builder =>
  builder
    .addCase(setHistory, (state, { payload }) => {
      state.history.rounds = payload;
      const reviewIndex = state.history.rounds.findIndex(
        round => round.roundId === state.history.review.roundId
      );
      if (reviewIndex === -1) state.history.review = initialReview;
    })
    .addCase(review, (state, { payload }) => {
      state.history.review = payload;
    })
    .addCase(setHistoryWinner, (state, { payload }) => {
      state.history.winner = payload;
    })
    .addCase(setRound, (state, { payload }) => {
      var temp: number[][] = new Array(state.game.height)
        .fill(-1)
        .map(() => new Array(payload.difficulty.blocksInRow).fill(-1));
      state.game.tower = temp;
      state.game.roundId = payload.roundId;
      state.game.difficulty = payload.difficulty;
      state.game.betAmount = payload.betAmount;
      state.game.bets = payload.bets;
      state.game.status = payload.status;
      state.game.multiplier = payload.multiplier;
      state.game.nextMultiplier = payload.nextMultiplier;
      state.game.profit =
        payload.maxWinning < payload.betAmount * payload.multiplier
          ? payload.maxWinning
          : payload.betAmount * payload.multiplier;
      payload.bets.forEach((bet: number, index: number) => {
        state.game.tower[index][bet] = 1;
      });
    })
    .addCase(setDifficulty, (state, { payload }) => {
      state.game.difficulty = payload;
      var temp: number[][] = new Array(state.game.height)
        .fill(-1)
        .map(() => new Array(state.game.difficulty.blocksInRow).fill(-1));
      state.game.tower = temp;
      state.game.bets = [];
      state.game.status = '';
      state.game.nextMultiplier =
        ((payload.blocksInRow / payload.starsInRow) * 97) / 100;
    })
    .addCase(setBetAmount, (state, { payload }) => {
      state.game.betAmount = payload;
    })
    .addCase(reset, state => {
      var temp: number[][] = new Array(state.game.height)
        .fill(-1)
        .map(() => new Array(state.game.difficulty.blocksInRow).fill(-1));
      state.game.tower = temp;
      state.game.bets = [];
      state.game.status = '';
      state.auto.accumulated = 0;
      state.auto.betCount = undefined;
    })
    .addCase(clear, () => initialState)
    .addCase(updateStatus, (state, { payload }) => {
      state.game.roundId = payload.roundId;
      state.game.multiplier = payload.multiplier;
      state.game.status = payload.status;
      state.game.profit =
        payload.maxWinning < state.game.betAmount * payload.multiplier
          ? payload.maxWinning
          : state.game.betAmount * payload.multiplier;

      if (payload.status !== 'playing') {
        state.game.tower = payload.tower;
      } else if (payload.status === 'playing') {
        state.game.nextMultiplier = payload.nextMultiplier;
      }
    })
    .addCase(raise, (state, { payload }) => {
      if (
        payload.status === 'playing' ||
        payload.status === 'win' ||
        payload.status === 'cashout'
      )
        state.game.tower[state.game.bets.length][payload.bet] = 1;
      state.game.bets.push(payload.bet);
    })

    .addCase(setAutoPath, (state, { payload }) => {
      if (state.game.bets.length <= payload.height) {
        state.game.bets.push(payload.index);
      } else if (state.game.bets[payload.height] === payload.index) {
        state.game.bets = state.game.bets.slice(0, payload.height);
      } else {
        state.game.bets[payload.height] = payload.index;
      }
      var nextMultiplier =
        state.game.difficulty.blocksInRow / state.game.difficulty.starsInRow;
      nextMultiplier = Math.pow(nextMultiplier, state.game.bets.length + 1);
      nextMultiplier =
        nextMultiplier -
        (nextMultiplier * 3 * (state.game.bets.length + 1)) / 100;
      nextMultiplier = Math.floor(nextMultiplier * 100) / 100;
      state.game.nextMultiplier = nextMultiplier;
    })
    .addCase(setAutoStatus, (state, { payload }) => {
      state.auto.status = payload;
      if (payload === '') {
        var temp: number[][] = new Array(state.game.height)
          .fill(-1)
          .map(() => new Array(state.game.difficulty.blocksInRow).fill(-1));
        state.game.tower = temp;
        state.game.status = '';
        state.auto.accumulated = 0;
        state.game.betAmount = state.auto.betAmount;
      }
    })
    .addCase(setAutoBetAmount, (state, { payload }) => {
      state.auto.betAmount = payload;
    })
    .addCase(setAutoBetCount, (state, { payload }) => {
      state.auto.betCount = payload;
    })
    .addCase(setChangeBetOnWin, (state, { payload }) => {
      state.auto.changeBetOnWin = payload;
    })
    .addCase(setChangeBetOnLoss, (state, { payload }) => {
      state.auto.changeBetOnLoss = payload;
    })
    .addCase(setStopProfit, (state, { payload }) => {
      state.auto.stopProfit = payload;
    })
    .addCase(setStopLoss, (state, { payload }) => {
      state.auto.stopLoss = payload;
    })
    .addCase(setAccumulated, (state, { payload }) => {
      state.auto.accumulated = payload;
    })
    .addDefaultCase(state => state)
);
