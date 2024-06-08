import { PublicKey } from '@solana/web3.js';

export const MASTER_WALLET = new PublicKey(
  process.env.REACT_APP_MASTER_WALLET_PUBLIC_KEY || ''
);

// const IMAGEPROXY_URL = window.origin + '/imageproxy/';
const IMAGEPROXY_URL =
  process.env.REACT_APP_STAGE === 'beta'
    ? 'https://duel.win/imageproxy/'
    : 'https://staging.duel.win/imageproxy/';

export const imageProxy = (size: number = 100) => {
  return IMAGEPROXY_URL + size + '/';
};

export * from './emojis';
export * from './commands';
export * from './countrycodes';
