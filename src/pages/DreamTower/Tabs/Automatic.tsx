import { useState, useCallback, ChangeEvent, useMemo, useEffect } from 'react';
import { NumericFormat } from 'react-number-format';
import styled from 'styled-components';
import { toast } from 'react-toastify';

import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { ReactComponent as CoinBlueIcon } from 'assets/imgs/coins/coin-blue.svg';
import { ReactComponent as TriangleArrowIcon } from 'assets/imgs/icons/triangle-arrow.svg';

import { Box, Flex, Text, Grid, useModal, BetAllModal } from 'components';
import { InputBox } from 'components/InputBox';
import { Count } from 'pages/Coinflip/components/List/List.styles';
import { BetAllButton } from 'pages/Coinflip/components/CreateGame/CreateGame.styles';

import state, { useAppSelector } from 'state';
import { DreamtowerDifficulty } from 'api/types/dreamtower';
import {
  setDifficulty,
  setBetAmount,
  setAutoBetCount,
  setChangeBetOnWin,
  setStopProfit,
  setChangeBetOnLoss,
  setStopLoss,
  setAutoBetAmount
} from 'state/dreamtower/actions';
import { useCustomSWR } from 'hooks';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

import { DownBtn, ModeBtn, StyledInput, UpperBtn } from './components';
// import { BetAllModal } from '../components/Modal';
import { LevelButtonContainer } from '../styles';

