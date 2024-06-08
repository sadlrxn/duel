import { useState, useCallback, ChangeEvent, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';

import coin from 'assets/imgs/coins/coin.svg';
import coinBlue from 'assets/imgs/coins/coin-blue.svg';
import percent from 'assets/imgs/icons/percent2.svg';

import { Box, Flex, Span, InputItem, useModal, BetAllModal } from 'components';
import { useAppSelector } from 'state';
import { toast } from 'utils/toast';
import { Players } from '../components';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';
import { useCrash } from 'hooks';
import { CrashAutoBet } from 'api/types/crash';

import AutomaticButtons from '../components/Button/AutomaticButtons';
import AutoCashout from './AutoCashout';

interface AutoBetProps {
  gameHeight: number;
  gameWidth: number;
  bottomGap: number;
}

export default function AutoBet({
  gameHeight,
  gameWidth,
  bottomGap
}: AutoBetProps) {
  const user = useAppSelector(state => state.user);
  const { betBalanceType, crashAnimation: usingAnimation } = useAppSelector(
    state => state.user
  );
  const { minCashOutAt, minBetAmount, maxBetAmount, maxCashOut } =
    useAppSelector(state => state.meta.crash);
  const {
    showStatus,
    handleToggleAnimation,
    currentAutoBet,
    autoBetEnable,
    setAutoBets,
    currentAutoBetIndex
  } = useCrash();

  const [amount, setAmount] = useState('');
  const [cashout, setCashout] = useState('');
  const [roundCount, setRoundCount] = useState('');
  const [stopProfit, setStopProfit] = useState('');
  const [stopLoss, setStopLoss] = useState('');
  const [onWin, setOnWin] = useState('');
  const [onLoss, setOnLoss] = useState('');

  useEffect(() => {
    const {
      betAmount,
      cashOutAt,
      rounds,
      onLoss,
      onWin,
      stopProfit,
      stopLoss
    } = currentAutoBet as CrashAutoBet;
    setAmount(convertBalanceToChip(betAmount).toFixed(2));
    setCashout(cashOutAt.toFixed(2));
    setRoundCount(rounds > 0 ? rounds.toFixed(0) : '');
    setOnLoss(onLoss ? onLoss.toString() : '');
    setOnWin(onWin ? onWin.toString() : '');
    setStopProfit(
      stopProfit ? convertBalanceToChip(stopProfit).toFixed(2) : ''
    );
    setStopLoss(stopLoss ? convertBalanceToChip(stopLoss).toFixed(2) : '');
  }, [currentAutoBet]);

  const duplicateBet = useCallback(() => {
    let v = +amount * 2;

    if (v <= convertBalanceToChip(minBetAmount))
      v = convertBalanceToChip(minBetAmount);
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    if (v > convertBalanceToChip(maxBetAmount))
      v = convertBalanceToChip(maxBetAmount);
    setAmount(v.toFixed(2).toString());
  }, [amount, minBetAmount, maxBetAmount, user.balance]);

  const divideBetInHalf = useCallback(() => {
    let v = +amount / 2;

    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    if (v > convertBalanceToChip(maxBetAmount))
      v = convertBalanceToChip(maxBetAmount);
    if (v < convertBalanceToChip(minBetAmount))
      v = convertBalanceToChip(minBetAmount);
    setAmount(v.toFixed(2).toString());
  }, [amount, maxBetAmount, minBetAmount, user.balance]);

  const onChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      e.preventDefault();
      if (+e.target.value < 0) return;
      if (convertChipToBalance(+e.target.value) > maxBetAmount) return;
      if (convertChipToBalance(+e.target.value) > user.balance) return;
      setAmount(e.target.value);
    },
    [maxBetAmount, user.balance]
  );

  const duplicateCashout = useCallback(() => {
    let v = cashout;
    if (v === '') v = minCashOutAt.toFixed(2);
    else if (+v < minCashOutAt) v = minCashOutAt.toFixed(2);
    else v = (+v + 0.1).toFixed(2);
    setCashout(v);
  }, [cashout, minCashOutAt]);

  const divideCashoutInHalf = useCallback(() => {
    let v = cashout;
    if (v === '') return;
    else if (+v < minCashOutAt + 0.1) v = '';
    else v = (+v - 0.1).toFixed(2);
    setCashout(v);
  }, [cashout, minCashOutAt]);

  const onCashoutChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    if (+e.target.value <= 0) {
      setCashout('');
      return;
    }
    const val = +e.target.value;
    if (val !== Math.floor(val * 100) / 100) return;
    setCashout(e.target.value);
  }, []);

  const increaseRoundCount = useCallback(() => {
    let count = +roundCount + 1;
    setRoundCount(count.toFixed(0));
  }, [roundCount]);

  const decreaseRoundCount = useCallback(() => {
    const count = +roundCount - 1;
    if (count <= 0) setRoundCount('');
    else setRoundCount(count.toFixed(0));
  }, [roundCount]);

  const onRoundCountChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const count = +e.target.value;
    if (Math.floor(count) !== count) return;
    if (+count < 0) return;
    setRoundCount(e.target.value);
  }, []);

  const increaseStopProfit = useCallback(() => {
    let val = +stopProfit + 1;
    setStopProfit(val.toFixed(2));
  }, [stopProfit]);

  const decreaseStopProfit = useCallback(() => {
    let val = +stopProfit - 1;
    if (val < 1) setStopProfit('');
    else setStopProfit(val.toFixed(2));
  }, [stopProfit]);

  const onStopProfitChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setStopProfit(e.target.value);
  }, []);

  const increaseStopLoss = useCallback(() => {
    let val = +stopLoss + 1;
    if (val < 1) val = 1;
    if (convertChipToBalance(val) > user.balance)
      val = convertBalanceToChip(user.balance);
    setStopLoss(val.toFixed(2));
  }, [stopLoss, user.balance]);

  const decreaseStopLoss = useCallback(() => {
    let val = +stopLoss - 1;
    if (val < 1) setStopLoss('');
    else setStopLoss(val.toFixed(2));
  }, [stopLoss]);

  const onStopLossChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setStopLoss(e.target.value);
  }, []);

  const increaseOnWin = useCallback(() => {
    let val = +onWin + 5;
    if (val === 0) setOnWin('');
    else setOnWin(val.toFixed(0));
  }, [onWin]);

  const decreaseOnWin = useCallback(() => {
    let val = +onWin - 5;
    if (val === 0) setOnWin('');
    else setOnWin(val.toFixed(0));
  }, [onWin]);

  const onOnWinChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    let val = +e.target.value;
    if (val % 1 !== 0) setOnWin(Math.floor(val).toFixed(0));
    else setOnWin(e.target.value);
  }, []);

  const increaseOnLoss = useCallback(() => {
    let val = +onLoss + 5;
    if (val === 0) setOnLoss('');
    else setOnLoss(val.toFixed(0));
  }, [onLoss]);

  const decreaseOnLoss = useCallback(() => {
    let val = +onLoss - 5;
    if (val === 0) setOnLoss('');
    else setOnLoss(val.toFixed(0));
  }, [onLoss]);

  const onOnLossChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    let val = +e.target.value;
    if (val % 1 !== 0) setOnLoss(Math.floor(val).toFixed(0));
    else setOnLoss(e.target.value);
  }, []);

  const handleMax = useCallback(() => {
    setAmount(convertBalanceToChip(user.balance).toFixed(2));
  }, [user.balance]);

  const [onBetAll] = useModal(
    <BetAllModal setValue={setAmount} />,
    true,
    true,
    true,
    'BetAllModal'
  );

  const handleBetMax = useCallback(() => {
    if (user.name === '') {
      toast.info('Please sign in.');
      return;
    }
    onBetAll();
  }, [user, onBetAll]);

  const newAutoBet = useMemo(() => {
    if (currentAutoBet.isBetted) return currentAutoBet;
    const newAutoBet: CrashAutoBet = {
      betAmount: convertChipToBalance(+amount),
      cashOutAt: +cashout,
      pnl: currentAutoBet.pnl,
      paidBalanceType: user.betBalanceType,
      bettedRounds: 0,
      betId: -1,
      roundId: -1,
      profit: 0,
      isComplete: false,
      isBetted: false,
      rounds: +roundCount <= 0 ? -1 : +roundCount,
      onLoss: +onLoss,
      onWin: +onWin,
      stopProfit: convertChipToBalance(+stopProfit),
      stopLoss: convertChipToBalance(+stopLoss)
    };
    return newAutoBet;
  }, [
    amount,
    currentAutoBet,
    cashout,
    onLoss,
    onWin,
    roundCount,
    stopLoss,
    stopProfit,
    user.betBalanceType
  ]);

  useEffect(() => {
    if (autoBetEnable) return;
    setAutoBets((prev: CrashAutoBet[]) => {
      if (prev.length <= currentAutoBetIndex || currentAutoBetIndex < 0)
        return prev;
      let newBets = [
        ...prev.slice(0, currentAutoBetIndex),
        newAutoBet,
        ...prev.slice(currentAutoBetIndex + 1)
      ];
      return newBets;
    });
  }, [autoBetEnable, currentAutoBetIndex, newAutoBet, setAutoBets]);

  return (
    <>
      <Container>
        <Box position="relative">
          <Flex
            className={'crash_autobet_tab' + (showStatus ? '--hide' : '')}
            flexDirection="column"
            gap={15}
          >
            <InputItem
              label="Bet Amount"
              iconUrl={betBalanceType === 'chip' ? coin : coinBlue}
              placeholder={'0.00'}
              allInButton
              upDownButton
              upButtonLabel="2x"
              downButtonLabel="1/2x"
              handleAllIn={handleBetMax}
              onInputChange={onChange}
              handleUp={duplicateBet}
              handleDown={divideBetInHalf}
              handleMax={handleMax}
              inputValue={amount}
              description={`Max Win: ${convertBalanceToChip(maxCashOut).toFixed(
                2
              )}`}
              type="number"
              tabIndex={0}
              readOnly={autoBetEnable}
            />

            <InputItem
              label="Cashout At"
              placeholder={'1.00x'}
              upDownButton
              upButtonLabel="+0.1"
              downButtonLabel="-0.1"
              onInputChange={onCashoutChange}
              handleUp={duplicateCashout}
              handleDown={divideCashoutInHalf}
              inputValue={cashout}
              inputSecondValue={(+amount * +cashout).toFixed(2)}
              inputSecondIcon={betBalanceType === 'chip' ? coin : coinBlue}
              type="number"
              step="any"
              tabIndex={1}
              readOnly={autoBetEnable}
            />

            <HorizontalDivider />

            <InputItem
              label="Number of Rounds"
              placeholder="∞"
              upDownButton
              upButtonLabel="+1"
              downButtonLabel="-1"
              onInputChange={onRoundCountChange}
              handleUp={increaseRoundCount}
              handleDown={decreaseRoundCount}
              inputValue={roundCount}
              type="number"
              step="1"
              readOnly={autoBetEnable}
            />

            <HorizontalDivider />

            <Flex width="100%" gap={10}>
              <InputItem
                label="On Win"
                status={
                  +onWin === 0
                    ? 'Disabled'
                    : +onWin > 0
                    ? 'Increase Bet By'
                    : 'Decrease Bet By'
                }
                iconUrl={percent}
                placeholder=""
                upDownButton
                upButtonLabel="+5%"
                downButtonLabel="-5%"
                onInputChange={onOnWinChange}
                handleUp={increaseOnWin}
                handleDown={decreaseOnWin}
                inputValue={onWin}
                type="number"
                step="1"
                width="100%"
                readOnly={autoBetEnable}
              />
              <InputItem
                label="Stop Profit"
                status={+stopProfit === 0 ? 'Disabled' : 'At'}
                iconUrl={betBalanceType === 'chip' ? coin : coinBlue}
                placeholder=""
                upDownButton
                upButtonLabel="+1"
                downButtonLabel="-1"
                onInputChange={onStopProfitChange}
                handleUp={increaseStopProfit}
                handleDown={decreaseStopProfit}
                inputValue={stopProfit}
                type="number"
                width="100%"
                readOnly={autoBetEnable}
              />
            </Flex>

            <HorizontalDivider />

            <Flex width="100%" gap={10}>
              <InputItem
                label="On Loss"
                status={
                  +onLoss === 0
                    ? 'Disabled'
                    : +onLoss > 0
                    ? 'Increase Bet By'
                    : 'Decrease Bet By'
                }
                iconUrl={percent}
                placeholder=""
                upDownButton
                upButtonLabel="+5%"
                downButtonLabel="-5%"
                onInputChange={onOnLossChange}
                handleUp={increaseOnLoss}
                handleDown={decreaseOnLoss}
                inputValue={onLoss}
                type="number"
                step="1"
                width="100%"
                readOnly={autoBetEnable}
              />
              <InputItem
                label="Stop Loss"
                status={+stopLoss === 0 ? 'Disabled' : 'At'}
                iconUrl={betBalanceType === 'chip' ? coin : coinBlue}
                placeholder=""
                upDownButton
                upButtonLabel="+1"
                downButtonLabel="-1"
                onInputChange={onStopLossChange}
                handleUp={increaseStopLoss}
                handleDown={decreaseStopLoss}
                inputValue={stopLoss}
                type="number"
                width="100%"
                readOnly={autoBetEnable}
              />
            </Flex>
          </Flex>

          <CustomPlayers
            className={'crash_autobet_list' + (showStatus ? '' : '--hide')}
          />
        </Box>

        <AutomaticButtons
          gameHeight={gameHeight + bottomGap}
          gameWidth={gameWidth}
        />

        <Flex
          alignItems="center"
          gap={8}
          fontSize="12px"
          fontWeight={600}
          height="18px"
        >
          <StyledToggle
            defaultChecked={usingAnimation}
            icons={false}
            onChange={handleToggleAnimation}
          />
          <Span color="#B2D1FF">
            {'Animations ' + (usingAnimation ? 'On' : 'Off')}
          </Span>
        </Flex>
      </Container>
      <Cashout />
    </>
  );
}

