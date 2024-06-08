import { useConnection } from '@solana/wallet-adapter-react';
import { PublicKey } from '@solana/web3.js';
import { getAssociatedTokenAddress } from '@solana/spl-token';
import { useCallback, useEffect, useState } from 'react';
import { shallowEqual } from 'react-redux';
import { useAppSelector } from 'state';
import { formatLamports2Sol } from 'utils/format';

export default function useUserTokenAmounts() {
  const walletAddress = useAppSelector(
    state => state.user.walletAddress,
    shallowEqual
  );
  const tokens = useAppSelector(state => state.token.tokens, shallowEqual);
  const { connection } = useConnection();

  const [balances, setBalances] = useState<number[]>([]);

  const fetchTokens = useCallback(async () => {
    let balances: number[] = [];
    for (let i = 0; i < tokens.length; i++) {
      try {
        if (tokens[i].keyword.toLowerCase().slice(0, 3) === 'sol') {
          let lamports = await connection.getBalance(
            new PublicKey(walletAddress)
          );
          balances.push(formatLamports2Sol(lamports));
        } else {
          let tokenAccount = await getAssociatedTokenAddress(
            new PublicKey(tokens[i].mintAddress!),
            new PublicKey(walletAddress)
          );
          let splBal = await connection.getTokenAccountBalance(tokenAccount);
          balances.push(splBal.value.uiAmount!);
        }
      } catch (err: any) {
        if (balances.length < i + 1) {
          balances.push(0);
        }
      }
    }
    setBalances(balances);
  }, []);

  // const fetchTokens = async () => {
  //   let lamports = await connection.getBalance(new PublicKey(walletAddress));

  //   setSolBalance(formatLamports2Sol(lamports));

  //   let bonkTokenAccount = await getAssociatedTokenAddress(
  //     new PublicKey('DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263'),
  //     new PublicKey(walletAddress)
  //   );
  //   let bonkBal = await connection.getTokenAccountBalance(bonkTokenAccount);

  //   setBonkBalance(bonkBal.value.uiAmount!);

  //   let usdcTokenAccount = await getAssociatedTokenAddress(
  //     new PublicKey('EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v'),
  //     new PublicKey(walletAddress)
  //   );
  //   let usdcBal = await connection.getTokenAccountBalance(usdcTokenAccount);

  //   setUsdcBalance(usdcBal.value.uiAmount!);
  // };

  useEffect(() => {
    fetchTokens();
  }, []);

  // return useMemo(
  //   () => [solBalance, bonkBalance, usdcBalance],
  //   [solBalance, bonkBalance, usdcBalance]
  // );
  return balances;
}
