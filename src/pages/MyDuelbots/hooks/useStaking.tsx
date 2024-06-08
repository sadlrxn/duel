import { useFetchUserInfo } from 'hooks';
import { useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useAppDispatch } from 'state';
import { cancelSelectedBots, loadBots } from 'state/staking/actions';
import { updateBalance } from 'state/user/actions';
import useSWR, { useSWRConfig } from 'swr';
import api from 'utils/api';

const useStaking = () => {
  const dispatch = useAppDispatch();
  const { fetchUserNfts } = useFetchUserInfo();
  const { mutate } = useSWRConfig();

  const stakeDuelBots = useCallback(async (mints: string[]) => {
    try {
      await api.post('/bot/stake', { mints });
      toast.success('Successfully Staked!');
      mutate(`/bot/duel-bots`, (data: any) => {
        dispatch(loadBots(data));
      });

      fetchUserNfts();
      dispatch(cancelSelectedBots());
    } catch (err: any) {
      if (err.response.status === 429) {
        toast.error(err.response.data.message);
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else console.log('Internal Server Err!');
    }
  }, []);

  const unStakeDuelBots = useCallback(async (mints: string[]) => {
    try {
      const { data: reward } = await api.post('/bot/unstake', { mints });
      toast.success('Successfully Unstaked!');
      mutate(`/bot/duel-bots`, (data: any) => {
        dispatch(loadBots(data));
      });

      fetchUserNfts();
      dispatch(updateBalance({ type: 1, usdAmount: reward }));
      dispatch(cancelSelectedBots());
    } catch (err: any) {
      if (err.response.status === 429) {
        toast.error(err.response.data.message);
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else console.log('Internal Server Err!');
    }
  }, []);

  const claimRewards = useCallback(async (mints: string[]) => {
    try {
      const { data: reward } = await api.post('/bot/claim', { mints });

      toast.success('Successfully Claimed!');
      mutate(`/bot/duel-bots`, (data: any) => {
        dispatch(loadBots(data));
      });

      dispatch(updateBalance({ type: 1, usdAmount: reward }));
      dispatch(cancelSelectedBots());
    } catch (err: any) {
      if (err.response.status === 429) {
        toast.error(err.response.data.message);
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else console.log('Internal Server Err!');
    }
  }, []);

  return { stakeDuelBots, unStakeDuelBots, claimRewards };
};

export default useStaking;
