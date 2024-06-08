// export type PaidBalanceType = 'chip' | 'coupon';

export const BalanceTypes = ['chip', 'coupon'] as const;

export type PaidBalanceType = typeof BalanceTypes[number];

export interface Balance {
  balance: number;
  code?: string;
  wagered?: number;
  wagerLimit?: number;
  expiredTime?: number;
}
