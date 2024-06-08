import { useState, useCallback, useMemo } from 'react';
import {
  useFloating,
  autoUpdate,
  offset,
  flip,
  shift,
  useDismiss,
  useRole,
  useHover,
  useInteractions,
  FloatingFocusManager,
  useId,
  safePolygon
} from '@floating-ui/react';
import styled from 'styled-components';
import { Text, Span } from 'components/Text';
import { Box, Flex } from 'components/Box';
import { Button } from 'components/Button';
import { lightColors } from '../../theme';
import { CountDown } from 'components/CountDown';

import CoinBlueImg from 'assets/imgs/coins/coin-blue.svg';
import CoinImg from 'assets/imgs/coins/coin.svg';
import { useAppSelector } from 'state';
import api from 'utils/api';
import { updateBalance, deleteCoupon } from 'state/user/actions';
import { useAppDispatch } from 'state';
import { toast } from 'utils/toast';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

const CouponButton = styled.div`
  position: relative;
  display: flex;
  height: 38px;
  align-items: center;
  background: linear-gradient(90deg, #004150 0%, rgba(0, 65, 80, 0.25) 100%);
  border: none;
  ::before {
    position: absolute;
    background: #4be9ff;
    content: '';
    display: block;
    height: 45%;
    left: 0;
    pointer-events: none;
    width: 2px;
  }

  border-radius: 5px;
  padding: 0px 12px;

  img {
    width: 14px;
    height: 14px;
    margin-right: 8px;
  }

  span {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 600;
    font-size: 16px;
    line-height: 19px;

    color: #caffff;
  }
  &:hover {
    cursor: pointer;
  }
`;

const CouponClaimBtn = styled(Button)`
  display: flex;
  padding: 7px 20px;
  margin: 13px auto 0;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background: #14587d;
  border-radius: 5px;
  border: none;
  gap: 4px;
`;

const CouponModal = styled.div`
  position: relative;
  max-width: 195px;
  background: linear-gradient(180deg, #6a7f9e 0%, rgba(106, 127, 158, 0) 100%);
  text-align: center;
  border-top-left-radius: 10px;
  border-top-right-radius: 10px;
  border-bottom-left-radius: 6px;
  border-bottom-right-radius: 6px;
  padding: 1px 1px 0 1px;
`;

const CouponInner = styled.div`
  padding: 14px 12px 0 12px;
  background: linear-gradient(
    180deg,
    rgba(19, 32, 49, 0.9) 0%,
    rgba(26, 41, 60, 0.92) 46.47%
  );
  backdrop-filter: blur(10px);
  border-top-left-radius: 10px;
  border-top-right-radius: 10px;
  border-bottom-left-radius: 6px;
  border-bottom-right-radius: 6px;
  overflow: hidden;
`;

const ProgressBarContainer = styled(Box)`
  position: absolute;
  left: -1px;
  bottom: 0;
  background: #14587d;
  width: calc(100% + 2px);
  height: 4px;
`;

const ProgressBarContent = styled(Box)`
  height: 100%;
  background: #4be9ff;
`;

const Triangle = styled.div`
  position: absolute;
  left: 50%;
  top: 1px;
  z-index: 1;

  transform: translate(-50%, -50%) rotate(-45deg);

  border: 1px solid rgba(106, 127, 158);
  border-left-width: 0px;
  border-bottom-width: 0px;
  border-radius: 2px;
  background: rgba(26, 41, 60);

  width: 10px;
  height: 10px;
`;

