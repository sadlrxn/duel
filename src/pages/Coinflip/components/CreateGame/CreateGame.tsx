import { useCallback, useState, ChangeEvent, useMemo, useEffect } from 'react';
import styled, { useTheme } from 'styled-components';
import ClipLoader from 'react-spinners/ClipLoader';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';

import { Box, Flex, Text, useModal, BetAllModal } from 'components';

import { useAppDispatch, useAppSelector } from 'state';
import { setCreateRequest } from 'state/coinflip/actions';
import { updateBalance, setBalanceType } from 'state/user/actions';
import state from 'state';
import { sendMessage } from 'state/socket';
import { toast } from 'utils/toast';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';
import { useSound } from 'hooks';

import {
  Container,
  SideContainer,
  Sides,
  Side,
  Title,
  // DataContainer,
  // PrivateButton,
  BetAllButton,
  CreateButton,
  FlexBox
} from './CreateGame.styles';
import RoundInput from '../Input/RoundInput';
import { BetInput } from '../Input';
import { Coin } from '../Coin';

export default function CreateGame() {
  const dispatch = useAppDispatch();
  const theme = useTheme();
  const { buttonPlay } = useSound();

  const user = useAppSelector(state => state.user);
  const game = useAppSelector(state => state.coinflip);
  const meta = useAppSelector(state => state.meta.coinflip);
  const [value, setValue] = useState('');
  const [count, setCount] = useState('1');
  const [side, setSide] = useState<'duel' | 'ana' | 'rnd'>('duel');
  const [opponent, setOpponent] = useState<'dueler' | 'bot'>('dueler');
  const [userMaxBet, setUserMaxBet] = useState(0);

  const request = useMemo(() => game.createRequest > 0, [game]);

  useEffect(() => {
    if (
      opponent === 'bot' &&
      user.balances.coupon.balance >= convertChipToBalance(0.01)
    ) {
      dispatch(setBalanceType('coupon'));
    } else {
      dispatch(setBalanceType('chip'));
    }
  }, [dispatch, opponent, user.balances.coupon.balance, user.betBalanceType]);

  useEffect(() => {
    setUserMaxBet(user.balance);
  }, [user.balance]);

  const onChangeBet = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    if (+e.target.value * +count > convertBalanceToChip(userMaxBet)) return;
    setValue(e.target.value);
  };

  const duplicateBet = useCallback(() => {
    let v = +value * 2;
    if (v <= 0) v = 1;
    if (v <= convertBalanceToChip(meta.minBetAmount))
      v = convertBalanceToChip(meta.minBetAmount);
    if (v > convertBalanceToChip(meta.maxBetAmount))
      v = convertBalanceToChip(meta.maxBetAmount);
    if (v > convertBalanceToChip(userMaxBet))
      v = convertBalanceToChip(userMaxBet);
    if (v * +count > convertBalanceToChip(userMaxBet)) return;
    setValue(v.toFixed(2).toString());
  }, [meta, userMaxBet, value, count]);

  const divideBetInHalf = useCallback(() => {
    let v = +value / 2;
    if (v > convertBalanceToChip(meta.maxBetAmount))
      v = convertBalanceToChip(meta.maxBetAmount);
    if (v <= convertBalanceToChip(meta.minBetAmount))
      v = convertBalanceToChip(meta.minBetAmount);
    if (v > convertBalanceToChip(userMaxBet))
      v = convertBalanceToChip(userMaxBet);
    if (v * +count > convertBalanceToChip(userMaxBet)) return;
    setValue(v.toFixed(2).toString());
  }, [meta, value, count, userMaxBet]);

  const onChangeCount = (e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    if (+e.target.value !== Math.floor(+e.target.value))
      e.target.value = Math.floor(+e.target.value).toString();
    if (+e.target.value > meta.createRoundLimit) return;
    // console.log(+e.target.value * +value);
    if (+e.target.value * +value > convertBalanceToChip(userMaxBet)) return;
    setCount((+e.target.value).toFixed());
  };

  const increaseCount = useCallback(() => {
    if (+count === meta.createRoundLimit) return;
    if ((+count + 1) * +value > convertBalanceToChip(userMaxBet)) return;
    setCount((+count + 1).toString());
  }, [count, value, userMaxBet, meta]);

  const decreaseCount = useCallback(() => {
    if (+count === 1) return;
    setCount((+count - 1).toString());
  }, [count]);

  const handleCreate = useCallback(
    (e: any) => {
      e.preventDefault();
      if (request) return;
      if (user.name === '') {
        toast.info('Please sign in.');
        return;
      }
      if (user.balance < meta.minBetAmount) {
        toast.error('Insufficient funds');
        return;
      }
      if (+value < convertBalanceToChip(meta.minBetAmount)) {
        toast.warning(
          `Minimum bet is ${convertBalanceToChip(meta.minBetAmount)} CHIPs.`
        );
        setValue('');
        return;
      }
      if (+value > convertBalanceToChip(meta.maxBetAmount)) {
        toast.warning(
          `Max bet is ${convertBalanceToChip(meta.maxBetAmount)} CHIPs.`
        );
        return;
      }

      const paidBalanceType = user.betBalanceType;

      if (+value * +count > convertBalanceToChip(user.balance)) {
        toast.error('Insufficient funds');
        return;
      }

      setValue((+value).toFixed(2).toString());
      const amount = Math.floor(+convertChipToBalance(+value));

      state.dispatch(
        updateBalance({
          type: -1,
          usdAmount: amount * +count,
          wagered: amount * +count,
          balanceType: paidBalanceType
        })
      );

      var contents = [];
      for (var i = 0; i < +count; i++) {
        var coinSide =
          side === 'duel' ? 'heads' : side === 'ana' ? 'tails' : '';
        if (coinSide === '') {
          coinSide = Math.random() > 0.5 ? 'heads' : 'tails';
        }
        contents.push({
          eventType: 'bet',
          amount,
          side: coinSide,
          opponent,
          paidBalanceType
        });
      }
      state.dispatch(
        sendMessage({
          type: 'event',
          room: 'coinflip',
          content: JSON.stringify(contents)
        })
      );
      state.dispatch(setCreateRequest(+count));
    },
    [request, user, meta, value, count, side, opponent]
  );

  const [onBetAll] = useModal(
    <BetAllModal
      setValue={setValue}
      count={+count}
      maxAmount={meta.maxBetAmount}
    />,
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

  const handleToggle = useCallback(() => {
    if (opponent === 'dueler') {
      setOpponent('bot');
      state.dispatch(setBalanceType('coupon'));
    } else {
      setOpponent('dueler');
      state.dispatch(setBalanceType('chip'));
    }
    buttonPlay && buttonPlay();
  }, [opponent, buttonPlay]);

  return (
    <form onSubmit={handleCreate}>
      <Container>
        <SideContainer>
          <Text
            color="#D0DAEB"
            fontSize="18px"
            fontWeight={600}
            letterSpacing="0.08em"
            textTransform="uppercase"
            mb="20px"
          >
            Choose Your Side
          </Text>
          <Sides>
            <Side active={side === 'duel'} onClick={() => setSide('duel')}>
              <Coin
                boxShadow={`0 0 30px ${theme.coinflip.duel}`}
                size={72}
                scale={0.64}
              />

              <Text
                letterSpacing="0.08em"
                fontSize={'10px'}
                fontWeight={600}
                color="#D0DAEB"
                mt="12px"
              >
                GREEN
              </Text>
            </Side>
            <Side active={side === 'ana'} onClick={() => setSide('ana')}>
              <Coin
                boxShadow={`0 0 30px ${theme.coinflip.ana}`}
                size={72}
                scale={0.64}
                side="ana"
              />
              <Text
                letterSpacing="0.08em"
                fontSize={'10px'}
                fontWeight={600}
                color="#D0DAEB"
                mt="12px"
              >
                PURPLE
              </Text>
            </Side>
            <Side active={side === 'rnd'} onClick={() => setSide('rnd')}>
              <Coin
                boxShadow={`0 0 30px ${theme.coinflip.rnd}`}
                size={72}
                scale={0.45}
                side="rnd"
              />
              <Text
                letterSpacing="0.08em"
                fontSize={'10px'}
                fontWeight={600}
                color="#D0DAEB"
                mt="12px"
              >
                RANDOM
              </Text>
            </Side>
          </Sides>
        </SideContainer>

        <Box>
          <Text
            textTransform="uppercase"
            color="#D0DAEB"
            fontSize="18px"
            fontWeight={600}
            letterSpacing="0.08em"
            mb="20px"
          >
            Place Your Bet
          </Text>

          <FlexBox>
            <Flex flexDirection="column" gap={10}>
              <Flex justifyContent="space-between">
                <Title>Number of Games</Title>
              </Flex>
              <RoundInput
                isCount
                tabIndex={1}
                value={count}
                onChange={onChangeCount}
                duplicateBet={increaseCount}
                divideBetInHalf={decreaseCount}
              />
            </Flex>

            <Text color={'#898F97'} fontSize="20px" mx="20px" mt="45px">
              X
            </Text>
            <Flex flexDirection="column" gap={10}>
              <Flex justifyContent="space-between">
                <Title>Bet Amount</Title>
                <BetAllButton type="button" onClick={handleBetMax}>
                  All-in
                </BetAllButton>
              </Flex>
              <BetInput
                tabIndex={2}
                value={value}
                onChange={onChangeBet}
                duplicateBet={duplicateBet}
                divideBetInHalf={divideBetInHalf}
              />
            </Flex>
          </FlexBox>
        </Box>

        <Box>
          <Text
            textTransform="uppercase"
            color={'#D0DAEB'}
            fontSize="18px"
            fontWeight={600}
            mb="20px"
            letterSpacing="0.08em"
          >
            Choose Opponent
          </Text>
          <FlexBox gap={18} mb="20px">
            <Text color={'#4F617B'} fontSize="16px">
              Dueler
            </Text>
            <StyledToggle
              tabIndex={3}
              defaultChecked={false}
              icons={false}
              onChange={handleToggle}
            />
            <Text color={'#4F617B'} fontSize="16px">
              DuelBot
            </Text>
          </FlexBox>
          <FlexBox>
            {/* <CreateButton tabIndex={3} onClick={request ? null : handleCreate}> */}
            <CreateButton tabIndex={4} type="submit">
              {request ? (
                <ClipLoader color="#ffffff" loading={request} size={20} />
              ) : (
                'Create Game'
              )}
            </CreateButton>
          </FlexBox>
        </Box>
      </Container>
    </form>
  );
}

const StyledToggle = styled(Toggle)`
  .react-toggle-track {
    width: 56px;
    height: 30px;
    background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0) 162.5%);
    background-color: transparent !important;
  }

  .react-toggle-thumb {
    width: 28px;
    height: 28px;
    background: #4fff8b;
  }

  .react-toggle--checked:hover:not(.react-toggle--disabled)
    .react-toggle-track {
    background-color: transparent;
  }
`;
