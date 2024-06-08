import React, { useCallback, useState, useEffect } from 'react';
import { TransitionGroup, CSSTransition } from 'react-transition-group';
import ClipLoader from 'react-spinners/ClipLoader';
import { Options } from 'react-select';

import state, { useAppSelector } from 'state';
import { api } from 'services';
import * as coinflipActions from 'state/coinflip/actions';
import { Button, CreatableSelect as Select } from 'components';
import { createOption } from 'utils/selectOption';

import { Row as GameRow } from '../Row';

import {
  Container,
  GameList,
  DataContainer,
  Heading,
  Text
} from './List.styles';

const COINFLIP_OPTIONS: Options<{ label: string; value: string }> = [
  { label: 'All Games', value: 'All Games' }
];

function CoinflipList() {
  const game = useAppSelector(state => state.coinflip);
  const { name: userName, role: userRole } = useAppSelector(
    state => state.user
  );
  const differ = useAppSelector(state => state.socket.differ);

  const { history, games } = game;

  const [isFetching, setIsFetching] = useState(false);
  const [options, setOptions] = useState(COINFLIP_OPTIONS);
  const [option, setOption] = useState(COINFLIP_OPTIONS[0]);

  useEffect(() => {
    if (userName)
      setOptions([COINFLIP_OPTIONS[0], createOption('My Games', userName)]);
    else setOptions([COINFLIP_OPTIONS[0]]);
  }, [userName]);

  useEffect(() => {
    const option = options.find(op => op.value === history.winner);
    if (option) setOption(option);
  }, [history.winner, options]);

  const handleCreate = useCallback((inputValue: string) => {
    const newOption = createOption(inputValue);
    setOptions(prev => [...prev, newOption]);
    setOption(newOption);
  }, []);

  const fetchHistory = useCallback(
    async (offset: number) => {
      setIsFetching(true);
      try {
        const { data } = await api.get(`/coinflip/history`, {
          params: {
            userName:
              option.value === COINFLIP_OPTIONS[0].value
                ? undefined
                : option.value,
            offset,
            count: 6
          }
        });
        const newHistory = data.history
          //@ts-ignore
          .map(item => ({
            ...item,
            status: 'ended',
            time: item.endedAt
              ? new Date(item.endedAt).getTime() + differ
              : Date.now()
          }))
          .filter(
            //@ts-ignore
            h =>
              history.games
                .slice(0, offset)
                .findIndex(item => item.roundId === h.roundId) === -1
          );
        state.dispatch(
          coinflipActions.setHistory([
            ...history.games.slice(0, offset),
            ...newHistory
          ])
        );
      } catch {}
      setIsFetching(false);
    },
    [history, option, differ]
  );

  useEffect(() => {
    fetchHistory(0);
    state.dispatch(coinflipActions.setHistoryWinner(option.value));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [option]);

  const handleShowMore = useCallback(async () => {
    fetchHistory(history.games.length);
  }, [fetchHistory, history.games.length]);

  return (
    <Container>
      <Heading>
        <DataContainer justifyContent="space-between" width="100%">
          <Text>Active Games</Text>
          <Select
            // background="#03060999"
            // hoverBackground="#29364a"
            // color="#ffffff"
            isSearchable={userRole !== 'user'}
            onCreateOption={handleCreate}
            width={200}
            isDisabled={isFetching}
            background="#192637"
            hoverBackground="#03060933"
            color="#B2D1FF"
            options={options}
            onChange={(selectedOption: any) => setOption(selectedOption)}
            value={option}
          />
        </DataContainer>
      </Heading>
      <GameList>
        <TransitionGroup appear>
          {games
            .filter(game => {
              const names: string[] = [];
              if (game.headsUser) names.push(game.headsUser.name);
              if (game.tailsUser) names.push(game.tailsUser.name);
              if (
                option.value === 'All Games' ||
                names.indexOf(option.value) !== -1
              )
                return true;
              return false;
            })
            .map(game => {
              return (
                <CSSTransition key={game.roundId} timeout={500}>
                  <GameRow game={game} />
                </CSSTransition>
              );
            })}
        </TransitionGroup>
      </GameList>
      <Heading mt="40px">
        <DataContainer justifyContent="space-between" width="100%">
          <Text>Recent Games</Text>
        </DataContainer>
      </Heading>

      <GameList>
        <TransitionGroup>
          {history.games.map(game => {
            return (
              <CSSTransition key={game.roundId} timeout={500}>
                <GameRow game={game} />
              </CSSTransition>
            );
          })}
        </TransitionGroup>
      </GameList>

      <Button
        variant="secondary"
        outlined
        scale="sm"
        width={153}
        background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
        color="#FFFFFF"
        borderColor="chipSecondary"
        marginX="auto"
        marginTop={20}
        onClick={isFetching ? undefined : handleShowMore}
      >
        {isFetching ? <ClipLoader size={20} color="#fff" /> : 'SHOW MORE'}
      </Button>
    </Container>
  );
}

export default React.memo(CoinflipList);
