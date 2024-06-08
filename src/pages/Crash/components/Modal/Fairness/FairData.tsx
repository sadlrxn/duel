import React from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

import { Box, Flex } from 'components';

import { CrashFairData } from 'api/types/crash';
import { formatUserName } from 'utils/format';

import { GrayButton, Table } from './styles';

import Row from './Row';

interface FairDataProps {
  roundId?: number;
  onDismiss?: any;
  game: CrashFairData;
}

const FairData: React.FC<FairDataProps> = ({ onDismiss, game }) => {
  return (
    <>
      <TableWrapper
        overflow="auto"
        minHeight="260px"
        height="270px"
        mt="25px"
        mb="26px"
        mr="10px"
        pl="2px"
      >
        <Table>
          <thead>
            <tr>
              <th>Player</th>
              <th>CHIPS Bet Amount</th>
              <th>CHIPS Cashed Out</th>
              <th>Multiplier</th>
            </tr>
          </thead>
          <tbody>
            {game.bets.length > 0 &&
              game.bets
                .sort((a, b) => (a.profit! > b.profit! ? -1 : 0))
                .map(player => {
                  return (
                    <Row
                      name={formatUserName(player.user.name)}
                      betAmount={player.betAmount}
                      profit={player.profit}
                      paidBalanceType={player.paidBalanceType}
                      multiplier={
                        player.payoutMultiplier
                          ? String(player.payoutMultiplier) + 'x'
                          : '-'
                      }
                      key={player.betId}
                    />
                  );
                })}
          </tbody>
        </Table>
      </TableWrapper>

      <Flex gap={20} flexWrap="wrap">
        <GrayButton width={['100%', '100%', '192px']} onClick={onDismiss}>
          <Link to="/fair">Provably Fair</Link>
        </GrayButton>
        <GrayButton width={['100%', '100%', '216px']} onClick={onDismiss}>
          <Link
            to="/fair"
            state={{ gameType: 'crash', gameData: { serverSeed: game.seed } }}
          >
            Advanced Verification
          </Link>
        </GrayButton>
      </Flex>
    </>
  );
};

export default React.memo(FairData);

const TableWrapper = styled(Box)`
  ::-webkit-scrollbar-track {
    margin-top: 37px;
  }

  ::-webkit-scrollbar-corner {
    background: transparent;
  }
`;
