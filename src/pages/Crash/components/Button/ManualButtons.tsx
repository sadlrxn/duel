import React, { useCallback } from 'react';
import styled from 'styled-components';

import { Flex, FlexProps, Button, Span, RepeatIcon } from 'components';
import { toast } from 'utils/toast';
import { useCrash } from 'hooks';

import { useAppSelector, useAppDispatch } from 'state';
import { sendMessage } from 'state/socket';
import { updateBalance } from 'state/user/actions';
import { setPending } from 'state/crash/actions';

import BetButton from './BetButton';
import { Players } from '../Players';
import { Persons } from '../Persons';
import { convertChipToBalance } from 'utils/balance';

interface ManualButtonsProps extends FlexProps {
  buttonRef: React.RefObject<HTMLButtonElement>;
  gameHeight: number;
  gameWidth: number;
  amount?: string;
  cashout?: string;
}

export default function ManualButtons({
  buttonRef,
  gameHeight,
  gameWidth,
  amount = '',
  cashout = '',
  ...props
}: ManualButtonsProps) {
  const dispatch = useAppDispatch();
  const user = useAppSelector(state => state.user);
  const { status, roundId } = useAppSelector(state => state.crash);
  const { minCashOutAt } = useAppSelector(state => state.meta.crash);
  const {
    userBets,
    repeatBet,
    repeatBetEnabled,
    setRepeatBet,
    showStatus,
    setShowStatus
  } = useCrash();

  const handleClick = useCallback(() => {
    if (user.name === '') {
      toast.info('Please sign in.');
      return;
    }
    if (cashout !== '' && +cashout !== 0 && +cashout < minCashOutAt) {
      toast.warning(`Minimal cashout is ${minCashOutAt}`);
      return;
    }
    if (status === 'bet') {
      const betAmount = +convertChipToBalance(+amount).toFixed(0);
      if (betAmount > user.balance) {
        toast.error(`Insufficient balance.`);
        return;
      }
      dispatch(
        sendMessage({
          type: 'event',
          room: 'crash',
          content: JSON.stringify({
            type: 'cash-in',
            content: JSON.stringify({
              amount: betAmount,
              balanceType: user.betBalanceType,
              roundId,
              cashOutAt:
                cashout === '' || +cashout === 0
                  ? undefined
                  : +(+cashout).toFixed(2)
            })
          })
        })
      );

      dispatch(
        updateBalance({
          type: -1,
          usdAmount: betAmount,
          wagered: betAmount,
          balanceType: user.betBalanceType
        })
      );
    }
    if (status === 'play') {
      userBets.forEach((bet: any) => {
        dispatch(
          sendMessage({
            type: 'event',
            room: 'crash',
            content: JSON.stringify({
              type: 'cash-out',
              content: JSON.stringify({
                roundId,
                betId: bet.betId
              })
            })
          })
        );

        dispatch(setPending({ betId: bet.betId, status: true }));
      });
    }
  }, [
    amount,
    cashout,
    dispatch,
    minCashOutAt,
    roundId,
    status,
    user.balance,
    user.betBalanceType,
    user.name,
    userBets
  ]);

  return (
    <Container {...props}>
      <CustomPlayers
        className={'crash_mobile_list' + (!showStatus ? '--hide' : '')}
        maxHeight={`calc(100vh - ${
          gameHeight - 20 - 20 + 16 + (status === 'bet' ? 104 : 0)
        }px)`}
        width={`${gameWidth - 16 - (status === 'bet' ? 18 : 90)}px`}
      />
      <ToggleButton
        variant="secondary"
        border="1px solid #2B4160"
        borderRadius="10px"
        background="#24354C"
        height="100%"
        width="50px"
        minWidth="50px"
        type="button"
        onClick={() => {
          setShowStatus(!showStatus);
        }}
      >
        {showStatus ? (
          <Span
            style={{
              transform: 'rotate(-90deg) scaleX(0.6)'
            }}
            fontWeight={600}
            color="#B2D1FF"
            fontSize="20px"
          >
            {'<'}
          </Span>
        ) : (
          <Persons />
        )}
      </ToggleButton>

      <BetButton
        ref={buttonRef}
        height="100%"
        type="button"
        onClick={handleClick}
      />

      <RepeatButton
        type="button"
        borderColor={repeatBet ? 'success' : '#2b4160'}
        disabled={!repeatBetEnabled}
        onClick={() => {
          setRepeatBet(!repeatBet);
        }}
      >
        <RepeatIcon color={repeatBet ? '#4FFF8B' : '#B2D1FF'} />
        <Span
          fontWeight={700}
          fontSize="8px"
          color={repeatBet ? 'success' : 'textWhite'}
        >
          Repeat Bet
        </Span>
      </RepeatButton>
    </Container>
  );
}

const ToggleButton = styled(Button)`
  .width_700 & {
    display: none;
  }
`;

const RepeatButton = styled(Button)`
  flex-direction: column;
  height: 100%;
  width: 50px;
  min-width: 50px;
  gap: 0px;
  padding: 0px;

  background: #24354c;
  border-width: 1px;
  border-style: solid;
  border-radius: 10px;
`;
RepeatButton.defaultProps = {
  variant: 'secondary'
};

const CustomPlayers = styled(Players)`
  position: absolute;
  left: 16px;
  top: 0px;
  transform: translate(0px, -100%);
  border-bottom-left-radius: 0px;
  border-bottom-right-radius: 0px;

  transition: all 0.5s;

  .width_700 & {
    display: none;
  }
`;

const Container = styled(Flex)`
  gap: 12px;

  position: absolute;
  transform: translate(0px, -100%);
  padding: 17px 16px 15px;
  background: #131e2d;
  height: 74px;
  width: 100%;

  top: -35px;
  left: 0px;

  margin-top: 0px;

  .width_700 & {
    position: relative;
    transform: translate(0, 0);
    background: transparent;
    top: 0px;
    padding: 0px;
    height: 40px;
    min-height: 40px;
    margin-top: 8px;
  }

  .crash_mobile_list--hide {
    display: none;
  }
`;
