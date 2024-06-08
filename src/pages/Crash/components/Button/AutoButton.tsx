import React, { useCallback } from 'react';
import { shallowEqual } from 'react-redux';

import { BetButton as BaseBetButton } from 'components/GameTabs/styles';
import { ButtonProps } from 'components';
import { useCrash } from 'hooks';
import { CrashAutoBet } from 'api/types/crash';
import { toast } from 'utils/toast';
import { useAppSelector } from 'state';

interface AutoButtonProps extends ButtonProps {}

const AutoButton = React.forwardRef<HTMLButtonElement, AutoButtonProps>(
  ({ ...props }, ref) => {
    const status = useAppSelector(state => state.crash.status, shallowEqual);
    const { minBetAmount, maxBetAmount, minCashOutAt } = useAppSelector(
      state => state.meta.crash
    );
    const { autoBetEnable, setAutoBetEnable, autoBets } = useCrash();

    const handleClick = useCallback(() => {
      if (autoBetEnable) setAutoBetEnable(false);
      else {
        if (
          autoBets.every((bet: CrashAutoBet) => {
            return (
              bet.betAmount >= minBetAmount &&
              bet.betAmount <= maxBetAmount &&
              bet.cashOutAt >= minCashOutAt
            );
          })
        )
          setAutoBetEnable(true);
        else {
          setAutoBetEnable(false);
          toast.warn('Invalid Input');
        }
      }
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [autoBetEnable, setAutoBetEnable]);

    return (
      <BaseBetButton
        {...props}
        ref={ref}
        fontWeight={700}
        fontSize="14px"
        borderRadius="12px"
        height="100%"
        disabled={(status === 'play' || status === 'ready') && autoBetEnable}
        onClick={handleClick}
      >
        {autoBetEnable ? 'Stop' : 'Start'} Auto Bet
      </BaseBetButton>
    );
  }
);

AutoButton.displayName = 'CrashAutoBetButton';

export default AutoButton;
