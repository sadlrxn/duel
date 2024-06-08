import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { shallowEqual } from 'react-redux';
import styled from 'styled-components';
import { useProSidebar } from 'react-pro-sidebar';
import Toggle from 'react-toggle';
import 'react-toggle/style.css';

import CoinImg from 'assets/imgs/coins/coin.svg';
import LogOutImg from 'assets/imgs/icons/logout.svg';
import NotificationImg from 'assets/imgs/icons/notification.svg';
import { ReactComponent as SidebarShowIcon } from 'assets/imgs/icons/sidebar-show.svg';
import { ReactComponent as SidebarHideIcon } from 'assets/imgs/icons/sidebar-hide.svg';

import { useAppDispatch, useAppSelector } from 'state';
import { readLog } from 'state/actions';
import { deleteCoupon, isLedger, toggleSound } from 'state/user/actions';
import { useLocalStorage, useMatchBreakpoints } from 'hooks';
import { formatUserName } from 'utils/format';

import { Flex } from 'components/Box';
import { CouponPopover } from 'components/Popover';
import { useModal } from 'components/Modal';
import {
  VolumeOn,
  VolumeOff,
  ArrowDownIcon,
  ArrowUpIcon,
  RewardsIcon
} from 'components/Icon';
import Avatar from 'components/Avatar';
import {
  TermsConditionsModal,
  LogModal,
  DepositWithdrawModal
} from 'components/Modals';
import { Text } from 'components';
import Logo from 'components/Icon/Logo';
import MobileLogo from 'components/Icon/MobileLogo';
import { convertBalanceToChip } from 'utils/balance';

import {
  ConnectWalletBtn,
  DepositBtn,
  IconButton,
  StyledHeader,
  UserBalance,
  UserName,
  WithdrawBtn,
  ToggleBtn,
  ButtonContainer,
  StyledLink,
  StyledAvatarContainer,
  RewardsBtn
} from './styles';

const StyledToggle = styled(Toggle)`
  .react-toggle-track {
    width: 40px;
    height: 20px;
    background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0) 162.5%);
    background-color: transparent !important;
  }

  &.react-toggle--checked .react-toggle-thumb {
    background: #4fff8b;
  }

  .react-toggle-thumb {
    width: 18px;
    height: 18px;
    /* background: #4fff8b; */
  }
`;

const StyledFlex = styled(Flex)`
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
  }
`;

const Appbar: FC<{
  login: () => Promise<void>;
  logout: () => Promise<void>;
}> = ({ login, logout }) => {
  const navigate = useNavigate();
  const { toggleSidebar, toggled } = useProSidebar();
  const { isAuthenticated, name, avatar, sound, usingLedger, betBalanceType } =
    useAppSelector(state => state.user);
  const balance = useAppSelector(
    state => state.user.balances.chip.balance,
    shallowEqual
  );
  const coupon = useAppSelector(
    state => state.user.balances.coupon,
    shallowEqual
  );

  const { couponExpiredTime, couponBalance } = useMemo(() => {
    return {
      couponExpiredTime: coupon.expiredTime
        ? coupon.expiredTime
        : Date.now() - 1000,
      couponBalance: coupon.balance
    };
  }, [coupon]);

  const { showDot } = useAppSelector(state => state.log);
  const dispatch = useAppDispatch();
  const { isMobile } = useMatchBreakpoints();

  const [showCoupon, setShowCoupon] = useState(false);

  useEffect(() => {
    if (couponExpiredTime < Date.now()) {
      dispatch(deleteCoupon());
      return;
    }
    const interval = setInterval(() => {
      setShowCoupon(couponExpiredTime >= Date.now());
    }, 1000);

    return () => {
      clearInterval(interval);
    };
  }, [couponExpiredTime, dispatch, showCoupon]);

  const [userAccepted, setUserAccepted] = useLocalStorage('terms', false);
  // const { pathname } = useLocation();

  const [onPresentDeposit] = useModal(
    <DepositWithdrawModal tabIndex={0} hideCloseButton={true} />,
    true
  );
  const [onPresentWithdraw] = useModal(
    <DepositWithdrawModal tabIndex={1} hideCloseButton={true} />,
    true
  );
  const [onPresentLog] = useModal(<LogModal />, true);

  const [onPresentTermsConditions] = useModal(
    <TermsConditionsModal setUserAccepted={setUserAccepted} login={login} />
  );

  const handleRewards = useCallback(() => {
    navigate('/profile?tab=rewards');
  }, [navigate]);

  const handleToggleSound = useCallback(() => {
    dispatch(toggleSound());
  }, [dispatch]);

  const handleToggle = (e: any) => {
    dispatch(isLedger(e.target.checked));
  };

  return (
    <StyledHeader>
      <Flex alignItems="center">
        <ToggleBtn onClick={toggleSidebar}>
          {toggled ? <SidebarHideIcon /> : <SidebarShowIcon />}
        </ToggleBtn>
        <StyledLink to={'/'}>{isMobile ? <MobileLogo /> : <Logo />}</StyledLink>
      </Flex>

      {isAuthenticated ? (
        <Flex alignItems="center" gap={12}>
          {showCoupon && <CouponPopover />}

          <UserBalance
            onClick={isMobile ? onPresentDeposit : () => {}}
            style={
              couponBalance > 0 && betBalanceType === 'coupon'
                ? { opacity: 0.6 }
                : {}
            }
          >
            <img src={CoinImg} alt="coin" />
            <span>
              {(Math.floor(convertBalanceToChip(balance) * 100) / 100).toFixed(
                2
              )}
            </span>
          </UserBalance>

          <ButtonContainer>
            <DepositBtn onClick={onPresentDeposit}>
              <ArrowDownIcon />
              Deposit
            </DepositBtn>
            <WithdrawBtn onClick={onPresentWithdraw}>
              <ArrowUpIcon />
              Withdraw
            </WithdrawBtn>
          </ButtonContainer>

          <RewardsBtn onClick={handleRewards} ml="18px">
            <RewardsIcon />
            Rewards
          </RewardsBtn>

          <Link to={'/profile?tab=profile'}>
            <StyledAvatarContainer>
              <Avatar
                image={avatar}
                padding="0px"
                size="46px"
                border="2px solid #768BAD"
                background="#0e1925"
              />

              <UserName>{formatUserName(name)}</UserName>
            </StyledAvatarContainer>
          </Link>

          <ButtonContainer>
            <IconButton onClick={handleToggleSound}>
              {sound ? <VolumeOn /> : <VolumeOff />}
            </IconButton>
            <IconButton
              onClick={() => {
                onPresentLog();
                dispatch(readLog());
              }}
            >
              <img src={NotificationImg} alt="notification" />
              {showDot && <span className="badge" />}
            </IconButton>
            <IconButton onClick={logout}>
              <img src={LogOutImg} alt="logout" />
            </IconButton>
          </ButtonContainer>
        </Flex>
      ) : (
        <Flex alignItems={'center'} gap={20}>
          <StyledFlex alignItems={'center'} gap={10}>
            <Text color={'#4F617B'} fontSize="14px" fontWeight={600}>
              Using Ledger?
            </Text>
            <StyledToggle
              defaultChecked={usingLedger}
              icons={false}
              onChange={handleToggle}
            />
          </StyledFlex>

          <ConnectWalletBtn
            onClick={userAccepted ? login : onPresentTermsConditions}
          >
            Connect Wallet
          </ConnectWalletBtn>
        </Flex>
      )}
    </StyledHeader>
  );
};

export default Appbar;
