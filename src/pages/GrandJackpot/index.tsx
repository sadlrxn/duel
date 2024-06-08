import React from 'react';
import { Flex, PageSpinner, Topbar } from 'components';

import { useAppSelector } from 'state';

import GrandJackpotGame from './GrandJackpotGame';
import History from 'pages/Jackpot/History';

export default function Jackpot() {
  const loading = useAppSelector(state => !state.grandJackpot.fetch);
  const { fee } = useAppSelector(state => state.meta.grandJackpot);

  if (loading) return <PageSpinner />;

  return (
    <Flex
      padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}
      flexDirection="column"
      gap={36}
    >
      <Topbar title={'24 HOUR GRAND JACKPOT'} fee={fee} />
      <GrandJackpotGame />
      <History isGrand />
    </Flex>
  );
}
