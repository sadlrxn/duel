import React, { useEffect, useMemo } from 'react';
import {
  Badge,
  Box,
  Button,
  Chip,
  DepositWithdrawModal,
  Flex,
  Grid,
  Span,
  Text,
  useModal
} from 'components';
import { useNavigate } from 'react-router-dom';
import useSWR from 'swr';
import { StyledTabs, StyledTabPanel } from './styles';
import { Tab, TabList } from 'react-tabs';
import { BotsImg, TopBox } from './styles';
import botsImg from 'assets/imgs/staking/duelbots.png';
import DuelBotsTab from './Tabs/DuelBotsTab';
import StakingTab from './Tabs/StakingTab';
import Coin from 'components/Icon/Coin';
import styled from 'styled-components';
import { useAppDispatch, useAppSelector } from 'state';
import { cancelSelectedBots, changeTab, loadBots } from 'state/staking/actions';
import api from 'utils/api';
import UnstakeModal from './components/UnstakeModal';
import useStaking from './hooks/useStaking';
import { formatNumber } from 'utils/format';
import { imageProxy } from 'config';
import { convertBalanceToChip } from 'utils/balance';

const StyledStakingFlex = styled(Flex)`
  flex-direction: column;
  gap: 15px;
  .width_800 & {
    flex-direction: row;
    gap: 20px;
  }
`;

const StyledNft = styled.img`
  width: 42px;
  height: 42px;
  border-radius: 10px;

  &:not(:first-child) {
    margin-left: -18px;
  }
`;

const BottomBox = styled(Flex)`
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 18px 30px;
  background: #1a2534;
  justify-content: space-between;
  flex-direction: column;
  gap: 15px;

  .width_800 & {
    position: absolute;
    flex-direction: row;
  }
`;

