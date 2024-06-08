import React from 'react';
import styled from 'styled-components';
import { TransitionGroup, CSSTransition } from 'react-transition-group';

import { useAppSelector } from 'state';

import Item from './Item';

export default function History() {
  const history = useAppSelector(state => state.plinko.history);

  return (
    <Container>
      {history.map(item => {
        return (
          <CSSTransition
            key={'crash_history_' + item.roundId}
            timeout={{
              enter: 2000,
              exit: 2000
            }}
            unmountOnExit
            mountOnEnter
            classNames="crash-history"
          >
            <Item
              key={item.roundId}
              multiplier={item.multiplier}
              roundId={item.roundId}
            />
          </CSSTransition>
        );
      })}
    </Container>
  );
}
const Container = styled(TransitionGroup)`
  display: flex;
  max-width: max-content;
  min-width: 800px;
  height: 40px;
  /* background: linear-gradient(
    90deg,
    rgba(9, 71, 117, 0.4) 0%,
    rgba(0, 148, 255, 0) 100%
  ); */
  border-radius: 9px;
  padding: 5px 6px;
  gap: 7px;

  .crash-history-enter {
    max-width: 0px;
    opacity: 0;
  }

  .crash-history-enter-active {
    max-width: 100px;
    opacity: 1;
    transition: max-width 1s, opacity 1s 1s;
  }

  .crash-history-exit {
    max-width: 100px;
    opacity: 1;
  }

  .crash-history-exit-active {
    max-width: 0px;
    opacity: 0;
    transition: opacity 1s, max-width 1s 1s;
  }
`;