export default function Automatic() {
  const user = useAppSelector(state => state.user);
  const game = useAppSelector(state => state.dreamtower.game);
  const auto = useAppSelector(state => state.dreamtower.auto);
  const meta = useAppSelector(state => state.meta.dreamtower);
  const { betBalanceType } = useAppSelector(state => state.user);

  const { data: maxWinData } = useCustomSWR({
    key: 'dreamtower_max_win',
    route: '/dreamtower/max-win',
    method: 'get'
  });

  const maxWinning = useMemo(() => {
    if (!maxWinData) return 0;
    return maxWinData;
  }, [maxWinData]);

  const [amount, setAmount] = useState(
    auto.betAmount > 0 ? convertBalanceToChip(auto.betAmount).toFixed(2) : ''
  );
  const [level, setLevel] = useState<DreamtowerDifficulty>(game.difficulty);
  const [count, setCount] = useState(
    auto.betCount === undefined ? '' : auto.betCount.toFixed()
  );
  const [onWin, setOnWin] = useState(
    auto.changeBetOnWin === undefined ? '' : auto.changeBetOnWin.toFixed()
  );
  const [onLoss, setOnLoss] = useState(
    auto.changeBetOnLoss === undefined ? '' : auto.changeBetOnLoss.toFixed()
  );
  const [profit, setProfit] = useState(
    auto.stopProfit === undefined
      ? ''
      : convertBalanceToChip(auto.stopProfit).toFixed(2)
  );
  const [loss, setLoss] = useState(
    auto.stopLoss === undefined
      ? ''
      : convertBalanceToChip(auto.stopLoss).toFixed(2)
  );

  const onChangeBet = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setAmount(e.target.value);
    // state.dispatch(setBetAmount(+e.target.value * 100));
    // state.dispatch(setAutoBetAmount(+e.target.value * 100));
  };

  const increaseBet = useCallback(() => {
    let v = +amount * 2;
    if (v <= 0) v = 1;
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    setAmount(v.toFixed(2));
    // state.dispatch(setBetAmount(v * 100));
    // state.dispatch(setAutoBetAmount(v * 100));
  }, [amount, user.balance]);

  const decreaseBet = useCallback(() => {
    let v = +amount / 2;
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    if (v < 1) v = 1;
    setAmount(v.toFixed(2));
    // state.dispatch(setBetAmount(v * 100));
    // state.dispatch(setAutoBetAmount(v * 100));
  }, [amount, user.balance]);

  const handleSetDifficulty = useCallback((value: DreamtowerDifficulty) => {
    setLevel(value);
    state.dispatch(setDifficulty(value));
  }, []);

  const onChangeCount = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    else if (+e.target.value === 0) {
      setCount('');
      state.dispatch(setAutoBetCount(undefined));
    } else {
      setCount((+e.target.value).toFixed());
      state.dispatch(setAutoBetCount(+(+e.target.value).toFixed()));
    }
  };

  const increaseCount = useCallback(() => {
    setCount((+count + 1).toFixed());
    state.dispatch(setAutoBetCount(+(+count + 1).toFixed()));
  }, [count]);

  const decreaseCount = useCallback(() => {
    if (+count === 1) {
      setCount('');
      state.dispatch(setAutoBetCount(undefined));
    } else if (+count > 0) {
      setCount((+count - 1).toFixed());
      state.dispatch(setAutoBetCount(+(+count - 1).toFixed()));
    }
  }, [count]);

  const onChangeOnWin = ({ value }: any) => {
    if (value === '') {
      setOnWin('');
      state.dispatch(setChangeBetOnWin(undefined));
    } else if (+value >= -100) {
      setOnWin((+value).toFixed());
      state.dispatch(setChangeBetOnWin(+(+value).toFixed()));
    }
  };

  const increaseOnWin = useCallback(() => {
    setOnWin((+onWin + 5).toFixed());
    state.dispatch(setChangeBetOnWin(+(+onWin + 5).toFixed()));
  }, [onWin]);

  const decreaseOnWin = useCallback(() => {
    setOnWin((+onWin - 5).toFixed());
    state.dispatch(setChangeBetOnWin(+(+onWin - 5).toFixed()));
  }, [onWin]);

  const onChangeOnLoss = ({ value }: any) => {
    if (value === '') {
      setOnLoss('');
      state.dispatch(setChangeBetOnLoss(undefined));
    } else if (+value >= -100) {
      setOnLoss((+value).toFixed());
      state.dispatch(setChangeBetOnLoss(+(+value).toFixed()));
    }
  };

  const increaseOnLoss = useCallback(() => {
    setOnLoss((+onLoss + 5).toFixed());
    state.dispatch(setChangeBetOnLoss(+(+onLoss + 5).toFixed()));
  }, [onLoss]);

  const decreaseOnLoss = useCallback(() => {
    setOnLoss((+onLoss - 5).toFixed());
    state.dispatch(setChangeBetOnLoss(+(+onLoss - 5).toFixed()));
  }, [onLoss]);

  const onChangeProfit = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    else if (+e.target.value === 0) {
      setProfit('');
      state.dispatch(setStopProfit(undefined));
    } else {
      setProfit(e.target.value);
      state.dispatch(setStopProfit(convertChipToBalance(+e.target.value)));
    }
  };

  const increaseProfit = useCallback(() => {
    setProfit((+profit + 1).toFixed(2));
    state.dispatch(setStopProfit(convertChipToBalance(+profit + 1)));
  }, [profit]);

  const decreaseProfit = useCallback(() => {
    if (+profit === 1) {
      setProfit('');
      state.dispatch(setStopProfit(undefined));
    } else {
      setProfit((+profit - 1).toFixed(2));
      state.dispatch(setStopProfit(convertChipToBalance(+profit - 1)));
    }
  }, [profit]);

  const onChangeLoss = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    else if (+e.target.value === 0) {
      setLoss('');
      state.dispatch(setStopLoss(undefined));
    } else {
      setLoss(e.target.value);
      state.dispatch(setStopLoss(convertChipToBalance(+e.target.value)));
    }
  };

  const increaseLoss = useCallback(() => {
    setLoss((+loss + 1).toFixed(2));
    state.dispatch(setStopLoss(convertChipToBalance(+loss + 1)));
  }, [loss]);

  const decreaseLoss = useCallback(() => {
    if (+loss === 1) {
      setLoss('');
      state.dispatch(setStopLoss(undefined));
    } else {
      setLoss((+loss - 1).toFixed(2));
      state.dispatch(setStopLoss(convertChipToBalance(+loss - 1)));
    }
  }, [loss]);

  const disabled = useMemo(
    () => game.status === 'playing' || auto.status === 'running',
    [game.status, auto.status]
  );

  useEffect(() => {
    if (auto.status === 'running' && auto.betCount !== undefined) {
      setCount(auto.betCount.toFixed());
    }
  }, [auto.betCount, auto.status]);

  useEffect(() => {
    if (auto.status === 'running') {
      setAmount(convertBalanceToChip(auto.betAmount).toFixed(2));
    }
  }, [auto.betAmount, auto.status]);

  useEffect(() => {
    state.dispatch(setBetAmount(+convertChipToBalance(+amount).toFixed(0)));
    state.dispatch(setAutoBetAmount(+convertChipToBalance(+amount).toFixed(0)));
  }, [amount]);

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

  return (
    <div className="container">
      <Flex justifyContent="space-between">
        <Text fontSize={16} color="#4F617B" fontWeight={400}>
          Bet Amount
        </Text>
        <BetAllButton onClick={handleBetMax} disabled={disabled}>
          All-in
        </BetAllButton>
      </Flex>

      <InputBox
        as={Grid}
        gridTemplateColumns="max-content auto max-content"
        mt="15px"
        gap={10}
        p="3px 3px 3px 15px"
      >
        {betBalanceType === 'coupon' ? <CoinBlueIcon /> : <CoinIcon />}

        <StyledInput
          type="number"
          placeholder="0.00"
          value={amount}
          onChange={onChangeBet}
          disabled={disabled}
          onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
            if (
              e.key === 'e' ||
              e.key === 'E' ||
              e.key === '+' ||
              e.key === '-'
            )
              e.preventDefault();
          }}
          step="any"
        />
        <Box>
          <UpperBtn type="button" onClick={increaseBet} disabled={disabled}>
            <TriangleArrowIcon />
          </UpperBtn>
          <DownBtn type="button" onClick={decreaseBet} disabled={disabled}>
            <TriangleArrowIcon />
          </DownBtn>
        </Box>
      </InputBox>
      <Flex justifyContent="end">
        <Text color="#4F617B" fontSize={16} fontWeight={400} mb="0px">
          Max Win: {convertBalanceToChip(maxWinning).toFixed(2)}
        </Text>
      </Flex>

      <Text color="#4F617B" fontSize={16} fontWeight={400} mt="30px">
        Game Difficulty
      </Text>

      <LevelButtonContainer>
        {meta.difficulties.map((v, i) => (
          <ModeBtn
            selected={level.level === v.level}
            onClick={() => handleSetDifficulty(v)}
            key={i}
            disabled={disabled}
          >
            {v.level}
          </ModeBtn>
        ))}
      </LevelButtonContainer>

      <Text color="#4F617B" fontSize={16} fontWeight={400} mt="30px">
        Number of Bets
      </Text>

      <InputBox mt="15px" gap={10} p="3px 3px 3px 15px">
        <StyledInput
          type="number"
          placeholder="∞"
          value={count}
          onChange={onChangeCount}
          disabled={disabled}
        />
        <Box>
          <UpperBtn type="button" onClick={increaseCount} disabled={disabled}>
            <TriangleArrowIcon />
          </UpperBtn>
          <DownBtn
            type="button"
            onClick={decreaseCount}
            disabled={disabled || +Count === 0}
          >
            <TriangleArrowIcon />
          </DownBtn>
        </Box>
      </InputBox>

      <Grid gridTemplateColumns={'1fr 1fr'} gap={30} mt="30px">
        <div>
          <Flex flexDirection={['column', 'row']} gap={5}>
            <Text color="#4F617B" fontSize={16} fontWeight={700}>
              On Win
            </Text>
            <Text color="#95A3B9" fontSize={16} fontWeight={700}>
              {onWin === ''
                ? 'Disabled'
                : +onWin === 0
                ? 'Reset'
                : +onWin > 0
                ? 'Increase Bet By'
                : 'Decrease Bet By'}
            </Text>
          </Flex>

          <InputBox
            background="#03060999"
            mt="10px"
            gap={10}
            p="3px 3px 3px 15px"
          >
            <NumericFormat
              suffix="%"
              onValueChange={onChangeOnWin}
              value={onWin !== '' ? onWin + '%' : ''}
              disabled={disabled}
              isAllowed={({ value }: any) => +value >= -100 || value === '-'}
            />

            <StyledBox>
              <UpperBtn onClick={increaseOnWin} disabled={disabled}>
                <TriangleArrowIcon />
              </UpperBtn>
              <DownBtn
                onClick={decreaseOnWin}
                disabled={disabled || +onWin < -95}
              >
                <TriangleArrowIcon />
              </DownBtn>
            </StyledBox>
          </InputBox>
        </div>

        <div>
          <Flex flexDirection={['column', 'row']} gap={5}>
            <Text color="#4F617B" fontSize={16} fontWeight={700}>
              On Loss
            </Text>
            <Text color="#95A3B9" fontSize={16} fontWeight={700}>
              {onLoss === ''
                ? 'Disabled'
                : +onLoss === 0
                ? 'Reset'
                : +onLoss > 0
                ? 'Increase Bet By'
                : 'Decrease Bet By'}
            </Text>
          </Flex>

          <InputBox
            background="#03060999"
            mt="10px"
            gap={10}
            p="3px 3px 3px 15px"
          >
            <NumericFormat
              suffix="%"
              onValueChange={onChangeOnLoss}
              value={onLoss !== '' ? onLoss + '%' : ''}
              disabled={disabled}
              isAllowed={({ value }: any) => +value >= -100 || value === '-'}
            />

            <StyledBox>
              <UpperBtn onClick={increaseOnLoss} disabled={disabled}>
                <TriangleArrowIcon />
              </UpperBtn>
              <DownBtn
                onClick={decreaseOnLoss}
                disabled={disabled || +onLoss < -95}
              >
                <TriangleArrowIcon />
              </DownBtn>
            </StyledBox>
          </InputBox>
        </div>

        <div>
          <Flex flexDirection={['column', 'row']} gap={5}>
            <Text color="#4F617B" fontSize={16} fontWeight={700}>
              Stop Profit
            </Text>
            <Text color="#95A3B9" fontSize={16} fontWeight={700}>
              {profit === '' ? 'Disabled' : 'At'}
            </Text>
          </Flex>

          <InputBox
            background="#03060999"
            mt="10px"
            gap={10}
            p="3px 3px 3px 15px"
          >
            <StyledBox>
              {betBalanceType === 'coupon' ? <CoinBlueIcon /> : <CoinIcon />}
            </StyledBox>
            <StyledInput
              type="number"
              onChange={onChangeProfit}
              value={profit}
              disabled={disabled}
            />
            <StyledBox>
              <UpperBtn onClick={increaseProfit} disabled={disabled}>
                <TriangleArrowIcon />
              </UpperBtn>
              <DownBtn
                onClick={decreaseProfit}
                disabled={disabled || profit === ''}
              >
                <TriangleArrowIcon />
              </DownBtn>
            </StyledBox>
          </InputBox>
        </div>

        <div>
          <Flex flexDirection={['column', 'row']} gap={5}>
            <Text color="#4F617B" fontSize={16} fontWeight={700}>
              Stop Loss
            </Text>
            <Text color="#95A3B9" fontSize={16} fontWeight={700}>
              {loss === '' ? 'Disabled' : 'At'}
            </Text>
          </Flex>

          <InputBox
            background="#03060999"
            mt="10px"
            gap={10}
            p="3px 3px 3px 15px"
          >
            <StyledBox>
              {betBalanceType === 'coupon' ? <CoinBlueIcon /> : <CoinIcon />}
            </StyledBox>
            <StyledInput
              type="number"
              onChange={onChangeLoss}
              value={loss}
              disabled={disabled}
            />
            <StyledBox>
              <UpperBtn onClick={increaseLoss} disabled={disabled}>
                <TriangleArrowIcon />
              </UpperBtn>
              <DownBtn
                onClick={decreaseLoss}
                disabled={disabled || loss === ''}
              >
                <TriangleArrowIcon />
              </DownBtn>
            </StyledBox>
          </InputBox>
        </div>
      </Grid>
    </div>
  );
}

const StyledBox = styled(Box)`
  flex: none;
`;
