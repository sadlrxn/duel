import {
  createContext,
  FC,
  PropsWithChildren,
  useState,
  useRef,
  useEffect,
  useCallback
} from 'react';
import { useLocation } from 'react-router-dom';
import { shallowEqual } from 'react-redux';

import {
  useModal,
  AuthorizeModal,
  JurisdictionsModal,
  PASSWORD
} from 'components';

import { useAppSelector, useAppDispatch } from 'state';
import { sendMessage, wsConnect } from 'state/socket';
import { setTokenPrices } from 'state/token/action';
import * as coinflipActions from 'state/coinflip/actions';
import * as jackpotActions from 'state/jackpot/actions';
import * as grandJpActions from 'state/grandJackpot/actions';
import * as crashActions from 'state/crash/actions';
// import * as dreamActions from 'state/dreamtower/actions';
// import * as plinkoActions from 'state/plinko/actions';

import { useTokenPrices } from 'hooks';
import { fetchServerConfig, fetchTokens } from 'utils/fetchData';

export const MainContext = createContext<any>({});

export const MainProvider: FC<PropsWithChildren> = ({ children }) => {
  const { pathname } = useLocation();
  const dispatch = useAppDispatch();

  const connected = useAppSelector(
    state => state.socket.connected,
    shallowEqual
  );
  const code = useAppSelector(state => state.user.code, shallowEqual);
  const config = useAppSelector(state => state.user.config, shallowEqual);

  const { prices } = useTokenPrices();

  const [authorized, setAuthorized] = useState(false);

  const [onPresentModal] = useModal(
    <AuthorizeModal setAuthorized={setAuthorized} />,
    false
  );

  const [onJurisdictionsModal] = useModal(
    <JurisdictionsModal hideCloseButton />,
    false
  );

  const convertChipToBalance = useCallback(
    (amount: number) => {
      return Math.floor(amount * 10 ** 2) * 10 ** (config.balanceDecimals - 2);
    },
    [config.balanceDecimals]
  );

  const convertBalanceToChip = useCallback(
    (amount: number) => {
      return Math.floor(amount / 10 ** (config.balanceDecimals - 2)) / 10 ** 2;
    },
    [config.balanceDecimals]
  );

  useEffect(() => {
    fetchServerConfig();
    fetchTokens();
    dispatch(wsConnect());
  }, [dispatch]);

  useEffect(() => {
    if (!prices) return;
    dispatch(setTokenPrices(prices));
  }, [dispatch, prices]);

  useEffect(() => {
    if (code) {
      onJurisdictionsModal();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [code]);

  useEffect(() => {
    if (!PASSWORD) return;
    if (localStorage.getItem('duel') === PASSWORD) {
      setAuthorized(true);
      return;
    }
    // if (window.location.host !== 'duel.win') return;
    onPresentModal();

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authorized]);

  const connectedRef = useRef<boolean>(false);
  useEffect(() => {
    if (!connected) {
      dispatch(coinflipActions.reset());
      dispatch(jackpotActions.reset());
      dispatch(crashActions.reset());
      return;
    }
    if (connected && connectedRef.current === false) {
      dispatch(grandJpActions.setFetch(false));
      dispatch(sendMessage({ type: 'visit', room: 'grandJackpot' }));
    }
    connectedRef.current = connected;
    if (pathname.includes('coinflip')) {
      dispatch(coinflipActions.reset());
      dispatch(coinflipActions.setFetch(false));
      dispatch(sendMessage({ type: 'visit', room: 'coinflip' }));
    } else if (pathname.includes('jackpot')) {
      dispatch(jackpotActions.reset());
      // dispatch(grandJackpotReset());
      if (pathname.includes('grandjackpot')) {
        dispatch(grandJpActions.setFetch(false));
        dispatch(sendMessage({ type: 'visit', room: 'grandJackpot' }));
      } else {
        dispatch(jackpotActions.setFetch({ room: 'low', fetch: false }));
        dispatch(jackpotActions.setFetch({ room: 'medium', fetch: false }));
        dispatch(jackpotActions.setFetch({ room: 'wild', fetch: false }));
        dispatch(sendMessage({ type: 'visit', room: 'jackpot' }));
      }
    } else if (pathname.includes('crash')) {
      dispatch(crashActions.reset());
      dispatch(sendMessage({ type: 'visit', room: 'crash' }));
    } else {
      dispatch(sendMessage({ type: 'visit', room: '' }));
    }
  }, [pathname, connected, dispatch]);

  return (
    <MainContext.Provider
      value={{
        convertBalanceToChip,
        convertChipToBalance
      }}
    >
      {children}
    </MainContext.Provider>
  );
};
