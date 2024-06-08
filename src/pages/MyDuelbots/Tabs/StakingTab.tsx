import React, { useCallback, useMemo } from 'react';
import { InputBox } from 'components/InputBox';
import { ReactComponent as SearchIcon } from 'assets/imgs/icons/search.svg';
import { Button, Flex } from 'components';
import DuelCard from '../components/DuelCard';
import { useAppDispatch, useAppSelector } from 'state';
import { selectBot } from 'state/staking/actions';
import { StyledDuelbotsContainer } from './styles';
import useStaking from '../hooks/useStaking';

export default function StakingTab() {
  const { bots, selectedBots } = useAppSelector(state => state.staking);
  const { unStakeDuelBots } = useStaking();
  const dispatch = useAppDispatch();
  const selectDuelBots = useCallback(
    (nft: any) => {
      dispatch(selectBot(nft));
    },
    [dispatch]
  );

  const getDuelBots = useMemo(() => {
    return bots
      .filter(bot => bot.status === 'staked')
      .map(bot => (
        <DuelCard
          key={bot.mintAddress}
          mintAddress={bot.mintAddress}
          name={bot.name}
          image={bot.image}
          staked={true}
          totalEarned={bot.totalEarned}
          stakingReward={bot.stakingReward}
          selected={
            selectedBots.findIndex(
              item => item.mintAddress === bot.mintAddress
            ) !== -1
          }
          onClick={() => selectDuelBots(bot)}
        />
      ));
  }, [bots, selectedBots, selectDuelBots]);

  const handleUnstake = () => {
    unStakeDuelBots(
      bots.filter(bot => bot.status === 'staked').map(bot => bot.mintAddress)
    );
  };

  return (
    <div className="container">
      <div className="box">
        <Flex gap={30}>
          <InputBox background={'#131d28 !important'} gap={20} p="10px 20px">
            <SearchIcon width={25} height={25} />
            <input type={'text'} name="search" placeholder="Search DuelBots" />
          </InputBox>
          <Button
            background={'#242F42'}
            p="12px 30px"
            borderRadius={'5px'}
            fontWeight={600}
            color="#768BAD"
            fontSize="14px"
            onClick={handleUnstake}
            // disabled={true}
          >
            Unstake All
          </Button>
        </Flex>

        <StyledDuelbotsContainer>{getDuelBots}</StyledDuelbotsContainer>
      </div>
    </div>
  );
}
