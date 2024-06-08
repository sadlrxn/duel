import { sortBets } from './sort';
import { CrashBet } from 'api/types/crash';

describe('CrashBet sort', () => {
  const bets: CrashBet[] = [
    {
      roundId: 2943,
      betId: 2467,
      user: {
        avatar:
          'https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png',
        id: 30,
        name: 'Bomber',
        role: 'admin',
        count: 1
      },
      betAmount: 8000,
      paidBalanceType: 'chip',
      cashOutAt: 1.8,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2468,
      user: {
        avatar:
          'https://duelana-bucket-test.s3.us-east-2.amazonaws.com/avatar/8AgbhSZqHWjiNBsKkyuoRsTG3a3mubKZze9U7WFvCg28.png',
        id: 6,

        name: 'zog1',
        role: 'ambassador',
        count: 1
      },
      betAmount: 9000,
      paidBalanceType: 'chip',
      cashOutAt: 1.6,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2469,
      user: {
        avatar:
          'https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png',
        id: 30,

        name: 'Bomber',
        role: 'admin',
        count: 1
      },
      betAmount: 8000,
      paidBalanceType: 'chip',
      cashOutAt: 1.8,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2470,
      user: {
        avatar:
          'https://duelana-bucket-test.s3.us-east-2.amazonaws.com/avatar/8AgbhSZqHWjiNBsKkyuoRsTG3a3mubKZze9U7WFvCg28.png',
        id: 6,

        name: 'zog1',
        role: 'ambassador',
        count: 1
      },
      betAmount: 9000,
      paidBalanceType: 'chip',
      cashOutAt: 1.6,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2471,
      user: {
        avatar:
          'https://duelana-bucket-test.s3.us-east-2.amazonaws.com/avatar/8AgbhSZqHWjiNBsKkyuoRsTG3a3mubKZze9U7WFvCg28.png',
        id: 6,

        name: 'zog1',
        role: 'ambassador',
        count: 1
      },
      betAmount: 9000,
      paidBalanceType: 'chip',
      cashOutAt: 1.6,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2472,
      user: {
        avatar:
          'https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png',
        id: 30,

        name: 'Bomber',
        role: 'admin',
        count: 1
      },
      betAmount: 8000,
      paidBalanceType: 'chip',
      cashOutAt: 1.8,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2473,
      user: {
        avatar:
          'https://duelana-bucket-test.s3.us-east-2.amazonaws.com/avatar/8AgbhSZqHWjiNBsKkyuoRsTG3a3mubKZze9U7WFvCg28.png',
        id: 6,

        name: 'zog1',
        role: 'ambassador',
        count: 1
      },
      betAmount: 10000,
      paidBalanceType: 'chip',
      cashOutAt: 1.5,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2474,
      user: {
        avatar:
          'https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png',
        id: 30,

        name: 'Bomber',
        role: 'admin',
        count: 1
      },
      betAmount: 8000,
      paidBalanceType: 'chip',
      cashOutAt: 1.8,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2475,
      user: {
        avatar:
          'https://duelana-bucket-test.s3.us-east-2.amazonaws.com/avatar/8AgbhSZqHWjiNBsKkyuoRsTG3a3mubKZze9U7WFvCg28.png',
        id: 6,

        name: 'zog1',
        role: 'ambassador',
        count: 1
      },
      betAmount: 10000,
      paidBalanceType: 'chip',
      cashOutAt: 1.5,
      pending: false
    },
    {
      roundId: 2943,
      betId: 2476,
      user: {
        avatar:
          'https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png',
        id: 30,

        name: 'Bomber',
        role: 'admin',
        count: 1
      },
      betAmount: 8000,
      paidBalanceType: 'chip',
      cashOutAt: 1.8,
      pending: false
    }
  ];

  it('should sort crash bets', () => {
    const userId = 30;
    const sorted = bets.slice().sort((a, b) => sortBets(a, b, userId));

    console.dir(
      sorted.map(bet => ({ userId: bet.user.id })),
      { depth: null }
    );
  });
});
