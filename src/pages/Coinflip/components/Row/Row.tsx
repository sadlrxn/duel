import { useMemo, useState, useRef, FC, memo } from 'react';
import padlock from 'assets/imgs/icons/padlock.svg';

import { CoinflipRoundData as Game } from 'api/types/coinflip';
import { useAppSelector } from 'state';

import { Span, FlexProps, FairnessIcon, useModal } from 'components';
import { FairnessModal } from '../Modal';
import Middle from '../Middle/Middle';
import { Side } from '../Side';
import {
  Container,
  RowContainer,
  FairnessButton,
  MiddleContainer,
  PrivateGame
} from './Row.styles';

interface CoinflipRowProps extends FlexProps {
  game: Game;
}

const CoinflipRow: FC<CoinflipRowProps> = ({ game, ...props }) => {
  const user = useAppSelector(state => state.user);

  const rowRef = useRef(null);
  const [isPrivate] = useState(false);
  const [end, setEnd] = useState(false);

  const [onFairnessModal] = useModal(
    <FairnessModal roundId={game.roundId} />,
    true
  );

  const { head, tail, status, winner, time, request } = useMemo(() => {
    const status = game.status;
    const time = game.time;
    const request = game.request;

    const winner: 'duel' | 'ana' =
      game.winnerId === game.headsUser?.id ? 'duel' : 'ana';

    let head = {
      userId: game.headsUser ? game.headsUser.id : 0,
      name: game.headsUser
        ? game.headsUser.id === user.id
          ? user.name
          : game.headsUser.name
        : '',
      avatar: game.headsUser ? game.headsUser.avatar : '',
      prize: game.prize,
      winner: game.winnerId !== 0 && game.winnerId === game.headsUser?.id,
      user: user.id !== 0 && game.winnerId === user.id,
      paidBalanceType: game.paidBalanceType
    };

    let tail = {
      userId: game.tailsUser ? game.tailsUser.id : 0,
      name: game.tailsUser
        ? game.tailsUser.id === user.id
          ? user.name
          : game.tailsUser.name
        : '',
      avatar: game.tailsUser ? game.tailsUser.avatar : '',
      prize: game.prize,
      winner: game.winnerId !== 0 && game.winnerId === game.tailsUser?.id,
      user: user.id !== 0 && game.winnerId === user.id,
      paidBalanceType: game.paidBalanceType
    };

    return { head, tail, status, winner, time, request };
  }, [game, user]);

  const { creator, login, play } = useMemo(() => {
    const creator = game.creatorId !== 0 && game.creatorId === user.id;
    const login = user.id !== 0;
    const play =
      user.id !== 0 &&
      (user.id === game.headsUser?.id || user.id === game.tailsUser?.id);
    return { creator, login, play };
  }, [game, user]);

  return (
    <Container>
      <RowContainer ref={rowRef} {...props}>
        <Side side="duel" {...head} end={status === 'ended' || end} />
        <MiddleContainer>
          {isPrivate && (
            <PrivateGame>
              <img src={padlock} alt="" width={10} height={12} />
              <Span fontWeight={600}>PRIVATE GAME</Span>
            </PrivateGame>
          )}
          <Middle
            creator={creator}
            roundId={game.roundId}
            amount={game.amount}
            status={status}
            request={request}
            winner={winner}
            time={time}
            login={login}
            play={play}
            rowRef={rowRef}
            setEnd={setEnd}
          />
          {/* {status === 'ended' && (
            <>
              <Flex
                gap={19}
                flexWrap="wrap"
                fontSize={['10px', '10px', '10px', '12px']}
                lineHeight="1.3em"
                letterSpacing="0.1em"
                style={{
                  position: 'absolute',
                  left: '50%',
                  bottom: '0',
                  transform: 'translate(-60%, 10%)'
                }}
              >
                <Span color="white" fontWeight={400}>
                  #{game.roundId}
                </Span>
                <Span color="white" fontWeight={400}>
                  {dayjs(time).format('MMM DD, hh:mm A')}
                </Span>
              </Flex>
            </>
          )} */}
        </MiddleContainer>
        <Side side="ana" {...tail} end={status === 'ended' || end} />
      </RowContainer>

      <FairnessButton onClick={onFairnessModal}>
        <FairnessIcon size={14} />
      </FairnessButton>
    </Container>
  );
};

export default memo(CoinflipRow);
