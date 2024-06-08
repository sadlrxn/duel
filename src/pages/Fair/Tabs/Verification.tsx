import React, {
  useState,
  useRef,
  useEffect,
  useMemo,
  useCallback
} from 'react';
import styled from 'styled-components';
import { shallowEqual } from 'react-redux';
import Select, { Options, StylesConfig } from 'react-select';
import copy from 'copy-to-clipboard';

import copyIcon from 'assets/imgs/icons/copy.svg';

import { Box, Text, Flex, Grid, Button } from 'components';
import { JackpotFairData } from 'api/types/jackpot';
import { CoinflipRoundData } from 'api/types/coinflip';
import { useAppSelector } from 'state';

import {
  coinflipAlgorithm,
  jackpotAlgorithm,
  dreamtowerAlgorithm,
  crashAlgorithm
} from './constant';
import {
  getCoinflipInput,
  getCoinflipResult,
  getCrashInput,
  getCrashResult,
  getDreamTowerInput,
  getDreamTowerResult,
  getJackpotInput,
  getJackpotResult
} from './utils';
import { defaultTower, DreamtowerFairData } from 'api/types/dreamtower';
import Tower from 'pages/DreamTower/components/Modal/Fairness/Tower';

const VERIFICATION_OPTIONS: Options<{ label: string; value: string }> = [
  { label: 'Coin Flip', value: 'coinflip' },
  { label: 'Jackpot', value: 'jackpot' },
  { label: 'Dream Tower', value: 'dreamtower' },
  { label: 'Crash', value: 'crash' }
];

const customStyles: StylesConfig = {
  control: provided => ({
    ...provided,
    padding: '10px 20px',
    background: '#03060999',
    border: 0,
    borderRadius: 11,
    cursor: 'pointer',
    boxShadow: 'none',
    color: 'white',
    maxWidth: 400,
    '&:hover': {
      background: '#030609cc'
    }
  }),
  option: provided => ({
    ...provided,
    background: '#030609',
    color: 'white',
    padding: '15px 25px 12px',
    cursor: 'pointer',
    '&:hover': {
      background: '#080a0f'
    }
  }),
  input: base => ({
    ...base
    // display: "none",
  }),

  singleValue: provided => ({
    ...provided,
    color: 'white'
  }),
  indicatorSeparator: () => ({ display: 'none' }),
  dropdownIndicator: (provided, state) => ({
    ...provided,
    color: 'white',
    transition: '0.5s',
    transform: state.selectProps.menuIsOpen ? 'scaleY(-1)' : 'scaleY(1)',
    '&:hover': {
      color: 'white'
    }
  }),
  menu: provided => ({
    ...provided,
    background: 'transparent',
    boxShadow: 'none',
    maxWidth: 400,
    marginTop: '5px'
  }),
  menuList: provided => ({
    ...provided,
    background: '#030609',
    borderRadius: 11,
    transition: '0.5s',
    maxWidth: 400,
    padding: 0
  })
};

interface VerificationProps {
  gameType?: string;
  gameData?:
    | JackpotFairData
    | CoinflipRoundData
    | DreamtowerFairData
    | { serverSeed: string };
}

