import { useCallback, useMemo } from 'react';
import styled from 'styled-components';
import { shallowEqual } from 'react-redux';

import { Box, Flex, Chip, Span, Button } from 'components';
import { useAppSelector } from 'state';
import { CrashAutoBet, CrashGameStatus } from 'api/types/crash';
import { useCrash } from 'hooks';
import { convertBalanceToChip } from 'utils/balance';

const AddSlot = ({ handleAddSlot, disabled }: any) => {
  return (
    <CashoutItemWrapper padding="1px">
      <Flex
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        width="100%"
        borderRadius="0px 7px 7px 0px"
        background="linear-gradient(180deg, rgba(19, 32, 49, 0.75) 0%, rgba(26, 41, 60, 0.75) 67.97%);"
        padding="8px 10px"
      >
        <Button
          padding="6px 8px"
          borderRadius="5px"
          fontSize="1em"
          fontWeight={700}
          onClick={handleAddSlot}
          color="black"
          disabled={disabled}
          minWidth="max-content"
        >
          Add Slot
        </Button>
      </Flex>
    </CashoutItemWrapper>
  );
};

const CashoutItem = ({
  bet,
  multiplier,
  maxCashOut,
  status,
  itemIndex = 0,
  totalCount
}: {
  bet: CrashAutoBet;
  multiplier: number;
  status: CrashGameStatus;
  maxCashOut: number;
  itemIndex: number;
  totalCount: number;
}) => {
  const {
    setCurrentAutoBet,
    setCurrentAutoBetIndex,
    currentAutoBetIndex,
    handleReset,
    handleRemove,
    autoBetEnable
  } = useCrash();

  const handleClick = useCallback(() => {
    setCurrentAutoBetIndex(itemIndex);
    setCurrentAutoBet(bet);
  }, [bet, itemIndex, setCurrentAutoBet, setCurrentAutoBetIndex]);

  const [_, color] = useMemo(() => {
    let winType = -1;
    let color = 'warning';
    if (!bet.isBetted) {
      winType = 0;
      color = 'white';
    } else if (status === 'bet' || status === 'ready') {
      winType = 0;
      color = 'white';
    } else if (status === 'play') {
      if (bet.profit > 0) {
        winType = 1;
        color = 'success';
      } else {
        winType = 0;
        color = 'white';
      }
    } else if (bet.profit > 0) {
      winType = 1;
      color = 'success';
    }
    return [winType, color];
  }, [bet.profit, status, bet.isBetted]);

  const price = useMemo(() => {
    if (!bet.isBetted) return convertBalanceToChip(bet.betAmount);
    if (bet.profit > 0) return convertBalanceToChip(bet.profit);
    return status === 'play'
      ? convertBalanceToChip(Math.min(bet.betAmount * multiplier, maxCashOut))
      : status === 'bet' || status === 'ready'
      ? convertBalanceToChip(bet.betAmount)
      : bet.profit
      ? convertBalanceToChip(bet.profit)
      : convertBalanceToChip(-bet.betAmount);
  }, [bet.betAmount, bet.isBetted, bet.profit, maxCashOut, multiplier, status]);

  return (
    <CashoutItemWrapper
      onClick={handleClick}
      cursor="pointer"
      padding={itemIndex === currentAutoBetIndex ? '0px' : '1px'}
    >
      <Flex
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        width="100%"
        borderRadius="0px 7px 7px 0px"
        background="linear-gradient(180deg, rgba(19, 32, 49, 0.75) 0%, rgba(26, 41, 60, 0.75) 67.97%);"
        padding="8px 10px"
        mb="8px"
        boxShadow={
          itemIndex === currentAutoBetIndex
            ? 'inset 0px 0px 8px rgba(79, 255, 139, 0.5)'
            : 'none'
        }
        borderColor="success"
        borderStyle="solid"
        borderWidth={itemIndex === currentAutoBetIndex ? '1px' : '0px'}
      >
        {bet.isComplete ? (
          <Span color="success" fontSize="1em" fontWeight={700}>
            Complete
          </Span>
        ) : (
          <Chip
            decimal={2}
            chipType={bet.paidBalanceType}
            price={price}
            fontSize="1em"
            fontWeight={700}
            color={color}
          />
        )}

        <Span
          color="#fffc"
          fontSize="9px"
          fontWeight={500}
          lineHeight={1}
          my="2px"
        >
          {'Auto ' + bet.cashOutAt?.toFixed(2) + 'x'}
        </Span>
        <Span
          color="#fffc"
          fontSize="9px"
          fontWeight={500}
          lineHeight={1}
          my="2px"
        >
          {'Rounds ' +
            (bet.rounds > 0 ? Math.max(bet.rounds - bet.bettedRounds, 0) : '∞')}
        </Span>
        <Span
          color={
            bet.isComplete || !autoBetEnable
              ? bet.pnl > 0
                ? 'success'
                : bet.pnl < 0
                ? 'warning'
                : '#fffc'
              : '#fffc'
          }
          fontSize="9px"
          fontWeight={500}
          lineHeight={1}
          my="2px"
        >
          {'PnL ' +
            (bet.pnl < 0 ? '' : '+') +
            convertBalanceToChip(bet.pnl).toFixed(2)}
        </Span>
        <Button
          padding="6px 8px"
          marginBottom="-16px"
          borderRadius="5px"
          fontSize="1em"
          fontWeight={700}
          color="black"
          disabled={autoBetEnable}
          onClick={() => {
            if (totalCount === 1 && itemIndex === 0) handleReset();
            else handleRemove(itemIndex);
          }}
        >
          {totalCount === 1 && itemIndex === 0 ? 'Reset' : 'Remove'}
        </Button>
      </Flex>
    </CashoutItemWrapper>
  );
};

export default function AutoCashout({ ...props }: any) {
  const status = useAppSelector(state => state.crash.status, shallowEqual);
  const { maxCashOut } = useAppSelector(state => state.meta.crash);
  const { autoBets, multiplier, handleAddSlot, autoBetEnable } = useCrash();

  return (
    <Container {...props}>
      {(autoBets as CrashAutoBet[]).map((bet, index) => {
        return (
          <CashoutItem
            key={`crash_automatic_${index}`}
            bet={bet}
            itemIndex={index}
            totalCount={autoBets.length}
            multiplier={multiplier}
            status={status}
            maxCashOut={maxCashOut}
          />
        );
      })}
      {autoBets.length < 5 && (
        <AddSlot handleAddSlot={handleAddSlot} disabled={autoBetEnable} />
      )}
    </Container>
  );
}

const Container = styled(Flex)`
  flex-direction: column;
  gap: 10px;

  font-size: 12px;

  .width_1100 & {
    font-size: 14px;
  }
`;

const CashoutItemWrapper = styled(Box)`
  background: linear-gradient(180deg, #4f617b -0.27%, rgba(26, 41, 61, 0) 100%);
  border-radius: 0px 7px 7px 0px;
  min-width: max-content;
  min-height: max-content;

  & > div {
    /* backdrop-filter: blur(4px); */
  }
`;
