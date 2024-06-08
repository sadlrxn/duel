import { useCallback, useState, useEffect, useMemo } from 'react';
import styled, { css } from 'styled-components';
import { margin, MarginProps } from 'styled-system';
import dayjs from 'dayjs';
import ClipLoader from 'react-spinners/ClipLoader';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import { useNavigate } from 'react-router-dom';
import { Options } from 'react-select';

import { api } from 'services';
import state, { useAppSelector } from 'state';
import * as jackpotActions from 'state/jackpot/actions';
import * as grandJackpotActions from 'state/grandJackpot/actions';
import { imageProxy } from 'config';

import {
  Box,
  Text,
  Flex,
  Avatar,
  Badge,
  Chip,
  Button,
  useModal,
  CreatableSelect as Select,
  Notification
} from 'components';
import { VerifyButton, FairnessModal } from './components';

import { formatUserName } from 'utils/format';
import { createOption } from 'utils/selectOption';
import { User } from 'api/types/jackpot';
import { StyledNotification } from './components/UserStatus/styles';
import { useQuery } from 'hooks';
import { convertBalanceToChip } from 'utils/balance';

const JACKPOT_OPTIONS: Options<{ label: string; value: string }> = [
  { label: 'All Games', value: 'All Games' }
];

interface WinnerProps extends MarginProps {
  avatar?: string;
  level: number;
  name: string;
  userId?: number;
}

function Winner({ userId, avatar, name, ...props }: WinnerProps) {
  return (
    <WinnerWrapper {...props}>
      <Avatar
        userId={userId}
        name={name}
        image={avatar}
        border="none"
        borderRadius="5px"
        padding="0px"
        size="36px"
      />
      {formatUserName(name)}
    </WinnerWrapper>
  );
}

interface RowProps {
  roundId: number;
  ticketId: string;
  players: User[];
  signedString: string;
  winner: User;
  chance: number;
  prize: number;
  time: number;
  userId: number;
  userName: string;
  isGrand?: boolean;
}

const Row = ({
  roundId,
  players,
  winner,
  chance,
  prize,
  time,
  userId,
  userName,
  isGrand
}: RowProps) => {
  const [onFairnessModal] = useModal(<FairnessModal roundId={roundId} />, true);
  const navigate = useNavigate();
  const query = useQuery();
  const [roundIdFromLink, setRoundIdFromLink] = useState(0);

  useEffect(() => {
    const roundId = query.get('roundId');
    if (!roundId) setRoundIdFromLink(0);
    else setRoundIdFromLink(+roundId);
  }, [query]);

  const hanleHistoryClick = () => {
    return navigate(
      `/${isGrand ? 'grandjackpot' : 'jackpot'}?roundId=${roundId}`
    );
  };

  return (
    <Tr onClick={hanleHistoryClick} selected={roundId === roundIdFromLink}>
      <Td>#{roundId}</Td>
      <Td>{dayjs(time).format('MMM DD, hh:mm A')}</Td>
      <Td>
        <Flex
          flexDirection="row"
          gap={4}
          alignItems="center"
          marginLeft={20}
          marginRight={20}
          position="relative"
          height={30}
          justifyContent="start"
        >
          {players.length > 0 &&
            players
              .slice(0, Math.min(players.length, 3))
              .map((player, index) => (
                <StackedImage
                  key={roundId + '_' + player.id + '_' + index}
                  src={imageProxy() + player.avatar}
                  alt={'player'}
                />
              ))}
          {players.length > 3 && (
            <StyledNotification>{players.length - 3}</StyledNotification>
          )}
        </Flex>
      </Td>
      <Td>
        <Winner
          userId={winner.id}
          avatar={winner.avatar}
          level={5}
          name={winner.id === userId ? userName : winner.name}
          margin={2}
        />
      </Td>
      <Td>
        <Badge
          variant={chance < 7 ? 'secondary' : 'primary'}
          margin="auto"
          fontSize="14px"
        >
          {chance.toFixed(2)}%
          {chance < 7 && <Notification variant="secondary">Snipe</Notification>}
        </Badge>
      </Td>
      <Td>
        <Chip price={convertBalanceToChip(prize)} margin="auto" />
      </Td>
      <Td>
        <VerifyButton ml="auto" onClick={onFairnessModal} />
      </Td>
    </Tr>
  );
};

export interface HistoryProps {
  isGrand?: boolean;
}

