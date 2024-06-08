import React, { useState, useCallback, ChangeEvent, useMemo } from 'react';
import styled from 'styled-components';
import { Range, getTrackBackground } from 'react-range';

import coin from 'assets/imgs/coins/coin.svg';
import coinBlue from 'assets/imgs/coins/coin-blue.svg';

import {
  useModal,
  BetAllModal,
  InputItem,
  Label,
  Flex,
  Grid,
  Box
} from 'components';
import state, { useAppSelector } from 'state';
import { toast } from 'utils/toast';
import * as plinkoActions from 'state/plinko/actions';
import { convertBalanceToChip } from 'utils/balance';
import { usePlinko } from 'hooks';

import { BetButton, ModeButton } from '../components';
import { plinkoRows } from '../config';

export default function Manual() {
  const user = useAppSelector(state => state.user);
  const betBalanceType = useAppSelector(state => state.user.betBalanceType);
  const { difficulties } = useAppSelector(state => state.meta.plinko);
  const { balls } = useAppSelector(state => state.plinko);
  const {
    gameMode,
    gameRows,
    setGameMode,
    setGameRows,
    dropSpeed,
    setDropSpeed
  } = usePlinko();

  const [id, setId] = useState(1);
  const [amount, setAmount] = useState('');

  const isPlaying = useMemo(() => balls.length > 0, [balls.length]);

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
  }, [amount, user.balance]);

  const divideBetInHalf = useCallback(() => {
    let v = +amount / 2;
    if (v > convertBalanceToChip(user.balance))
      v = convertBalanceToChip(user.balance);
    setAmount(v.toFixed(2).toString());
  }, [amount, user.balance]);

  const handleBet = useCallback(() => {
    let path: string = '';

    for (let i = 0; i < gameRows; i++) {
      if (Math.random() < 0.5) path = path + 'L';
      else path = path + 'R';
    }

    state.dispatch(
      plinkoActions.addBall({
        roundId: id,
        path,
        betAmount: 500,
        lines: gameRows,
        level: gameMode,
        time: Date.now(),
        multiplier: 0.5
      })
    );
    setId(id + 1);
  }, [gameMode, gameRows, id]);

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

  const handleMax = useCallback(() => {
    const max = user.balance;
    setAmount(convertBalanceToChip(max).toFixed(2));
  }, [user.balance]);

  return (
    <Container>
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
        type="number"
      />

      <Flex flexDirection="column" gap={8}>
        <Label
          ml="2px"
          color="#4F617B"
          fontWeight={400}
          lineHeight={1}
          fontSize="1em"
        >
          Game Mode
        </Label>
        <Grid
          gridTemplateColumns={'repeat(auto-fit,minmax(130px,1fr))'}
          gap={10}
        >
          {difficulties.map((difficulty, index) => {
            return (
              <ModeButton
                mode={difficulty}
                key={difficulty}
                selected={difficulty === gameMode}
                disabled={index === 3 || isPlaying}
                onClick={() => setGameMode(difficulty)}
              />
            );
          })}
        </Grid>
      </Flex>

      <Flex flexDirection="column" gap={8} mb="8px">
        <Label
          ml="2px"
          color="#4F617B"
          fontWeight={400}
          lineHeight={1}
          fontSize="1em"
        >
          Drop Speed
        </Label>
        <Flex
          justifyContent="space-between"
          width="100%"
          fontWeight={700}
          fontSize="10px"
          px="3px"
          my="-2px"
          style={{ color: 'white' }}
        >
          <Label width="50px" textAlign="left">
            Slow
          </Label>
          <Label width="50px" textAlign="center">
            Medium
          </Label>
          <Label width="50px" textAlign="right">
            Turbo
          </Label>
        </Flex>
        <Box px="8px" width="100%">
          <Range
            values={dropSpeed}
            step={0.1}
            min={0.5}
            max={1.5}
            onChange={values => setDropSpeed(values)}
            renderThumb={({ props }) => {
              return (
                <div
                  {...props}
                  style={{
                    ...props.style,
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center'
                  }}
                >
                  <div
                    style={{
                      height: '16px',
                      width: '16px',
                      backgroundColor: '#4FFF8B',
                      borderRadius: '100%'
                    }}
                  />
                </div>
              );
            }}
            renderTrack={({ props, children }) => {
              return (
                <Flex
                  onMouseDown={props.onMouseDown}
                  onTouchStart={props.onTouchStart}
                  style={{ ...props.style, height: '10px', width: '100%' }}
                >
                  <div
                    ref={props.ref}
                    style={{
                      height: '6px',
                      width: '100%',
                      borderRadius: '30px',
                      alignSelf: 'center',
                      background: getTrackBackground({
                        values: dropSpeed,
                        colors: ['#1A5032', 'rgba(3, 6, 9, 0.6)'],
                        min: 0.5,
                        max: 1.5
                      })
                    }}
                  >
                    {children}
                  </div>
                </Flex>
              );
            }}
          />
        </Box>
      </Flex>

      <Flex flexDirection="column" gap={8}>
        <Label
          ml="2px"
          color="#4F617B"
          fontWeight={400}
          lineHeight={1}
          fontSize="1em"
        >
          Rows
        </Label>
        <Flex flexWrap="wrap" gap={10} justifyContent="center">
          {plinkoRows.map(rows => {
            return (
              <ModeButton
                width="38px"
                mode={rows}
                key={rows}
                selected={rows === gameRows}
                disabled={isPlaying}
                onClick={() => setGameRows(rows)}
              />
            );
          })}
        </Flex>
      </Flex>

      <BetButton onClick={handleBet} />
    </Container>
  );
}

const Container = styled.form`
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: 12px;
`;
