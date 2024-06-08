import React, { useState, useCallback } from 'react';
import styled from 'styled-components';

import Chip1Img from 'assets/imgs/coins/chip-1.png';
import Chip2Img from 'assets/imgs/coins/chip-2.png';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';

import { Flex, Box } from 'components/Box';
import { Button } from 'components/Button';
import { InputBox } from 'components/InputBox';
import { Text, Label } from 'components/Text';
import useModal from 'components/Modal/useModal';
import { loadCoupon } from 'state/user/actions';
import { useAppDispatch } from 'state';
import { toast } from 'utils/toast';
import api from 'utils/api';

const StyledTextFlex = styled(Flex)`
  justify-content: center;
  align-items: center;
  display: none;

  ${({ theme }) => theme.mediaQueries.md} {
    justify-content: space-between;
    display: flex;
  }
`;

const StyledText = styled(Text)`
  font-size: 20px;
  font-weight: 600;
  color: white;
  letter-spacing: '0.18em';
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: block;
  }
`;

const StyledLabel = styled(Label)`
  font-weight: 600;
  color: #bad0ee;
  :hover {
    color: #4fff8b;
    cursor: pointer;
  }
`;

const StyledFlex = styled(Flex)`
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
  }
`;

export default function CouponTab() {
  const dispatch = useAppDispatch();
  const [, onDismiss] = useModal(<></>, false);
  const [couponCode, setCouponCode] = useState('');

  const onCouponChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCouponCode(e.target.value);
  };

  const handleRedeem = useCallback(
    async (e: any) => {
      e.preventDefault();
      try {
        const { data: couponInfo } = await api.post('/coupon/redeem', {
          code: couponCode
        });

        dispatch(loadCoupon(couponInfo.activeCoupon));
        onDismiss();
        toast.success('Successfully Redeemed!');
      } catch (error: any) {
        if (error.response.data.message) {
          toast.error(error.response.data.message);
        }
      }
    },
    [couponCode, dispatch, onDismiss]
  );

  return (
    <div className="container">
      <div className="box">
        <StyledTextFlex>
          <StyledText>COUPONS</StyledText>

          <CloseIcon color="#96A8C2" onClick={onDismiss} cursor="pointer" />
        </StyledTextFlex>

        <form onSubmit={handleRedeem}>
          <Flex
            flexDirection={'column'}
            position={'relative'}
            zIndex={20}
            flex={1}
          >
            <Text
              color={'#BAD0EE'}
              fontSize="16px"
              fontWeight={700}
              textAlign={'center'}
              mt="30px"
              mx={'auto'}
              maxWidth={'400px'}
            >
              Redeem you coupon code. The coupon value will be automatically
              credited to your account.
            </Text>
            <Flex justifyContent={'center'} mt="30px">
              <Box width="250px">
                <InputBox gap={20} p="10px 20px">
                  <input
                    type={'text'}
                    name="duel-points"
                    onChange={onCouponChange}
                    placeholder="Enter Coupon code..."
                  />
                </InputBox>
              </Box>
            </Flex>

            <StyledFlex justifyContent={'center'} mt="30px">
              <Button
                fontSize={'16px'}
                borderRadius="5px"
                fontWeight={600}
                p={'12px 30px'}
                type="submit"
                disabled={couponCode === ''}
              >
                Redeem Code
              </Button>
            </StyledFlex>

            <Text
              color={'#BAD0EE'}
              fontSize="16px"
              mt="25px"
              textAlign={'center'}
            >
              You may find coupons scattered in our{' '}
              <a
                href="https://discord.gg/duel"
                rel="noreferrer"
                target={'_blank'}
              >
                <StyledLabel>Discord</StyledLabel>
              </a>{' '}
              or{' '}
              <a
                href="https://twitter.com/DuelCasino"
                rel="noreferrer"
                target={'_blank'}
              >
                <StyledLabel>Twitter</StyledLabel>
              </a>
            </Text>

            <Flex
              flex={1}
              alignItems={'end'}
              justifyContent={'center'}
              mt="30px"
              display={['flex', 'flex', 'flex', 'none']}
            >
              <Button
                fontSize={'16px'}
                borderRadius="5px"
                fontWeight={600}
                p={'12px 30px'}
                type="submit"
                disabled={couponCode === ''}
              >
                Redeem Code
              </Button>
            </Flex>
          </Flex>
        </form>
        <img src={Chip1Img} className="chip-left" alt="chip-left" />
        <img src={Chip2Img} className="chip-right" alt="chip-right" />
      </div>
    </div>
  );
}
