import React, { useState, useMemo, useEffect, useRef } from 'react';
import { ClipLoader } from 'react-spinners';

import { Box, Flex } from 'components/Box';
import { Button } from 'components/Button';
import { ReactComponent as EqualIcon } from 'assets/imgs/icons/equal.svg';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';
import { InputBox } from 'components/InputBox';
import { Text } from 'components/Text';
import Select from 'components/Select';
import { toast } from 'utils/toast';
import useTrading from 'hooks/useTrading';
import { formatSol2Lamports, formatSplToken2Lamports } from 'utils/format';
import { useModal } from 'components/Modal';
import { CoinIcon } from 'components/Icon';
import { Chip } from 'components/Chip';
import styled from 'styled-components';
import { useAppSelector } from 'state';

import { Icon, SelectItem } from './Crypto';
import { useMatchBreakpoints } from 'hooks';
import useUserTokenAmounts from 'hooks/useUserTokenAmounts';

const StyledText = styled(Text)`
  font-size: 20px;
  font-weight: 600;
  color: white;
  letter-spacing: 0.18em;
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: block;
  }
`;

const StyledTextFlex = styled(Flex)`
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
  display: none;

  ${({ theme }) => theme.mediaQueries.md} {
    justify-content: space-between;
    display: flex;
  }
`;

const StyledInputFlex = styled(Flex)`
  flex-direction: column;
  justify-content: space-between;
  margin: 15px 0px;
  align-items: center;
  gap: 30px;

  ${({ theme }) => theme.mediaQueries.md} {
    flex-direction: row;
  }
`;

const StyledInputBox = styled(InputBox)`
  width: 100%;
  input {
    width: 150px;
    ${({ theme }) => theme.mediaQueries.md} {
      width: auto;
    }
  }
`;

