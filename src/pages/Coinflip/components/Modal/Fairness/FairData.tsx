import React, { useMemo } from 'react';
import { Link } from 'react-router-dom';

import { Box, Flex, Grid, Span } from 'components';

import { initialFair, CoinflipFairData } from 'api/types/coinflip';
import { useFetchRoundInfo } from 'hooks';
import { formatUserName } from 'utils/format';

import {
  Title,
  Description,
  Divider,
  VerifyButton,
  GrayButton
} from './styles';

import Detail from './Detail';
import { convertBalanceToChip } from 'utils/balance';

interface FairDataProps {
  roundId?: number;
  onDismiss?: any;
}

const FairData: React.FC<FairDataProps> = ({ roundId = 0, onDismiss }) => {
  const { data, error } = useFetchRoundInfo('coinflip', roundId);

  const { gameData } = useMemo(() => {
    const isError = error ? true : false;
    const isLoading = !isError && !data ? true : false;
    const gameData =
      !isError && !isLoading ? (data as CoinflipFairData) : initialFair;
    return { gameData };
  }, [data, error]);

  const {
    signedString,
    ticketId,
    headsUser,
    tailsUser,
    creatorId,
    winnerId: wId,
    amount
  } = gameData;

  const { creator, opponent, creatorSide, opponentSide, winner, winnerSide } =
    useMemo(() => {
      const winnerId = wId ?? 0;
      let creatorSide = 'Green';
      let opponentSide = 'Purple';
      let creator = headsUser;
      let opponent = tailsUser;
      if (headsUser?.id !== creatorId) {
        creatorSide = 'Purple';
        opponentSide = 'Green';
        creator = tailsUser;
        opponent = headsUser;
      }
      const winner = headsUser?.id === winnerId ? headsUser : tailsUser;
      const winnerSide =
        winnerId === 0 ? '' : headsUser?.id === winnerId ? 'Green' : 'Purple';
      return {
        creatorSide,
        creator,
        opponent,
        opponentSide,
        winner,
        winnerSide
      };
    }, [headsUser, tailsUser, creatorId, wId]);

  return (
    <>
      <Grid
        mt="14px"
        gridTemplateColumns="repeat(auto-fill, minmax(280px, 1fr))"
        gridColumnGap="30px"
        gridRowGap="17px"
      >
        <Detail
          title="Game Creator"
          text={formatUserName(creator ? creator.name : '')}
          readOnly
        />
        <Detail
          title="Creator's Side"
          text={creatorSide}
          side={creatorSide === 'Green' ? 'duel' : 'ana'}
          readOnly
        />
        <Detail
          title="Opponent"
          text={formatUserName(opponent ? opponent.name : '')}
          readOnly
        />
        <Detail
          title="Opponent's Side"
          text={opponentSide}
          side={opponentSide === 'Green' ? 'duel' : 'ana'}
          readOnly
        />
      </Grid>

      <Divider />

      <Box>
        <Title>True Randomness</Title>
        <Grid
          gridTemplateColumns={['auto', 'auto', 'auto', 'auto max-content']}
          gridGap={[14, 14, 14, 90]}
          mt="10px"
        >
          <Description>
            When a new game is created, it is assigned a
            <Span fontWeight={700}> Ticket ID</Span> by Random.org which is used
            to generate a truly
            <Span fontWeight={700}> Random String</Span> after the game has
            started.
          </Description>
          <VerifyButton disabled={ticketId === ''}>
            <a
              href={`https://api.random.org/tickets/form?ticket=${ticketId}`}
              rel="noreferrer"
              target={'_blank'}
            >
              Verify Randomness
            </a>
          </VerifyButton>
        </Grid>
      </Box>

      <Grid
        mt="14px"
        gridTemplateColumns="repeat(auto-fill, minmax(280px, 1fr))"
        gridColumnGap="30px"
        gridRowGap="17px"
      >
        <Detail title="Ticket ID" text={ticketId} enableCopy readOnly />
        <Detail
          title="Random String"
          text={signedString}
          placeholder="Game In Progress..."
          enableCopy
          readOnly
        />
      </Grid>

      <Divider />

      <Box>
        <Title>Winner Outcome</Title>
        <Grid
          gridTemplateColumns={['auto', 'auto', 'auto', 'auto max-content']}
          gridGap={[14, 14, 14, 90]}
          mt="10px"
        >
          <Description>
            The outcome is generated using the
            <Span fontWeight={700}> Random String</Span>. When the game is
            complete you can verify the outcome of the game by tapping
            <Span fontWeight={700}> Verify Outcome</Span>.
          </Description>
          <VerifyButton disabled={signedString === ''} onClick={onDismiss}>
            <Link to="/fair" state={{ gameType: 'coinflip', gameData }}>
              {signedString === '' ? 'Game In Progress' : 'Verify Outcome'}
            </Link>
          </VerifyButton>
        </Grid>
      </Box>

      <Grid
        mt="14px"
        gridTemplateColumns="repeat(auto-fill, minmax(250px, 1fr))"
        gridColumnGap="25px"
        gridRowGap="17px"
      >
        <Detail
          title="Winning Side"
          text={winnerSide}
          side={
            winnerSide === 'Green'
              ? 'duel'
              : winnerSide === 'Purple'
              ? 'ana'
              : ''
          }
          readOnly
        />
        <Detail title="Winner" text={winner ? winner.name : ''} readOnly />
        <Detail
          title="Total Prize Value*"
          text={convertBalanceToChip(2 * amount).toString()}
          readOnly
        />
      </Grid>

      <Divider />

      <Flex gap={30} flexWrap="wrap">
        <GrayButton width={['100%', '100%', '192px']} onClick={onDismiss}>
          <Link to="/fair">Provably Fair</Link>
        </GrayButton>
        <GrayButton width={['100%', '100%', '216px']} onClick={onDismiss}>
          <Link to="/fair" state={{ gameType: 'coinflip', gameData }}>
            Advanced Verification
          </Link>
        </GrayButton>
      </Flex>
    </>
  );
};

export default React.memo(FairData);