export default function History({ isGrand = false }: HistoryProps) {
  const {
    id: userId,
    name: userName,
    role: userRole
  } = useAppSelector(state => state.user);
  const differ = useAppSelector(state => state.socket.differ);
  const grandJackpot = useAppSelector(state => state.grandJackpot);
  const jackpot = useAppSelector(state => state.jackpot);

  const history = useMemo(
    () => (isGrand ? grandJackpot.history : jackpot.history),
    [grandJackpot.history, isGrand, jackpot]
  );

  const [options, setOptions] = useState(JACKPOT_OPTIONS);
  const [option, setOption] = useState(JACKPOT_OPTIONS[0]);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    if (userName)
      setOptions([JACKPOT_OPTIONS[0], createOption('My Games', userName)]);
    else setOptions([JACKPOT_OPTIONS[0]]);
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
        const route = isGrand ? '/grand-jackpot/history' : `/jackpot/history`;

        const { data } = await api.get(route, {
          params: {
            userName:
              option.value === JACKPOT_OPTIONS[0].value
                ? undefined
                : option.value,
            offset,
            count: 10
          }
        });
        const newHistory = data.history
          //@ts-ignore
          .map(item => ({
            ...item,
            time: new Date(item.endedAt).getTime() + differ
          }))
          .filter(
            //@ts-ignore
            h =>
              history.games
                .slice(0, offset)
                .findIndex(item => item.roundId === h.roundId) === -1
          );
        if (isGrand)
          state.dispatch(
            grandJackpotActions.setHistory([
              ...history.games.slice(0, offset),
              ...newHistory
            ])
          );
        else
          state.dispatch(
            jackpotActions.setHistory([
              ...history.games.slice(0, offset),
              ...newHistory
            ])
          );
      } catch {}
      setIsFetching(false);
    },
    [isGrand, option.value, history.games, differ]
  );

  useEffect(() => {
    fetchHistory(0);
    if (isGrand)
      state.dispatch(grandJackpotActions.setHistoryWinner(option.value));
    else state.dispatch(jackpotActions.setHistoryWinner(option.value));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [option]);

  const handleShowMore = useCallback(async () => {
    fetchHistory(history.games.length);
  }, [fetchHistory, history.games.length]);

  return (
    <HistoryWrapper>
      <Flex mb={40} justifyContent="space-between" alignItems="center">
        <Text color="#fff" fontSize={25} fontWeight={600}>
          Game History
        </Text>
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
      </Flex>
      <TableWrapper>
        <Box p={10}>
          <StyledTable>
            <thead>
              <tr>
                <Td>GAME ID</Td>
                <Td>DATE</Td>
                <StyledTd>DUELERS</StyledTd>
                <StyledTd>WINNER</StyledTd>
                <Td>CHANCE</Td>
                <Td>JACKPOT</Td>
                <Td>
                  HISTORY
                  <br />
                  /FAIRNESS
                </Td>
              </tr>
            </thead>
            <tbody>
              {history.games.map(item => (
                <Row
                  isGrand={isGrand}
                  {...item}
                  userId={userId}
                  userName={userName}
                  key={item.roundId}
                />
              ))}
            </tbody>
          </StyledTable>
        </Box>
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
      </TableWrapper>
    </HistoryWrapper>
  );
}

const StyledTable = styled.table`
  width: 100%;
  min-width: 900px;
  border-collapse: separate;
  border-spacing: 0 10px;
  color: ${({ theme }) => theme.colors.textWhite};
  font-size: 14px;

  thead {
    padding: 4px;
    letter-spacing: 0.17em;
  }

  td {
    text-align: center;
    width: 15%;
  }

  tr td:first-child {
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    text-align: left;
    padding-left: 20px;
    width: 13%;
  }
  tr td:nth-child(2) {
    text-align: left;
  }
  tr td:last-child {
    border-top-right-radius: 10px;
    border-bottom-right-radius: 10px;
    text-align: right;
  }
`;

const Tr = styled.tr<{ selected: boolean }>`
  background: #182738;
  border-radius: 8px;
  color: ${({ theme }) => theme.colors.text};
  padding: 4px;
  margin: 4px;
  cursor: pointer;
  transition: 0.3s;
  &:hover {
    background: #263449;
  }
  td:first-child {
    padding: 15px;
  }
  td:last-child {
    padding: 15px;
  }
  ${({ selected }) => {
    if (selected)
      return css`
        border: 1px solid #ffe87f;
        box-shadow: 0px 0px 16px rgba(255, 232, 127, 0.75);
        td {
          border-top: 1px solid #ffe87f;
          border-bottom: 1px solid #ffe87f;
        }
        td:first-child {
          border-left: 1px solid #ffe87f;
        }
        td:last-child {
          border-right: 1px solid #ffe87f;
        }
      `;
  }}
`;

const Td = styled.td`
  margin: 4px;
  /* border: 2px solid transparent; */
`;

const HistoryWrapper = styled.div`
  width: 100%;
`;

const WinnerWrapper = styled(Flex)`
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 10px;
  float: left;
  transform: translateX(20px);
  ${margin}
`;

const TableWrapper = styled.div`
  background: linear-gradient(180deg, #0f1a26 0%, #0f1a26 0.01%, #0f1a26 100%);
  border-radius: 13px;
  padding-bottom: 35px;

  & > div {
    &::-webkit-scrollbar {
      height: 4px;
    }

    overflow: auto;
  }
`;

const StyledTd = styled(Td)`
  text-align: left !important;
  transform: translateX(22px);
`;

export const StackedImage = styled(LazyLoadImage)`
  & + & {
    margin-left: -20px;
  }

  width: 38px;
  height: 38px;

  min-width: 38px;
  min-height: 38px;

  border: 1px solid #0f1a26;
  border-radius: 5px;
`;