const StyledToggle = styled(Toggle)`
  .react-toggle-track {
    width: 30px;
    height: 18px;
    /* background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0) 162.5%); */
    /* background-color: transparent !important; */
    background-color: #242f42;
    opacity: 1;
    border: 2px solid #070b10c0;
  }

  &.react-toggle--checked {
    .react-toggle-thumb {
      background-color: #4fff8b;
      left: 14px;
    }

    .react-toggle-track {
      background-color: #1a5032;
    }
  }

  &.react-toggle:hover:not(.react-toggle--disabled) .react-toggle-track {
    background-color: #242f42;
  }

  &.react-toggle--checked:hover:not(.react-toggle--disabled)
    .react-toggle-track {
    background-color: #1a5032;
  }

  &.react-toggle:active:not(.react-toggle--disabled) .react-toggle-thumb {
    box-shadow: none;
  }

  &.react-toggle--focus .react-toggle-thumb {
    box-shadow: none;
  }

  .react-toggle-thumb {
    width: 14px;
    height: 14px;
    top: 2px;
    left: 2px;
    background-color: #768bad;
    /* background: #4fff8b; */
  }
`;

const Cashout = styled(AutoCashout)`
  position: absolute;
  left: 0;
  top: 0;
  transform: translate(0, calc(-100% - 94px - 36px));

  max-height: calc(100vh - 175px - 30px - 94px - 36px);
  height: 100%;
  overflow: auto;

  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }

  .width_700 & {
    left: calc(100% + 2px);
    top: 50%;
    transform: translate(0, -50%);
    max-height: 100%;
    justify-content: center;
  }

  z-index: -1;
`;

const CustomPlayers = styled(Players)`
  width: 100%;
  min-height: 76px;
  opacity: 1;
  pointer-events: all;

  margin-top: 30px;
  max-height: 324px;
  min-height: 210px;

  .width_700 & {
    position: absolute;
    left: 0;
    top: 0;
    height: 100%;
    max-height: none;

    margin-top: 0px;
    transition: opacity 0.5s;
  }
`;

const HorizontalDivider = styled(Box)`
  height: 1px;
  background: #303c4f;
  width: 100%;
`;

const Container = styled(Flex)`
  min-width: 285px;
  width: 100%;
  flex-direction: column;
  gap: 16px;
  padding: 0px 2px;

  .width_700 & {
    .crash_autobet_tab--hide {
      opacity: 0;
      transition: opacity 0.5s;
    }
    .crash_autobet_list--hide {
      opacity: 0;
      transition: opacity 0.5s;
      pointer-events: none;
    }
  }
`;
