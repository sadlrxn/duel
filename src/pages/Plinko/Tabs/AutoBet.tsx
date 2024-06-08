import React, { useState, useCallback, ChangeEvent } from 'react';
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

export default function AutoBet() {
  const user = useAppSelector(state => state.user);
  const betBalanceType = useAppSelector(state => state.user.betBalanceType);
  const { difficulties } = useAppSelector(state => state.meta.plinko);
  const { gameMode, gameRows, setGameMode, setGameRows } = usePlinko();

  const [id, setId] = useState(1);
  const [amount, setAmount] = useState('');
  const [dropSpeed, setDropSpeed] = useState([0.1]);

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
    state.dispatch(
      plinkoActions.addBall({
        roundId: id,
        path: 'LRRLRRLRLRR',
        betAmount: 500,
        lines: 16,
        level: 'low',
        time: Date.now(),
        multiplier: 0.5
      })
    );
    setId(id + 1);
  }, [id]);

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
                disabled={index === 3}
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
        <Box px="8px" width="100%">
          <Range
            values={dropSpeed}
            step={0.1}
            min={0.1}
            max={1}
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
                        min: 0.1,
                        max: 1
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