export default function CouponPopover() {
  const dispatch = useAppDispatch();
  const [open, setOpen] = useState(false);

  const { balances, betBalanceType, config } = useAppSelector(
    state => state.user
  );

  const {
    couponCode,
    couponBalance,
    couponExpiredTime,
    couponWagered,
    couponWagerLimit
  } = useMemo(() => {
    const coupon = balances.coupon;
    return {
      couponCode: coupon.code,
      couponBalance: coupon.balance,
      couponExpiredTime: coupon.expiredTime
        ? coupon.expiredTime
        : Date.now() - 1000000,
      couponWagered: coupon.wagered ? coupon.wagered : 0,
      couponWagerLimit: coupon.wagerLimit ? coupon.wagerLimit : 0
    };
  }, [balances.coupon]);

  const { x, y, refs, strategy, context } = useFloating({
    open,
    onOpenChange: setOpen,
    middleware: [
      offset(10),
      flip({ fallbackAxisSideDirection: 'end' }),
      shift()
    ],
    whileElementsMounted: autoUpdate
  });

  const hover = useHover(context, {
    handleClose: safePolygon()
  });
  const dismiss = useDismiss(context);
  const role = useRole(context);

  const { getReferenceProps, getFloatingProps } = useInteractions([
    hover,
    dismiss,
    role
  ]);

  const headingId = useId();

  const handleCouponClaim = useCallback(async () => {
    try {
      const { data: couponClaim } = await api.post('/coupon/claim', {
        code: couponCode
      });
      dispatch(deleteCoupon());
      if (couponClaim.claimed >= config.couponMaxClaim) {
        dispatch(
          updateBalance({
            type: 1,
            usdAmount: config.couponMaxClaim
          })
        );
      } else {
        dispatch(
          updateBalance({
            type: 1,
            usdAmount: couponClaim.claimed
          })
        );
      }

      toast.success('Successfully Claimed!');
    } catch (error: any) {
      if (error.response.data.message) {
        toast.error(error.response.data.message);
      }
    }
  }, [couponCode, dispatch]);

  return (
    <>
      <CouponButton
        ref={refs.setReference}
        {...getReferenceProps()}
        style={
          couponBalance <= 0 || betBalanceType === 'chip'
            ? { opacity: 0.6 }
            : {}
        }
      >
        <img src={CoinBlueImg} alt="coin" />
        <span>
          {(
            Math.floor(convertBalanceToChip(couponBalance) * 100) / 100
          ).toFixed(2)}
        </span>
      </CouponButton>
      {open && (
        <FloatingFocusManager context={context} modal={false}>
          <CouponModal
            className="Popover"
            ref={refs.setFloating}
            style={{
              position: strategy,
              top: y ?? 0,
              left: x ?? 0
            }}
            aria-labelledby={headingId}
            {...getFloatingProps()}
          >
            <Triangle />
            <CouponInner>
              <Text
                color={lightColors.text}
                fontSize={11}
                fontWeight={500}
                lineHeight="14px"
                margin={0}
                id={headingId}
              >
                Wager your bonus balance {config.couponWagerTimes}x of the
                amount within {config.couponLifeTime} hours to claim it into
                your CHIP balance.
              </Text>
              {couponWagered >= couponWagerLimit && (
                <CouponClaimBtn
                  onClick={handleCouponClaim}
                  disabled={couponBalance <= 0}
                >
                  <Span color="#4BE9FF" fontWeight={600} fontSize={12}>
                    Claim
                  </Span>
                  <img src={CoinImg} alt="Coin" />
                  <Span color="#4BE9FF" fontWeight={600} fontSize={12}>
                    {couponBalance >= config.couponMaxClaim
                      ? convertBalanceToChip(config.couponMaxClaim)
                      : convertBalanceToChip(couponBalance).toFixed(2)}
                  </Span>
                </CouponClaimBtn>
              )}

              <Text
                color={lightColors.text}
                fontSize={10}
                fontWeight={600}
                marginTop="15px"
                marginBottom="1px"
              >
                TIME REMAINING
              </Text>
              <CountDown endedAt={couponExpiredTime} />
              <Text
                color={lightColors.text}
                fontSize={10}
                fontWeight={600}
                marginTop="15px"
                marginBottom="5px"
              >
                WAGER REMAINING
              </Text>
              <Flex
                justifyContent="center"
                alignItems="center"
                gap={4}
                paddingBottom="10px"
              >
                <img src={CoinBlueImg} alt="coin" />
                <Span color="#CAFFFF" fontWeight={600} fontSize={12}>
                  {convertBalanceToChip(couponWagered).toFixed(2)} /{' '}
                  {convertBalanceToChip(couponWagerLimit).toFixed(2)}
                </Span>
              </Flex>
              <ProgressBarContainer>
                <ProgressBarContent
                  width={`${(couponWagered / couponWagerLimit) * 100}%`}
                />
              </ProgressBarContainer>
            </CouponInner>
          </CouponModal>
        </FloatingFocusManager>
      )}
    </>
  );
}
