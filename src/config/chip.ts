import { Balance, PaidBalanceType } from 'api/types/chip';

export const chipColors = {
  chip: '#ffe24b',
  coupon: '#4BE9FF'
};

export const chipCoinColors = {
  chip: {
    background: '#ffe24b',
    border: '#ffb31f'
  },
  coupon: {
    background: '#4BE9FF',
    border: '#1FAEFF'
  }
};

export const chipWrapperBackgrounds = {
  chip: 'linear-gradient(90deg, #503b00 0%, #2f2814 100%)',
  coupon: 'linear-gradient(90deg, #004150 0%, rgba(0, 65, 80, 0.25) 100%)'
};

export const initialBalances: {
  [key in PaidBalanceType]: Balance;
} = {
  chip: { balance: 0 },
  coupon: { balance: 0 }
};
