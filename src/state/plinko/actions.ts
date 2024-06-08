import { createAction } from '@reduxjs/toolkit';

import { PlinkoBall, PlinkoHistory, PlinkoLevel } from 'api/types/plinko';

export const addBall = createAction<PlinkoBall>('plinko/add_ball');

export const removeBall = createAction<number>('plinko/remove_ball');

export const setAutoBet = createAction<number>('plinko/set_auto_bet');

export const setLines = createAction<number>('plinko/set_lines');

export const setLevel = createAction<PlinkoLevel>('plinko/set_level');

export const reset = createAction('plinko/reset');

export const setHistory = createAction<PlinkoHistory[]>('plinko/set_history');
