import { useMemo } from 'react';

import { ReactComponent as PersonIcon } from 'assets/imgs/icons/person.svg';

import { Flex, Span } from 'components';
import { useAppSelector } from 'state';
import { User } from 'api/types/user';

export default function Persons() {
  const bets = useAppSelector(state => state.crash.bets);

  const playerLength = useMemo(() => {
    return bets.reduce((players: User[], bet) => {
      const index = players.findIndex(player => player.id === bet.user.id);
      if (index === -1) players.push(bet.user);
      return players;
    }, []).length;
  }, [bets]);

  return (
    <Flex gap={3} alignItems="center">
      <Span color="#B2D1FF" fontWeight={700} fontSize="12px" lineHeight={1}>
        {playerLength}
      </Span>
      <PersonIcon stroke="#B2D1FF" width={10} height={11} />
    </Flex>
  );
}