export default function DepositSolTab() {
  const { depositToken } = useTrading();
  const { tokens, prices } = useAppSelector(state => state.token);
  const { isMobile } = useMatchBreakpoints();

  const firstInputRef = useRef<HTMLInputElement>(null);
  const tokenBalances = useUserTokenAmounts();

  const options = useMemo(() => {
    return tokens.map((token, index) => {
      return {
        label: SelectItem(
          token.keyword.toUpperCase(),
          token.image,
          24,
          tokenBalances[index]
        ),
        value: token.keyword,
        img: token.image,
        mintAddress: token.mintAddress,
        decimals: token.decimals,
        balance: tokenBalances[index],
        index
      };
    });
  }, [tokens, tokenBalances]);

  const [formData, setFormData] = useState({
    solAmount: '',
    duelAmount: ''
  });
  const [loading, setLoading] = useState(false);
  const [optionIndex, setOptionIndex] = useState(0);
  const [crypto, setCrypto] = useState(options[0].value);

  const option = useMemo(() => options[optionIndex], [options, optionIndex]);

  useEffect(() => {
    setTimeout(() => {
      if (!firstInputRef || !firstInputRef.current) return;
      firstInputRef.current.focus();
    }, 100);
  }, [optionIndex]);

  const tokenPrice = useMemo(() => {
    //@ts-ignore
    return prices[crypto] ?? 0;
  }, [prices, crypto]);

  const [, onDismissDepositSol] = useModal(<></>, false);

  const handleChange = (e: any) => {
    // eslint-disable-next-line no-throw-literal
    if (!tokenPrice) throw "can't fetch sol price!";

    if (e.target.value === '') {
      setFormData({ solAmount: '', duelAmount: '' });
      return;
    }

    if (e.target.name === 'solAmount') {
      const solAmount = e.target.value;
      const duelAmount = (parseFloat(e.target.value) * tokenPrice)
        .toFixed(2)
        .toString();

      setFormData({ solAmount, duelAmount });
    } else {
      const duelAmount = e.target.value;
      const solAmount = (parseFloat(e.target.value) / tokenPrice).toString();

      setFormData({ solAmount, duelAmount });
    }
  };

  const handleDeposit = async (e: any) => {
    if (loading) return;
    e.preventDefault();
    if (formData.solAmount === '') {
      toast.warning('Input amount!');
      return;
    }
    setLoading(true);
    if (crypto.toLowerCase().slice(0, 3) === 'sol')
      await depositToken(
        formatSol2Lamports(parseFloat(formData.solAmount)),
        option.decimals,
        tokenPrice
      );
    else
      await depositToken(
        formatSplToken2Lamports(
          parseFloat(formData.solAmount),
          option.decimals
        ),
        option.decimals,
        tokenPrice,
        option.mintAddress
      );

    setLoading(false);
    onDismissDepositSol();
  };

  const handleMax = () => {
    const token = ['SOL', 'Bonk', 'USDC'];

    const balance = tokenBalances[token.indexOf(crypto)].toString();

    const chip = (parseFloat(balance) * tokenPrice).toFixed(2).toString();

    setFormData({ solAmount: balance, duelAmount: chip });
  };

  useEffect(() => {
    setFormData({ solAmount: '', duelAmount: '' });
  }, [option]);

  return (
    <div className="container">
      <div className="box">
        <form onSubmit={handleDeposit}>
          <StyledTextFlex>
            <StyledText>DEPOSIT</StyledText>

            <CloseIcon
              color="#96A8C2"
              onClick={onDismissDepositSol}
              cursor="pointer"
            />
          </StyledTextFlex>

          <Flex
            alignItems="center"
            flexDirection={isMobile ? 'column' : 'row'}
            gap={isMobile ? 18 : 34}
          >
            <Select
              // defaultValue={options[0]}

              value={options[optionIndex]}
              options={options}
              isSearchable={false}
              fontSize="16px"
              fontWeight={400}
              background="#0A1119"
              width={isMobile ? '100%' : 330}
              hoverBackground="#0A1119"
              color="#FFF"
              onChange={(selectedOption: any) => {
                setOptionIndex(selectedOption.index);
                setCrypto(selectedOption.value);
              }}
            />
            <Text color={'#B2D1FF'} fontWeight={600}>
              1 {crypto} ={' '}
              <Chip
                color="#B2D1FF"
                price={tokenPrice && tokenPrice!.toFixed(2)}
              />{' '}
              = ${tokenPrice && tokenPrice!.toFixed(9)} USD
            </Text>
          </Flex>

          <Box p={'25px'} borderRadius="12px" background={'#0F1B2B'} mt="25px">
            <Text color={'#BAD0EE'}>VALUE CALCULATOR</Text>

            <StyledInputFlex>
              <StyledInputBox gap={20} p="4px 4px 4px 10px">
                <Icon img={option.img} />
                <input
                  ref={firstInputRef}
                  type={'number'}
                  name="solAmount"
                  value={formData.solAmount}
                  onChange={handleChange}
                  placeholder={`Enter ${crypto} amount`}
                />
                <Button
                  background={'#1A293D'}
                  color="#768BAD"
                  p="10px"
                  borderRadius="5px"
                  onClick={handleMax}
                >
                  Max
                </Button>
              </StyledInputBox>
              <EqualIcon />
              <StyledInputBox gap={20} p="4px 4px 4px 10px">
                <CoinIcon />
                <input
                  type={'number'}
                  name="duelAmount"
                  value={formData.duelAmount}
                  onChange={handleChange}
                  placeholder="Enter CHIP amount"
                />
                <Button
                  background={'#1A293D'}
                  color="#768BAD"
                  p="10px"
                  borderRadius="5px"
                  onClick={handleMax}
                >
                  Max
                </Button>
              </StyledInputBox>
            </StyledInputFlex>

            <Text color={'#4D6384'} fontWeight={500} mt="30px">
              The value of Solana may vary between now and the time we receive
              your payment.
            </Text>
          </Box>

          <Flex flex={1} alignItems={'end'} justifyContent={'center'} mt="30px">
            <Button
              fontSize={'16px'}
              fontWeight={600}
              p={'12px 40px'}
              borderRadius="5px"
              type="submit"
              disabled={loading || formData.solAmount === ''}
            >
              {loading ? <ClipLoader color="#fff" size={20} /> : 'Deposit'}
            </Button>
          </Flex>
        </form>
      </div>
    </div>
  );
}