const Verification: React.FC<VerificationProps> = ({ gameType, gameData }) => {
  const crashClientSeed = useAppSelector(
    state => state.user.config.crashClientSeed,
    shallowEqual
  );

  const inputRef = useRef<any>(null);
  const algorithmRef = useRef<any>(null);
  const outputRef = useRef<any>(null);

  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [option, setOption] = useState(() => {
    let index = 0;
    if (gameType === 'jackpot') index = 1;
    else if (gameType === 'dreamtower') index = 2;
    else if (gameType === 'crash') index = 3;
    return VERIFICATION_OPTIONS[index];
  });
  const [towerData, setTowerData] = useState({
    tower: defaultTower,
    blocksInRow: 4
  });

  const gameAlgorithm = useMemo(() => {
    switch (option.value) {
      case 'coinflip':
        return coinflipAlgorithm;
      case 'jackpot':
        return jackpotAlgorithm;
      case 'dreamtower':
        return dreamtowerAlgorithm;
      case 'crash':
        return crashAlgorithm;
    }
    return '';
  }, [option]);

  const handleRun = useCallback(async () => {
    try {
      const data = JSON.parse(input.replace(/\s+|\r?\n|\r/g, ''));
      console.info(data);
      let result = '';
      if (option.value === 'coinflip') result = await getCoinflipResult(data);
      else if (option.value === 'jackpot')
        result = await getJackpotResult(data);
      else if (option.value === 'dreamtower') {
        let dresult = await getDreamTowerResult(data);
        result = dresult.result;
        for (let i = 0; i < dresult.towerData.tower.length / 2; i++) {
          const j = dresult.towerData.tower.length - i - 1;
          const temp = dresult.towerData.tower[i];
          dresult.towerData.tower[i] = dresult.towerData.tower[j];
          dresult.towerData.tower[j] = temp;
        }
        setTowerData(dresult.towerData);
      } else if (option.value === 'crash') result = await getCrashResult(data);

      setOutput(result);
    } catch {
      setOutput('Invalid input');
    }
  }, [input, option.value]);

  const [coinflipInput, jackpotInput, dreamtowerInput, crashInput] =
    useMemo(() => {
      return [
        getCoinflipInput(gameType, gameData),
        getJackpotInput(gameType, gameData),
        getDreamTowerInput(gameType, gameData),
        getCrashInput(
          gameType,
          gameData
            ? //@ts-ignore
              { serverSeed: gameData.serverSeed, clientSeed: crashClientSeed }
            : undefined
        )
      ];
    }, [gameType, gameData, crashClientSeed]);

  useEffect(() => {
    if (option.value === 'coinflip') setInput(coinflipInput);
    else if (option.value === 'jackpot') setInput(jackpotInput);
    else if (option.value === 'dreamtower') setInput(dreamtowerInput);
    else if (option.value === 'crash') setInput(crashInput);
  }, [option.value, coinflipInput, jackpotInput, dreamtowerInput, crashInput]);

  useEffect(() => {
    if (!inputRef) return;
    inputRef.current.style.height = 'auto';
    inputRef.current.style.height = inputRef.current.scrollHeight + 5 + 'px';
  }, [input]);

  useEffect(() => {
    if (!algorithmRef) return;
    algorithmRef.current.style.height = 'auto';
    algorithmRef.current.style.height =
      algorithmRef.current.scrollHeight + 5 + 'px';
  }, [gameAlgorithm]);

  useEffect(() => {
    if (!outputRef) return;
    outputRef.current.style.height = 'auto';
    outputRef.current.style.height = outputRef.current.scrollHeight + 5 + 'px';
  }, [output]);

  return (
    <Box px="4px">
      <Text
        textTransform="uppercase"
        fontSize="18px"
        lineHeight="22px"
        fontWeight={500}
        color="white"
      >
        advanced fairness verification
      </Text>

      <RowContainer mt="35px">
        Game
        <Select
          styles={customStyles}
          options={VERIFICATION_OPTIONS}
          value={option}
          isSearchable={false}
          onChange={(selectedOption: any) => {
            if (selectedOption.value === option.value) return;
            setOption(selectedOption);
            setOutput('');
          }}
        />
      </RowContainer>

      <RowContainer mt="20px">
        <Flex justifyContent="space-between">
          <Text>Algorithm Input (Editable)</Text>
          <Text fontSize="12px">JSON Object</Text>
        </Flex>
        <InputContainer background="#03060999">
          <TextAreaContainer>
            <TextArea
              ref={inputRef}
              value={input}
              onKeyDown={e => {
                if (e.keyCode === 9) {
                  e.preventDefault();
                  e.currentTarget.setRangeText(
                    '  ',
                    e.currentTarget.selectionStart,
                    e.currentTarget.selectionStart,
                    'end'
                  );
                }
              }}
              onChange={e => {
                setInput(e.target.value);
              }}
            />
          </TextAreaContainer>
          <CopyButton onClick={() => copy(input)} />
        </InputContainer>
      </RowContainer>

      <VerifyButton mt="20px" onClick={handleRun}>
        Run Algorithm
      </VerifyButton>

      <RowContainer mt="30px">
        <Flex justifyContent="space-between">
          <Text>{option.label} Algorithm</Text>
          <Text fontSize="12px">Javascript</Text>
        </Flex>
        <InputContainer background="#0306094b">
          <TextAreaContainer>
            <TextArea ref={algorithmRef} value={gameAlgorithm} readOnly />
          </TextAreaContainer>
          <CopyButton onClick={() => copy(gameAlgorithm)} />
        </InputContainer>
      </RowContainer>

      <RowContainer mt="20px" minWidth="350px">
        Algorithm Output
        <InputContainer background="#0306094b">
          <TextAreaContainer>
            <TextArea readOnly ref={outputRef} value={output} />
          </TextAreaContainer>
          <CopyButton onClick={() => copy(output)} />
          {option.value === 'dreamtower' && (
            <Flex justifyContent="center">
              <Tower
                tower={towerData.tower}
                blocksInRow={towerData.blocksInRow}
              />
            </Flex>
          )}
        </InputContainer>
      </RowContainer>

      <VerifyButton mt="23px" onClick={handleRun}>
        Run Algorithm
      </VerifyButton>
    </Box>
  );
};

export default React.memo(Verification);

const TextAreaContainer = styled(Box)`
  overflow-y: auto;
`;

const TextArea = styled.textarea`
  border: none;
  outline: none;
  background: transparent;

  font-weight: 400;
  font-size: 14px;
  line-height: 17px;
  color: white;

  resize: none;
  white-space: pre;
  width: 100%;
`;

const InputContainer = styled(Grid)`
  border-radius: 11px;
  border-bottom-right-radius: 0px;
  padding: 20px 14px 20px 25px;

  grid-template-columns: auto max-content;
  gap: 15px;
  height: auto;
`;

const RowContainer = styled(Flex)`
  flex-direction: column;
  gap: 7px;

  color: #768bad;
  font-size: 16px;
  font-weight: 400;
  line-height: 1.2;
`;

const CopyButton = styled(Button)`
  width: 28px;
  height: 28px;
  position: relative;

  background: linear-gradient(180deg, #2a3d57 0%, #2a3d57 100%);
  border-radius: 6px;

  &:after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    background: url(${copyIcon}) no-repeat;
    background-size: 12px 12px;
    background-position: center;
  }
`;

CopyButton.defaultProps = { variant: 'secondary' };

const VerifyButton = styled(Button)`
  border: 2px solid ${({ theme }) => theme.colors.success};
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0.3) 100%);
  border-radius: 7px;

  font-size: 14px;
  font-weight: 600;
  line-height: 17px;
  letter-spacing: 16%;
  color: white;

  padding: 9px 14px;

  text-transform: uppercase;
`;

VerifyButton.defaultProps = { variant: 'secondary' };
