import { useState, useCallback, ChangeEvent, useMemo, useEffect } from 'react';

import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { ReactComponent as CoinBlueIcon } from 'assets/imgs/coins/coin-blue.svg';
import { ReactComponent as TriangleArrowIcon } from 'assets/imgs/icons/triangle-arrow.svg';

import { Box, Grid, Flex, Text, useModal, BetAllModal } from 'components';
import { InputBox } from 'components/InputBox';
import { BetAllButton } from 'pages/Coinflip/components/CreateGame/CreateGame.styles';

import { DreamtowerDifficulty } from 'api/types/dreamtower';
import state, { useAppSelector } from 'state';
import { setDifficulty, setBetAmount, clear } from 'state/dreamtower/actions';
import { useCustomSWR } from 'hooks';
import { toast } from 'utils/toast';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

import { DownBtn, ModeBtn, StyledInput, UpperBtn } from './components';
// import { BetAllModal } from '../components/Modal';
import { LevelButtonContainer } from '../styles';

export default function Manual() {
  const user = useAppSelector(state => state.user);
  const game = useAppSelector(state => state.dreamtower.game);
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
    game.betAmount > 0 ? convertBalanceToChip(game.betAmount).toFixed(2) : ''
  );
  const [level, setLevel] = useState<DreamtowerDifficulty>(game.difficulty);

  useEffect(() => {
    return () => {
      state.dispatch(clear());
    };
  }, []);

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setAmount(e.target.value);
    // state.dispatch(setBetAmount(+e.target.value * 100));
  };

  const duplicateBet = useCallback(() => {
    let v = +amount * 2;
    if (v <= 0) v = 1;
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    setAmount(v.toFixed(2).toString());
    // state.dispatch(setBetAmount(v * 100));
  }, [amount, user.balance]);

  const divideBetInHalf = useCallback(() => {
    let v = +amount / 2;
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    if (v < 1) v = 1;
    setAmount(v.toFixed(2).toString());
    // state.dispatch(setBetAmount(v * 100));
  }, [amount, user.balance]);

  useEffect(() => {
    state.dispatch(setBetAmount(+convertChipToBalance(+amount).toFixed(0)));
  }, [amount]);

  const handleSetDifficulty = useCallback((value: DreamtowerDifficulty) => {
    setLevel(value);
    state.dispatch(setDifficulty(value));
  }, []);

  const disabled = useMemo(() => game.status === 'playing', [game.status]);

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
          onChange={onChange}
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
          <UpperBtn type="button" onClick={duplicateBet} disabled={disabled}>
            <TriangleArrowIcon />
          </UpperBtn>
          <DownBtn type="button" onClick={divideBetInHalf} disabled={disabled}>
            <TriangleArrowIcon />
          </DownBtn>
        </Box>
      </InputBox>
      <Flex justifyContent="end">
        <Text color="#4F617B" fontSize={16} fontWeight={400} mb="0px">
          Max Win: {convertBalanceToChip(maxWinning).toFixed(2)}
        </Text>
      </Flex>

      <Text color="#4F617B" fontSize={16} fontWeight={400} mt="10px">
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
    </div>
  );
}
