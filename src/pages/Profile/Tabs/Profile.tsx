import React, { useCallback } from 'react';
import styled from 'styled-components';
// import { Options } from "react-select";

import {
  Avatar,
  useModal,
  Text,
  Span,
  Button,
  // Select,
  Box,
  Flex,
  Grid
} from 'components';
import { useFetchUserInfo, useFastRefreshEffect } from 'hooks';
import { useAppDispatch, useAppSelector } from 'state';

import { ReactComponent as DeathIcon } from 'assets/imgs/icons/death.svg';
import { ReactComponent as EyeIcon } from 'assets/imgs/icons/eye.svg';
import { ReactComponent as AngledSwordIcon } from 'assets/imgs/icons/angled-sword.svg';
import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
// import { ReactComponent as MsgIcon } from "assets/imgs/icons/msg.svg";
import { formatNumber, formatUserName } from 'utils/format';

import EditProfileModal from '../components/EditProfileModal';
import { api } from 'services';
import { useWallet } from '@solana/wallet-adapter-react';
import { logoutUser, requestLogin } from 'state/user/actions';
import { convertBalanceToChip } from 'utils/balance';

export default function Profile() {
  const { name, avatar, statistics, walletAddress } = useAppSelector(
    state => state.user
  );
  const { disconnect } = useWallet();
  const dispatch = useAppDispatch();
  const [onPresentEditProfile] = useModal(<EditProfileModal />);
  const { fetchUserStatistics } = useFetchUserInfo();

  useFastRefreshEffect(() => {
    fetchUserStatistics();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const logout = useCallback(async () => {
    try {
      dispatch(requestLogin(false));
      await api.get('/user/logout');
    } catch (error) {
      console.info(error);
    } finally {
      disconnect();

      dispatch(logoutUser());
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="container">
      <div className="box">
        <Flex
          flexDirection={['column', 'column', 'column', 'row']}
          justifyContent={'space-between'}
          gap={20}
        >
          <Flex alignItems={'center'} gap={16}>
            <Avatar image={avatar} size="80px" padding="2px" />

            <Box>
              <Text
                color={'#96A8C2'}
                fontSize="18px"
                fontWeight={600}
                lineHeight="28px"
              >
                {formatUserName(name)}
              </Text>
              <Text
                color={'#556273'}
                fontSize="13px"
                fontWeight={600}
                lineHeight="28px"
              >
                {formatUserName(walletAddress)}
              </Text>
            </Box>

            <Button
              border="2px solid #4F617B"
              borderRadius="0px"
              borderWidth="0px 0px 0px 2px"
              background="#070C12"
              color="#4F617B"
              fontWeight={600}
              p="8px 10px"
              ml="15px"
              nonClickable={true}
            >
              <EyeIcon />
              {statistics.private_profile
                ? 'Private Profile'
                : 'Public Profile'}
            </Button>
          </Flex>
          <Flex flexDirection={'row'} alignItems="center" gap={20}>
            <Button
              background={'#242F42'}
              borderRadius={'5px'}
              color="#768BAD"
              fontSize={'14px'}
              fontWeight={600}
              px="15px"
              py="8px"
              width={['100%', '100%', '100%', 'auto']}
              onClick={onPresentEditProfile}
            >
              Edit Profile
            </Button>

            <Button
              background={'#242F42'}
              borderRadius={'5px'}
              color="#768BAD"
              fontSize={'14px'}
              fontWeight={600}
              px="15px"
              py="8px"
              width={['100%', '100%', '100%', 'auto']}
              onClick={logout}
            >
              Sign Out
            </Button>
          </Flex>
        </Flex>

        {/* <Flex mt="40px" gap={35}>
          <AchievementCard
            background="#2B1F12"
            color="#FFAE34"
            subject="Real Gladiator"
            content="Play more than 4,000 on Duelana"
            current="4,643"
          />

          <AchievementCard
            background="#122B28"
            color="#17EF97"
            subject="Real Sniper"
            content="Won a game with less than 0.1% chance"
            current="0,09%"
          />

          <AchievementCard
            background="#1A2433"
            color="#96A8C2"
            subject="Duel Lover"
            content="plays 10 games choosing the Duel side"
            current="16"
          />
        </Flex> */}

        <Flex mt="45px" justifyContent={'space-between'}>
          <Text color="white" fontSize={'30px'} fontWeight={500}>
            Statistics
          </Text>

          {/* <Select
            background="#32425A"
            hoverBackground="#29364a"
            color="#96A8C2"
            options={STATISTICS_OPTIONS}
          /> */}
        </Flex>

        <DetailContainer>
          <Detail>
            <Text fontWeight={600} letterSpacing={'2px'}>
              TOTAL GAMES
            </Text>
            <Text
              color={'white'}
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {formatNumber(statistics.total_rounds)}
            </Text>
          </Detail>

          <Detail>
            <Text fontWeight={600} letterSpacing={'2px'}>
              GAMES WON
            </Text>
            <Text
              color="white"
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {formatNumber(statistics.winned_rounds)}
            </Text>
          </Detail>

          <Detail>
            <Text fontWeight={600} letterSpacing={'2px'}>
              GAMES LOST
            </Text>
            <Text
              color="white"
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {formatNumber(statistics.lost_rounds)}
            </Text>
          </Detail>

          <Detail>
            <Text fontWeight={600} letterSpacing={'2px'}>
              WIN RATIO
            </Text>
            <Text
              color="success"
              fontWeight={700}
              fontSize={['1.66em', '1.5em', '1.5em']}
              display="flex"
            >
              {formatNumber(
                statistics.total_rounds === 0
                  ? 0
                  : +(
                      (statistics.winned_rounds / statistics.total_rounds) *
                      100
                    )
              )}
              %
            </Text>
          </Detail>
        </DetailContainer>

        <Divider />

        <Grid
          gap={20}
          my="20px"
          gridTemplateColumns={'repeat(auto-fit, minmax(250px, 1fr))'}
          color="#96A8C2"
        >
          <Detail2>
            <Text fontWeight={600} letterSpacing={'2px'}>
              BEST & WORST STREAKS
            </Text>
            <Flex alignItems={'center'}>
              <AngledSwordIcon width={'19px'} />
              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {formatNumber(statistics.best_streaks)}
              </Span>
              <Span
                color="#4F617B"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={500}
                mx="15px"
              >
                /
              </Span>

              <DeathIcon />
              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {formatNumber(statistics.worst_streaks)}
              </Span>
            </Flex>
          </Detail2>
          <Detail2>
            <Text fontWeight={600} letterSpacing={'2px'}>
              WAGERED
            </Text>

            <Flex alignItems={'center'}>
              <CoinIcon />

              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {formatNumber(convertBalanceToChip(statistics.total_wagered))}
              </Span>
            </Flex>
          </Detail2>
        </Grid>

        {/* <Grid
          gap={20}
          gridTemplateColumns={'repeat(auto-fit, minmax(250px, 1fr))'}
          color="#96A8C2"
        >
          <Detail2>
            <Text fontWeight={600} letterSpacing={'2px'}>
              MAX PROFIT
            </Text>
            <Flex alignItems={'center'}>
              <CoinIcon />

              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {formatNumber(convertBalanceToChip(statistics.max_profit))}
              </Span>
            </Flex>
          </Detail2>
          <Detail2>
            <Text fontWeight={600} letterSpacing={'2px'}>
              TOTAL PROFIT
            </Text>

            <Flex alignItems={'center'}>
              <CoinIcon />

              <Span
                color={statistics.total_profit < 0 ? 'warning' : 'success'}
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {`${statistics.total_profit < 0 ? '' : '+'}${formatNumber(
                  convertBalanceToChip(statistics.total_profit )
                )}`}
              </Span>
            </Flex>
          </Detail2>
        </Grid> */}
      </div>
    </div>
  );
}

const Divider = styled(Box)`
  height: 2px;
  margin-top: 20px;
  background-color: #354d6d;

  .width_800 & {
    display: none;
  }
`;

const DefaultDetail = styled(Flex)`
  flex-direction: column;
  background: #1a2433;
  border-radius: 9px;
  padding: 15px;
  font-size: 12px;
  gap: 4px;
  color: #96a8c2;
`;

const Detail = styled(DefaultDetail)`
  align-items: center;

  .width_800 & {
    align-items: start;
    padding: 0px;
    border-radius: 0px;
    font-size: 16px;
  }
`;

const Detail2 = styled(DefaultDetail)`
  .width_800 & {
    padding: 25px 50px 17px;
    border-radius: 13px;
    font-size: 16px;
  }
`;

const DetailContainer = styled(Grid)`
  grid-template-columns: repeat(2, 1fr);
  background: transparent;
  border-radius: 13px;
  margin-top: 30px;
  gap: 20px;

  .width_800 & {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    padding: 15px 50px;
    background: #1a2433;
  }
`;
