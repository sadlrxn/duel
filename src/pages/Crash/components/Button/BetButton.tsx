import React from 'react';

import { BetButton as BaseBetButton } from 'components/GameTabs/styles';
import { Flex, Chip, ButtonProps } from 'components';
import { useAppSelector } from 'state';
import { useCrash } from 'hooks';
import { convertBalanceToChip } from 'utils/balance';

interface BetButtonProps extends ButtonProps {}

const BetButton = React.forwardRef<HTMLButtonElement, BetButtonProps>(
  ({ ...props }, ref) => {
    const status = useAppSelector(state => state.crash.status);
    const betCountLimit = useAppSelector(
      state => state.meta.crash.betCountLimit
    );
    const { userBets, totalBet, couponTotalBet, userBetCount } = useCrash();

    return (
      <BaseBetButton
        ref={ref}
        fontWeight={700}
        fontSize="14px"
        borderRadius="12px"
        height="100%"
        disabled={
          (status === 'bet' && userBetCount >= betCountLimit) ||
          status === 'ready' ||
          status === 'explosion' ||
          status === 'back' ||
          (status === 'play' && userBets.length === 0)
        }
        {...props}
      >
        {status === 'play' ? (
          <Flex gap={8} alignItems="center">
            {userBets.length === 0 ? (
              <>Place Bet</>
            ) : (
              <>
                Cash Out
                {userBets.length !== 1 && ' All'}
                <Flex flexWrap="wrap" gap={2}>
                  {!!couponTotalBet && (
                    <Chip
                      decimal={2}
                      chipType="coupon"
                      price={convertBalanceToChip(couponTotalBet)}
                      fontSize="12px"
                      fontWeight={700}
                      color="black"
                    />
                  )}
                  <Chip
                    decimal={2}
                    price={convertBalanceToChip(totalBet)}
                    fontSize="12px"
                    fontWeight={700}
                    color="black"
                  />
                </Flex>
              </>
            )}
          </Flex>
        ) : (
          'Place Bet'
        )}
      </BaseBetButton>
    );
  }
);

BetButton.displayName = 'CrashBaseBetButton';

export default BetButton;
