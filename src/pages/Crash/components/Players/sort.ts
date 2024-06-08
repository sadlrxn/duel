import { CrashBet } from 'api/types/crash';

export const sortBets = (a: CrashBet, b: CrashBet, userId: number): number => {
  if (a.user.id === userId) return 0;
  if (b.user.id === userId) return 1;
  return 0;
  // if (a.user.id === userId) {
  //   if (b.user.id === userId) {
  //     return (
  //       b.betAmount * (b.payoutMultiplier ? b.payoutMultiplier : 1) -
  //       a.betAmount * (a.payoutMultiplier ? a.payoutMultiplier : 1)
  //     );
  //   }
  //   return 0;
  // }
  // if (b.user.id === userId) return 1;
  // return b.betAmount - a.betAmount;
};
