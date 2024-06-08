import { useMemo } from 'react';
import styled from 'styled-components';

import { useCallback } from 'react';
import { Box, Flex, Chip, Span, Button } from 'components';
import { CrashBet } from 'api/types/crash';
import state, { useAppSelector } from 'state';
import { sendMessage } from 'state/socket';
import { setPending } from 'state/crash/actions';
import { convertBalanceToChip } from 'utils/balance';

const CashoutItem = ({
  bet,
  handleClick,
  multiplier,
  maxCashOut,
  show = false
}: {
  bet: CrashBet;
  handleClick: any;
  multiplier: number;
  maxCashOut: number;
  show?: boolean;
}) => {
  return (
    <CashoutItemWrapper>
      <Flex
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        width="100%"
        borderRadius="0px 7px 7px 0px"
        background="linear-gradient(180deg, rgba(19, 32, 49, 0.75) 0%, rgba(26, 41, 60, 0.75) 67.97%);"
        padding="8px 10px"
      >
        <Chip
          decimal={2}
          chipType={bet.paidBalanceType}
          price={convertBalanceToChip(
            Math.min(bet.betAmount * multiplier, maxCashOut)
          )}
          fontSize="1em"
          fontWeight={700}
          color="white"
        />
        <Span
          color="#fffc"
          fontSize="9px"
          fontWeight={500}
          lineHeight={1}
          my="2px"
        >
          {bet.cashOutAt ? 'Auto ' + bet.cashOutAt?.toFixed(2) + 'x' : ''}
        </Span>
        {show && (
          <Button
            padding="6px 8px"
            marginBottom="-16px"
            borderRadius="5px"
            fontSize="1em"
            fontWeight={700}
            color="black"
            onClick={() => handleClick(bet.betId)}
            disabled={bet.pending}
          >
            Cash Out
          </Button>
        )}
      </Flex>
    </CashoutItemWrapper>
  );
};

export default function Cashout({ multiplier = 1, ...props }: any) {
  const { roundId, bets, status } = useAppSelector(state => state.crash);
  const { maxCashOut } = useAppSelector(state => state.meta.crash);
  const user = useAppSelector(state => state.user);

  const handleClick = useCallback(
    (betId: number) => {
      state.dispatch(
        sendMessage({
          type: 'event',
          room: 'crash',
          content: JSON.stringify({
            type: 'cash-out',
            content: JSON.stringify({
              roundId,
              betId
            })
          })
        })
      );

      state.dispatch(setPending({ betId, status: true }));
    },
    [roundId]
  );

  const userBets = useMemo(() => {
    return bets
      .filter(bet => bet.user.id === user.id && !bet.payoutMultiplier)
      .sort((b1, b2) => (b2.cashOutAt ?? 1) - (b1.cashOutAt ?? 1));
  }, [bets, user.id]);

  return (
    <Container {...props}>
      {userBets.map(bet => {
        return (
          <CashoutItem
            key={bet.betId}
            bet={bet}
            multiplier={multiplier}
            maxCashOut={maxCashOut}
            handleClick={handleClick}
            show={status === 'play'}
          />
        );
      })}
    </Container>
  );
}

const Container = styled(Flex)`
  flex-direction: column;
  gap: 18px;

  font-size: 12px;

  .width_1100 & {
    font-size: 14px;
  }
`;

const CashoutItemWrapper = styled(Box)`
  background: linear-gradient(180deg, #4f617b -0.27%, rgba(26, 41, 61, 0) 100%);
  border-radius: 0px 7px 7px 0px;
  padding-top: 1px;

  & > div {
    backdrop-filter: blur(4px);
  }
`;
