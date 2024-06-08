import { useCallback, useState, useEffect } from 'react';
import { ClipLoader } from 'react-spinners';
import { Options } from 'react-select';

import { Flex, Text, CreatableSelect as Select, Button } from 'components';
import state, { useAppSelector } from 'state';
import * as dreamtowerActions from 'state/dreamtower/actions';
import { DreamtowerHistoryRound } from 'api/types/dreamtower';
import { api } from 'services';
import { createOption } from 'utils/selectOption';

import {
  HistoryList,
  HistoryListContainer
} from '../components/History/History.styles';
import { HistoryItem } from '../components/History';

const HISTORY_OPTIONS: Options<{ label: string; value: string }> = [
  { label: 'All Games', value: 'All Games' }
];

export default function History() {
  const { name: userName, role: userRole } = useAppSelector(
    state => state.user
  );
  const history = useAppSelector(state => state.dreamtower.history);

  const [options, setOptions] = useState(HISTORY_OPTIONS);
  const [option, setOption] = useState(HISTORY_OPTIONS[0]);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    if (userName)
      setOptions([HISTORY_OPTIONS[0], createOption('My Games', userName)]);
    else setOptions([HISTORY_OPTIONS[0]]);
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
        const { data } = await api.get(`/dreamtower/history`, {
          params: {
            userName:
              option.value === HISTORY_OPTIONS[0].value
                ? undefined
                : option.value,
            offset,
            count: 10
          }
        });
        const newHistory = data.history.filter(
          //@ts-ignore
          h =>
            history.rounds
              .slice(0, offset)
              .findIndex(item => item.roundId === h.roundId) === -1
        );
        state.dispatch(
          dreamtowerActions.setHistory([
            ...history.rounds.slice(0, offset),
            ...newHistory
          ])
        );
      } catch {}
      setIsFetching(false);
    },
    [history.rounds, option]
  );

  const handleSelectReview = (round: DreamtowerHistoryRound) => {
    state.dispatch(dreamtowerActions.review(round));
  };

  const handleShowMore = useCallback(async () => {
    fetchHistory(history.rounds.length);
  }, [fetchHistory, history.rounds.length]);

  useEffect(() => {
    fetchHistory(0);
    state.dispatch(dreamtowerActions.setHistoryWinner(option.value));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [option]);

  return (
    <div className="container">
      <Flex justifyContent={'space-between'} marginBottom="10px">
        <Text color="white" fontSize={'25px'} fontWeight={500}>
          Game History
        </Text>

        <Select
          isSearchable={userRole !== 'user'}
          onCreateOption={handleCreate}
          width={170}
          isDisabled={isFetching}
          background="#192637"
          hoverBackground="#03060933"
          color="#B2D1FF"
          options={options}
          onChange={(selectedOption: any) => setOption(selectedOption)}
          value={option}
        />
      </Flex>
      {history.rounds.length > 0 && (
        <HistoryListContainer>
          <HistoryList>
            {history.rounds.map(round => {
              return (
                <HistoryItem
                  selected={
                    history.review !== undefined &&
                    round.roundId === history.review.roundId
                  }
                  round={round}
                  onClick={() => handleSelectReview(round)}
                  key={round.roundId}
                  background={
                    'linear-gradient(90deg, #182738 0%, #182738 100%)'
                  }
                />
              );
            })}
          </HistoryList>
        </HistoryListContainer>
      )}
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
        onClick={isFetching ? null : handleShowMore}
      >
        {isFetching ? <ClipLoader size={20} color="#fff" /> : 'SHOW MORE'}
      </Button>
    </div>
  );
}
