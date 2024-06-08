import React, {
  useState,
  useCallback,
  ChangeEvent,
  useRef,
  useMemo
} from 'react';
import { shallowEqual } from 'react-redux';
import styled from 'styled-components';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';

import coin from 'assets/imgs/coins/coin.svg';
import coinBlue from 'assets/imgs/coins/coin-blue.svg';

import { Flex, Span, useModal, BetAllModal, InputItem } from 'components';
import { useAppSelector } from 'state';
import { toast } from 'utils/toast';
import { useCrash } from 'hooks';

import { Players, ManualButtons } from '../components';
import ManualCashout from './ManualCashout';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

interface ManualProps {
  gameHeight: number;
  gameWidth: number;
  bottomGap: number;
}

export default function Manual({
  gameHeight,
  gameWidth,
  bottomGap
}: ManualProps) {
  const user = useAppSelector(state => state.user);
  const status = useAppSelector(state => state.crash.status, shallowEqual);
  const { minCashOutAt, minBetAmount, maxBetAmount, maxCashOut } =
    useAppSelector(state => state.meta.crash);
  const { betBalanceType, crashAnimation: usingAnimation } = useAppSelector(
    state => state.user
  );
  const { handleToggleAnimation } = useCrash();

  const buttonRef = useRef<HTMLButtonElement>(null);
  const containerRef = useRef<HTMLFormElement>(null);
  const playersRef = useRef<HTMLDivElement>(null);

  const playerListHeight = useMemo(() => {
    let height = '210px';
    if (gameWidth > 700)
      height = `calc(100vh - ${
        gameHeight +
        (playersRef.current ? playersRef.current.offsetTop : 263) +
        50 +
        36
      }px)`;
    else height = '324px';
    return height;
  }, [gameHeight, gameWidth, playersRef]);

  const [amount, setAmount] = useState('');
  const [cashout, setCashout] = useState('');

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

  const handleMax = useCallback(() => {
    const max = Math.min(user.balance, maxBetAmount);
    setAmount(convertBalanceToChip(max).toFixed(2));
  }, [maxBetAmount, user.balance]);

  const [onBetAll] = useModal(
    <BetAllModal setValue={setAmount} maxAmount={maxBetAmount} />,
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

  const handleClick = useCallback(() => {
    // e.preventDefault();
    if (!buttonRef || !buttonRef.current) return;
    buttonRef.current.click();
  }, []);

  return (
    <>
      <Container
        onKeyDown={(e: React.KeyboardEvent<HTMLFormElement>) => {
          if (e.key === 'Enter') {
            handleClick();
          }
        }}
        ref={containerRef}
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
        />

        <ManualButtons
          buttonRef={buttonRef}
          gameHeight={gameHeight + bottomGap}
          gameWidth={gameWidth}
          amount={amount}
          cashout={cashout}
        />

        <Players ref={playersRef} maxHeight={playerListHeight} mt="8px" />

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
      {status !== 'explosion' && status !== 'back' && <Cashout />}
    </>
  );
}

const Cashout = styled(ManualCashout)`
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

const Container = styled.form`
  display: flex;
  flex-direction: column;
  min-width: 285px;
  width: 100%;
  gap: 12px;
  overflow: hidden auto;
  padding: 0px 2px;

  .width_700 & {
    height: 100%;
    max-height: none;
  }
`;
