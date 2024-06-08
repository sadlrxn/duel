import { useCallback } from 'react';
import { useConnection, useWallet } from '@solana/wallet-adapter-react';
import {
  PublicKey,
  SystemProgram,
  Transaction,
  TransactionSignature
} from '@solana/web3.js';
import {
  createAssociatedTokenAccountInstruction,
  createTransferInstruction,
  getAssociatedTokenAddress,
  TOKEN_PROGRAM_ID
} from '@solana/spl-token';
import { useAppDispatch } from 'state';
import { toast } from 'utils/toast';
import { MASTER_WALLET } from 'config';
import { updateBalance } from 'state/user/actions';
import { v4 as uuidv4 } from 'uuid';
import api from 'utils/api';
import { addLog } from 'state/log/actions';
import { useAppSelector } from 'state';
import { sendAndConfirmTransaction } from 'utils/solanaWeb3';
import { convertChipToBalance } from 'utils/balance';

const useTrading = () => {
  const { connection } = useConnection();
  const { publicKey, sendTransaction, signTransaction } = useWallet();
  const { walletAddress } = useAppSelector(state => state.user);
  const dispatch = useAppDispatch();

  const depositToken = useCallback(
    async (
      lamports: number,
      decimals: number,
      tokenPrice?: number,
      contractAddress?: string
    ) => {
      if (publicKey?.toBase58() !== walletAddress) {
        toast.warning(
          'Please make sure you are connected to the correct wallet and try again.'
        );
        return;
      }

      if (!tokenPrice) return;
      if (!publicKey) {
        toast.error(`Wallet not connected!`);
        console.info('error', `Send Transaction: Wallet not connected!`);
        return;
      }

      let signature: TransactionSignature = '';

      try {
        const transaction = new Transaction();

        if (contractAddress) {
          const mint = new PublicKey(contractAddress);
          const fromTokenAccount = await getAssociatedTokenAddress(
            mint,
            publicKey
          );

          const toTokenAccount = await getAssociatedTokenAddress(
            mint,
            MASTER_WALLET
          );

          const toTokenAccountBalance = await connection.getBalance(
            toTokenAccount
          );

          if (toTokenAccountBalance === 0) {
            transaction.add(
              createAssociatedTokenAccountInstruction(
                publicKey,
                toTokenAccount,
                MASTER_WALLET,
                mint
              )
            );
          }

          transaction.add(
            createTransferInstruction(
              fromTokenAccount, // source
              toTokenAccount, // dest
              publicKey,
              lamports,
              [],
              TOKEN_PROGRAM_ID
            )
          );
        } else {
          transaction.add(
            SystemProgram.transfer({
              fromPubkey: publicKey,
              toPubkey: MASTER_WALLET,
              lamports
            })
          );
        }

        signature = await sendAndConfirmTransaction(
          connection,
          transaction,
          publicKey,
          signTransaction,
          sendTransaction
        );

        dispatch(
          addLog({
            type: 'Deposit',
            signature,
            data: (lamports * tokenPrice) / 10 ** decimals,
            status: 'Pending',
            time: Date.now()
          })
        );

        toast.success(`Deposit request successful`);
      } catch (error: any) {
        toast.error(`Transaction failed!`);
        console.info(
          'error',
          `Transaction failed! ${error?.message}`,
          signature
        );
        return;
      }
    },
    [
      connection,
      dispatch,
      publicKey,
      signTransaction,
      sendTransaction,
      walletAddress
    ]
  );

  const depositNft = useCallback(
    async (nfts: any[]) => {
      if (publicKey?.toBase58() !== walletAddress) {
        toast.warning(
          'Please make sure you are connected to the correct wallet and try again.'
        );
        return;
      }

      const mintAddresses = nfts.map(value => value.mintAddress);
      const imgUrls = nfts.map(value => value.image);
      if (!publicKey) {
        toast.error(`Wallet not connected!`);
        console.info('error', `Send Transaction: Wallet not connected!`);
        return;
      }
      const transaction = new Transaction();

      for (let i = 0; i < mintAddresses.length; i++) {
        const mint = new PublicKey(mintAddresses[i]);

        const fromTokenAccount = await getAssociatedTokenAddress(
          mint,
          publicKey
        );

        const toTokenAccount = await getAssociatedTokenAddress(
          mint,
          MASTER_WALLET
        );

        const toTokenAccountBalance = await connection.getBalance(
          toTokenAccount
        );

        if (toTokenAccountBalance === 0) {
          transaction.add(
            createAssociatedTokenAccountInstruction(
              publicKey,
              toTokenAccount,
              MASTER_WALLET,
              mint
            )
          );
        }

        transaction.add(
          createTransferInstruction(
            fromTokenAccount, // source
            toTokenAccount, // dest
            publicKey,
            1,
            [],
            TOKEN_PROGRAM_ID
          )
        );
      }
      let signature: TransactionSignature = '';
      try {
        signature = await sendAndConfirmTransaction(
          connection,
          transaction,
          publicKey,
          signTransaction,
          sendTransaction
        );

        dispatch(
          addLog({
            type: 'Deposit',
            signature,
            status: 'Pending',
            data: imgUrls,
            time: Date.now()
          })
        );
        toast.success(`Deposit request successful`);
      } catch (error: any) {
        toast.error(`Transaction failed!`);
        console.info(
          'error',
          `Transaction failed! ${error?.message}`,
          signature
        );
        return;
      }
    },
    [
      publicKey,
      walletAddress,
      connection,
      signTransaction,
      sendTransaction,
      dispatch
    ]
  );

  const withdrawSol = useCallback(
    async (usdAmount: number, keyword: string) => {
      try {
        // const txId = uuidv4();
        dispatch(
          updateBalance({
            type: -1,
            usdAmount: convertChipToBalance(usdAmount)
          })
        );
        if (keyword !== 'SOL')
          dispatch(updateBalance({ type: -1, usdAmount: 10 }));
        const result = await api.post<{
          status: number;
          txId: string;
          amount: number;
        }>('/pay/withdraw/sol', {
          usdAmount: Math.ceil(convertChipToBalance(usdAmount)),
          targetToken: keyword
        });

        dispatch(
          addLog({
            type: 'Withdraw',
            signature: result.data.txId,
            status: 'Pending',
            data: usdAmount,
            time: Date.now()
          })
        );

        toast.success(`Withdraw request successful`);
      } catch (error: any) {
        dispatch(
          updateBalance({ type: 1, usdAmount: convertChipToBalance(usdAmount) })
        );
        if (keyword !== 'SOL')
          dispatch(updateBalance({ type: 1, usdAmount: 10000 }));
        if (error.response.status === 429) {
          toast.error(error.response.data.message);
        } else if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error(error.response.data.status);
      }
    },
    [dispatch]
  );

  const withdrawNft = useCallback(
    async (nfts: any[]) => {
      const mintAddresses = nfts.map(value => value.mintAddress);
      const imgUrls = nfts.map(value => value.image);
      try {
        const txId = uuidv4();
        const result = await api.post<{
          status: number;
          txId: string;
          mintAddresses: string[];
        }>('/pay/withdraw/nft', {
          mintAddresses,
          txId
        });
        dispatch(updateBalance({ type: -1, nfts: nfts }));

        dispatch(
          addLog({
            type: 'Withdraw',
            signature: result.data.txId,
            status: 'Pending',
            data: imgUrls,
            time: Date.now()
          })
        );
        toast.success(`Withdraw request successful`);
      } catch (error: any) {
        dispatch(updateBalance({ type: 1, nfts: nfts }));
        if (error.response.status === 429) {
          toast.error(error.response.data.message);
        } else if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error(error.response.data.status);
      }
    },
    [dispatch]
  );

  return { depositNft, withdrawSol, withdrawNft, depositToken };
};

export default useTrading;
