import { useCallback, useEffect } from 'react';
import { useConnection, useWallet } from '@solana/wallet-adapter-react';
import { useWalletModal } from '@solana/wallet-adapter-react-ui';
import { useAppDispatch, useAppSelector } from 'state';
import {
  logoutUser,
  requestLogin,
  setSelfExcludeTime
} from 'state/user/actions';
import { wsDisconnect } from 'state/socket';
import api from 'utils/api';
import useFetchUserInfo from './useFetchUserInfo';
import { fetchUserInfo } from 'utils/fetchData';
import { buildAuthTx } from 'utils/ledgerSign';
import RestrictModal from 'components/Modals/RestrictModal';
import { useModal } from 'components/Modal';

const useAuth = () => {
  const dispatch = useAppDispatch();
  const {
    connected,
    publicKey,
    disconnect,
    signMessage,
    signTransaction,
    wallet
  } = useWallet();
  const { connection } = useConnection();

  const { setVisible } = useWalletModal();

  const { fetchUserNfts, fetchUserStatistics } = useFetchUserInfo();

  const { usingLedger, request, isAuthenticated } = useAppSelector(
    state => state.user
  );

  const [onPresentRestrict] = useModal(<RestrictModal />, true);

  const checkAuth = useCallback(async () => {
    try {
      // dispatch(logoutUser());
      if (!publicKey) {
        return;
        // throw new Error('Wallet not connected!');
      }
      await fetchUserInfo();

      dispatch(wsDisconnect());
    } catch (error: any) {}
  }, [dispatch, publicKey]);

  useEffect(() => {
    checkAuth();
  }, [checkAuth, connected, publicKey]);

  useEffect(() => {
    const handleLogin = async () => {
      if (!publicKey || !wallet) throw new Error('Wallet not connected!');
      if (!signMessage)
        throw new Error('Wallet does not support message signing!');

      let nonce;
      try {
        const res = await api.post('/user/requestNonce', {
          walletAddress: publicKey.toBase58()
        });

        nonce = res.data;
      } catch (error: any) {
        // console.log(error);

        if (error.response.status === 403) {
          let timeVal = error.response.data.timeRemaining;
          dispatch(setSelfExcludeTime(timeVal));
          onPresentRestrict();
          dispatch(requestLogin(false));
          return;
        }
      }

      // ledger support

      if (usingLedger) {
        // (wallet.adapter.name === 'Ledger') {
        const authTx = buildAuthTx(`Sign in Duelana with nonce: ${nonce}`);
        authTx.feePayer = publicKey; // not sure if needed but set this properly
        authTx.recentBlockhash = (
          await connection.getLatestBlockhash()
        ).blockhash; //'574Xzgj9QSvXiJB6CcMbm6pyhZQuakE7737EpfwfYEAm'; // same as line above
        const signedAuthTx = await signTransaction!(authTx);

        await api.post('/user/login', {
          walletAddress: publicKey.toBase58(),
          sigOrTx: Array.from(signedAuthTx.serialize())
        });
      }
      // ledger end
      else {
        const message = new TextEncoder().encode(
          `Sign in Duelana with nonce: ${nonce}`
        );
        const signature = await signMessage(message);

        try {
          await api.post('/user/login', {
            walletAddress: publicKey.toBase58(),
            sigOrTx: Array.from(signature)
          });
        } catch (error) {}
      }

      try {
        await fetchUserInfo();
        dispatch(wsDisconnect());
      } catch (error) {
        console.info('occur error while login');
      }
    };

    const fetchUserData = async () => {
      if (!publicKey) throw new Error('Wallet not connected!');

      // Fetch User Nft Info
      fetchUserNfts();
      // Fetch User Statistic Info
      fetchUserStatistics();
    };

    if (connected && publicKey && !isAuthenticated && request) handleLogin();
    if (connected && publicKey && isAuthenticated) {
      fetchUserData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [connected, publicKey, isAuthenticated, request]);

  const login = useCallback(async () => {
    if (!connected) setVisible(true);

    dispatch(requestLogin(true));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [connected]);

  const logout = useCallback(async () => {
    try {
      dispatch(requestLogin(false));

      await api.get('/user/logout');
    } catch (error) {
      console.info(error);
    } finally {
      disconnect();
      // console.log('public key', publicKey);

      dispatch(logoutUser());
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { login, logout };
};

export default useAuth;
