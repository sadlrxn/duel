import React, { useMemo } from 'react';
import { Link } from 'react-router-dom';

import { Box, Flex, Grid, Span } from 'components';

import { initialFair, JackpotFairData } from 'api/types/jackpot';
import { useFetchRoundInfo } from 'hooks';
import { formatUserName } from 'utils/format';

import {
  Title,
  Description,
  Divider,
  VerifyButton,
  GrayButton,
  Table
} from './styles';

import Detail from './Detail';
import Row from './Row';
import { convertBalanceToChip } from 'utils/balance';

interface FairDataProps {
  roundId?: number;
  onDismiss?: any;
}

const FairData: React.FC<FairDataProps> = ({ roundId = 0, onDismiss }) => {
  const { data, error } = useFetchRoundInfo('jackpot', roundId);

  const { gameData } = useMemo(() => {
    const isError = error ? true : false;
    const isLoading = !isError && !data ? true : false;
    const gameData =
      !isError && !isLoading ? (data as JackpotFairData) : initialFair;
    return { gameData };
  }, [data, error]);

  const { ticketId, signedString, players, winner } = gameData;

  const { percents, wIndex, totalBetAmount } = useMemo(() => {
    const totalBetAmount = players.reduce(
      (sum: number, player) => sum + player.usdAmount + (player.nftAmount ?? 0),
      0
    );
    const percents = players.map(player => {
      return +(
        ((player.usdAmount + (player.nftAmount ?? 0)) / totalBetAmount) *
        100
      ).toFixed(2);
    });
    const wIndex = players.findIndex(player => player.id === winner.id);
    return { percents, wIndex, totalBetAmount };
  }, [players, winner.id]);

  return (
    <>
      <Box overflow="auto" height="150px">
        <Table>
          <thead>
            <tr>
              <th>Player</th>
              <th>CHIPS Bet Amount</th>
              <th>NFT Bet Amount</th>
              <th>Win Chance</th>
            </tr>
          </thead>
          <tbody>
            {players.map((player, index) => {
              return (
                <Row
                  name={formatUserName(player.name)}
                  usdAmount={convertBalanceToChip(player.usdAmount)}
                  nftAmount={convertBalanceToChip(player.nftAmount ?? 0)}
                  chance={percents[index]}
                  key={player.id}
                />
              );
            })}
          </tbody>
        </Table>
      </Box>

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
            <Link to="/fair" state={{ gameType: 'jackpot', gameData }}>
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
          title="Win Chance"
          text={wIndex !== -1 ? `${percents[wIndex].toString()}%` : ''}
          readOnly
        />
        <Detail
          title="Winner"
          text={wIndex !== -1 ? players[wIndex].name : ''}
          readOnly
        />
        <Detail
          title="Total Prize Value*"
          text={convertBalanceToChip(totalBetAmount).toString()}
          readOnly
          showChip
        />
      </Grid>

      <Divider />

      <Flex gap={30} flexWrap="wrap">
        <GrayButton width={['100%', '100%', '192px']} onClick={onDismiss}>
          <Link to="/fair">Provably Fair</Link>
        </GrayButton>
        <GrayButton width={['100%', '100%', '216px']} onClick={onDismiss}>
          <Link to="/fair" state={{ gameType: 'jackpot', gameData }}>
            Advanced Verification
          </Link>
        </GrayButton>
      </Flex>
    </>
  );
};

export default React.memo(FairData);