export default function Staking({ tabIndex }: { tabIndex: number }) {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { selectedBots, bots } = useAppSelector(state => state.staking);
  const { stakeDuelBots, claimRewards } = useStaking();
  const { data: botsInfo } = useSWR(`/bot/duel-bots`, async arg =>
    api.get(arg).then(res => res.data)
  );

  useEffect(() => {
    dispatch(changeTab(tabIndex === 0 ? 'STAKING' : 'UNSTAKING'));
  }, [tabIndex, dispatch]);

  useEffect(() => {
    if (!botsInfo) return;
    dispatch(loadBots(botsInfo));
  }, [botsInfo, dispatch]);

  const stakingInfo = useMemo(() => {
    let totalStakingRewards = 0;
    let stakedDuelBots = 0;
    let totalEarned = 0;
    bots
      .filter((bot: any) => bot.status === 'staked')
      .forEach((bot: any) => {
        totalStakingRewards += bot.stakingReward;
        totalEarned += bot.totalEarned;
        stakedDuelBots++;
      });

    return {
      stakedDuelBots,
      totalStakingRewards: convertBalanceToChip(totalStakingRewards),
      totalEarned: convertBalanceToChip(totalEarned)
    };
  }, [bots]);

  const [onPresentDepositNft] = useModal(
    <DepositWithdrawModal tabIndex={2} hideCloseButton={true} />,
    false
  );

  const [onPresentWithdrawNft] = useModal(
    <DepositWithdrawModal tabIndex={3} hideCloseButton={true} />,
    false
  );

  const StakingContent = (
    <StyledStakingFlex
      padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}
    >
      <TopBox>
        <Box position={'relative'} zIndex={1}>
          <Text color="#fff" fontWeight={700} fontSize={27}>
            Want daily free CHIPS?
          </Text>
          <Text color={'#b3b3b3'} fontWeight={600} fontSize={16}>
            Get a slice of the pie with your very own Duelbot.
          </Text>

          <a
            href="https://magiceden.io/marketplace/duelbots"
            target={'_blank'}
            rel="noreferrer"
          >
            <Button
              borderRadius={'5px'}
              p="8px 35px"
              fontWeight={700}
              fontSize={14}
              mt="20px"
            >
              Get a DuelBot
            </Button>
          </a>
        </Box>

        <BotsImg src={botsImg} />
      </TopBox>

      <Box>
        <StyledStakingFlex mb="15px">
          <Button
            background={'#1A5032'}
            p="12px 15px"
            color="#4FFF8B"
            fontSize={14}
            fontWeight={600}
            borderRadius="5px"
            onClick={onPresentDepositNft}
          >
            Deposit DuelBot
          </Button>

          <Button
            background={'#242F42'}
            p="12px 15px"
            color="#768BAD"
            fontSize={14}
            fontWeight={600}
            borderRadius="5px"
            onClick={onPresentWithdrawNft}
          >
            Withdraw DuelBot
          </Button>
        </StyledStakingFlex>
        <Flex
          background={'#121A25'}
          p="22px 30px"
          alignItems={'end'}
          gap={15}
          borderRadius="12px"
          justifyContent={'space-between'}
        >
          <Box>
            <Text
              color={'#96A8C2'}
              fontSize="12px"
              fontWeight={600}
              mb="10px"
              letterSpacing={'0.185em'}
              textTransform="uppercase"
            >
              Available to claim
            </Text>
            <Chip price={formatNumber(stakingInfo.totalStakingRewards)} />
          </Box>
          <Button
            background={'#1A5032'}
            p="6px 13px"
            color="#4FFF8B"
            borderRadius={'5px'}
            disabled={stakingInfo.totalStakingRewards < 0.1 ? true : false}
            onClick={() =>
              claimRewards(
                bots
                  .filter((bot: any) => bot.status === 'staked')
                  .map((bot: any) => bot.mintAddress)
              )
            }
          >
            Claim
          </Button>
        </Flex>
      </Box>
    </StyledStakingFlex>
  );

  const UnstakingContent = (
    <Grid gridTemplateColumns={'repeat(auto-fit, minmax(250px, 1fr))'} gap={20}>
      <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
        <Text
          textTransform="uppercase"
          color={'#96A8C2'}
          fontSize="12px"
          fontWeight={600}
          letterSpacing="0.185em"
        >
          DUELBOTS STAKED
        </Text>

        <Text
          textTransform="uppercase"
          color={'#fff'}
          fontSize="20px"
          fontWeight={700}
          mt="10px"
        >
          {stakingInfo.stakedDuelBots} of {bots && bots.length}
        </Text>
      </Box>

      <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
        <Text
          textTransform="uppercase"
          color={'#96A8C2'}
          fontSize="12px"
          fontWeight={600}
          letterSpacing="0.185em"
        >
          TOTAL EARNED
        </Text>
        <Flex alignItems={'center'} gap={5} mt="10px">
          <Coin />
          <Text
            textTransform="uppercase"
            color={'#fff'}
            fontSize="20px"
            fontWeight={700}
          >
            {formatNumber(stakingInfo.totalEarned)}
          </Text>
        </Flex>
      </Box>

      <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
        <Text
          textTransform="uppercase"
          color={'#96A8C2'}
          fontSize="12px"
          fontWeight={600}
          letterSpacing="0.185em"
        >
          Available to claim
        </Text>
        <Flex alignItems={'center'} gap={5} mt="10px">
          <Coin />
          <Text
            textTransform="uppercase"
            color={'#fff'}
            fontSize="20px"
            fontWeight={700}
          >
            {formatNumber(stakingInfo.totalStakingRewards)}
          </Text>

          <Button
            background={'#1A5032'}
            p="7px 25px"
            borderRadius={'5px'}
            fontWeight={600}
            color="#4FFF8B"
            fontSize="14px"
            ml="50px"
            disabled={stakingInfo.totalStakingRewards < 0.1 ? true : false}
            onClick={() =>
              claimRewards(
                bots
                  .filter((bot: any) => bot.status === 'staked')
                  .map((bot: any) => bot.mintAddress)
              )
            }
          >
            Claim
          </Button>
        </Flex>
      </Box>
    </Grid>
  );

  const Content = useMemo(
    () => (tabIndex === 0 ? StakingContent : UnstakingContent),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [tabIndex, bots]
  );

  const [onUnstake] = useModal(<UnstakeModal />);

  const handleStaking = () => {
    if (tabIndex === 0) stakeDuelBots(selectedBots.map(bot => bot.mintAddress));
    else onUnstake();
  };

  // api.post('/bot/unstake', { mints: selectedBots.map(bot => bot.mint) });

  const handleCancel = () => {
    dispatch(cancelSelectedBots());
  };

  const handleSelect = (index: number) => {
    const tabs = ['myduelbots', 'staking'];
    navigate(`/duelbots/${tabs[index]}`);
  };

  return (
    <div>
      <Text color={'#768BAD'} fontWeight={600} fontSize={22} mb="30px">
        DUELBOTS
      </Text>
      {Content}

      <Box my="50px">
        <StyledTabs selectedIndex={tabIndex} onSelect={handleSelect}>
          <TabList>
            <Flex>
              <Tab>
                My DUELBOTS
                <b />
              </Tab>
              <Tab>
                STAKING
                <b />
              </Tab>
            </Flex>
          </TabList>
          <StyledTabPanel>
            <DuelBotsTab />
          </StyledTabPanel>
          <StyledTabPanel>
            <StakingTab />
          </StyledTabPanel>
        </StyledTabs>
      </Box>

      <BottomBox>
        <Flex alignItems="center" gap={15}>
          <Badge
            background={'#18FF8E'}
            color="black"
            borderRadius={'50%'}
            p="3px 10px"
          >
            {selectedBots.length}
          </Badge>

          <Span
            color={'#687E9D'}
            fontSize="12px"
            fontWeight={500}
            width="52px"
            p="3px"
          >
            Selected Duelbots
          </Span>

          <Flex>
            {selectedBots.map(value => (
              <StyledNft
                src={imageProxy(300) + value.image}
                alt="nft"
                key={value.mintAddress}
              />
            ))}
          </Flex>
        </Flex>

        {selectedBots.length > 0 && (
          <Flex gap={20}>
            <Button
              background={'#242F42'}
              color="#768BAD"
              borderRadius="5px"
              fontSize={'14px'}
              fontWeight={600}
              p="12px 20px"
              onClick={handleCancel}
            >
              Cancel
            </Button>

            <Button
              background={'#4FFF8B'}
              color="black"
              borderRadius="5px"
              fontSize={'14px'}
              fontWeight={600}
              p="12px 20px"
              width={'100%'}
              onClick={handleStaking}
              // disabled={activeTab === 1 ? true : false}
            >
              {tabIndex === 0 ? 'Stake Selected' : 'Unstake Selected'}
            </Button>
          </Flex>
        )}
      </BottomBox>
    </div>
  );
}
