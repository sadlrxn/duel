import { useRef } from 'react';
import { shallowEqual } from 'react-redux';
import styled from 'styled-components';
import { TransitionGroup } from 'react-transition-group';

import { Flex } from 'components';
import { useAppSelector } from 'state';

import CashoutUser from './CashoutUser';
import CrashMain from './CrashMain';

interface CrashGameProps {
  gameHeight?: string;
}

export default function CrashGame({ gameHeight }: CrashGameProps) {
  const cashOuts = useAppSelector(state => state.crash.cashOuts, shallowEqual);
  const usingAnimation = useAppSelector(
    state => state.user.crashAnimation,
    shallowEqual
  );

  const graphRef = useRef<HTMLDivElement>(null);

  return (
    <Container ref={graphRef} height={gameHeight}>
      <TransitionGroup
        className="crash_cashout_users"
        style={{
          position: 'absolute',
          left: '0',
          top: '0',
          width: '100%',
          height: '100%'
        }}
      >
        {usingAnimation &&
          cashOuts.map(cashOut => {
            return (
              <CashoutUser
                key={'crash_cashout_' + cashOut.betId}
                cashOut={cashOut}
              />
            );
          })}
      </TransitionGroup>
      <CrashMain graphRef={graphRef} />
    </Container>
  );
}

const Container = styled(Flex)`
  position: relative;
  width: 100%;

  margin-bottom: 74px;

  .crash_cashout_user-exit {
    opacity: 1;
  }

  .crash_cashout_user-exit-active {
    opacity: 0;
  }
`;
