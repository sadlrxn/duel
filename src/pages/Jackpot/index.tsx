import { useEffect } from 'react';
import { Flex, PageSpinner, Topbar } from 'components';

import { useAppDispatch, useAppSelector } from 'state';
import { setBalanceType } from 'state/user/actions';

import History from './History';
import JackpotGame from './JackpotGame';

export default function Jackpot() {
  const dispatch = useAppDispatch();
  const lowLoading = useAppSelector(state => !state.jackpot.low.fetch);
  const wildLoading = useAppSelector(state => !state.jackpot.wild.fetch);

  useEffect(() => {
    dispatch(setBalanceType('chip'));
  }, [dispatch]);

  const { fee } = useAppSelector(
    state => state.meta.jackpot[state.jackpot.room]
  );

  if (lowLoading || wildLoading) return <PageSpinner />;

  return (
    <Flex
      flexDirection="column"
      gap={36}
      padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}
    >
      <Topbar title="JACKPOT" fee={fee} />
      <JackpotGame />
      <History />
    </Flex>
  );
}
