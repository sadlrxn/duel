import { ComponentMeta } from '@storybook/react';
import { CoinflipRoundData as Game } from 'api/types/coinflip';

import CoinflipRow from './Row';

export default {
  title: 'Coinflip/CoinflipRow',
  component: CoinflipRow,
  argTypes: {}
} as ComponentMeta<typeof CoinflipRow>;

export const Primary = () => {
  const game: Game = {
    status: 'created',
    roundId: 100,
    headsUser: {
      id: 100,
      name: 'ichiro',
      avatar: 'URL for avatar image',
      count: 1
    },
    tailsUser: null,
    amount: 100,
    prize: 194,
    ticketId: '5158795321',
    signedString: '',
    winnerId: 0,
    creatorId: 100,
    paidBalanceType: 'chip',
    time: Date.now(),
    request: false
  };

  return (
    <>
      <CoinflipRow game={game} />
    </>
  );
};
