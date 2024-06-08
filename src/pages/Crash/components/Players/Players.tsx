import React, { useMemo } from 'react';
import styled from 'styled-components';

import { Flex, Chip, FlexProps } from 'components';
import { useAppSelector } from 'state';
import { useCrash } from 'hooks';

import { UserStatus } from '../UserStatus';
import { Persons } from '../Persons';
import { convertBalanceToChip } from 'utils/balance';

interface PlayersProps extends FlexProps {}

const Players = React.forwardRef<HTMLDivElement, PlayersProps>(
  ({ ...props }, ref) => {
    const { bets, status } = useAppSelector(state => state.crash);
    const { id: userId } = useAppSelector(state => state.user);
    const { maxCashOut } = useAppSelector(state => state.meta.crash);
    const { multiplier } = useCrash();

    const [totalBet, couponTotalBet] = useMemo(() => {
      let totalBet: number = 0;
      let couponTotalBet: number = 0;
      bets.forEach(bet => {
        if (bet.paidBalanceType === 'chip') totalBet += bet.betAmount;
        else if (bet.paidBalanceType === 'coupon')
          couponTotalBet += bet.betAmount;
      });
      return [totalBet, couponTotalBet];
    }, [bets]);

    const userAllBets = useMemo(() => {
      return bets.slice().sort((b1, b2) => {
        if (b1.user.id === userId) {
          if (b2.user.id === userId) {
            return (
              b2.betAmount * (b2.payoutMultiplier ? b2.payoutMultiplier : 1) -
              b1.betAmount * (b1.payoutMultiplier ? b1.payoutMultiplier : 1)
            );
          }
          return -1;
        }
        if (b2.user.id === userId) return 1;
        return b2.betAmount - b1.betAmount;
      });
    }, [bets, userId]);

    return (
      <Container {...props} ref={ref}>
        <Flex
          alignItems="center"
          px="7px"
          justifyContent="space-between"
          mb="13px"
        >
          <Persons />
          <Flex gap={6} alignItems="center">
            {couponTotalBet > 0 && (
              <Chip
                decimal={2}
                price={convertBalanceToChip(couponTotalBet)}
                chipType="coupon"
                size={12}
                color="#B2D1FF"
                fontSize="12px"
                fontWeight={700}
              />
            )}
            <Chip
              decimal={2}
              price={convertBalanceToChip(totalBet)}
              size={12}
              color="#B2D1FF"
              fontSize="12px"
              fontWeight={700}
            />
          </Flex>
        </Flex>
        <ListContainer>
          {userAllBets.map(bet => {
            return (
              <UserStatus
                key={bet.betId}
                user={bet.user}
                betAmount={bet.betAmount}
                multiplier={multiplier}
                payoutMultiplier={bet.payoutMultiplier}
                status={status}
                maxCashOut={maxCashOut}
                profit={bet.profit}
                balanceType={bet.paidBalanceType}
                isUser={bet.user.id === userId}
              />
            );
          })}
        </ListContainer>
      </Container>
    );
  }
);
Players.displayName = 'CrashPlayerList';

export default Players;

const ListContainer = styled(Flex)`
  flex-direction: column;
  gap: 5px;

  overflow: hidden auto;
  padding: 3px 3px 3px 1px;
`;

const Container = styled(Flex)`
  flex-direction: column;
  background: linear-gradient(180deg, #0f1a26 0%, #0f1a26 0.01%, #0f1a26 100%);
  border: 2px solid #0f1a26;
  border-radius: 13px;
  overflow: hidden auto;
  padding: 16px 10px 12px;
  gap: 5px;
  min-height: 210px;
`;
