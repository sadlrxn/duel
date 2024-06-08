import React, { useCallback, useMemo } from 'react';
import { InputBox } from 'components/InputBox';
import { ReactComponent as SearchIcon } from 'assets/imgs/icons/search.svg';
import { Button, Flex } from 'components';
import DuelCard from '../components/DuelCard';
import { useAppDispatch, useAppSelector } from 'state';
import { selectBot } from 'state/staking/actions';
import { StyledDuelbotsContainer } from './styles';
import useStaking from '../hooks/useStaking';

export default function DuelBotsTab() {
  const { bots, selectedBots } = useAppSelector(state => state.staking);
  const { stakeDuelBots } = useStaking();
  const dispatch = useAppDispatch();
  const selectDuelBots = useCallback(
    (nft: any) => {
      dispatch(selectBot(nft));
    },
    [dispatch]
  );

  const getDuelBots = useMemo(() => {
    return bots
      .filter(bot => bot.status === 'normal')
      .map((bot, index) => (
        <DuelCard
          key={index}
          mintAddress={bot.mintAddress}
          name={bot.name}
          image={bot.image}
          selected={
            selectedBots.findIndex(
              item => item.mintAddress === bot.mintAddress
            ) !== -1
          }
          onClick={() => selectDuelBots(bot)}
        />
      ));
  }, [bots, selectedBots, selectDuelBots]);

  const handleStake = () => {
    stakeDuelBots(
      bots.filter(bot => bot.status === 'normal').map(bot => bot.mintAddress)
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
            background={'#1A5032'}
            p="12px 30px"
            borderRadius={'5px'}
            fontWeight={600}
            color="#4FFF8B"
            fontSize="14px"
            onClick={handleStake}
          >
            Stake All
          </Button>
        </Flex>

        <StyledDuelbotsContainer>{getDuelBots}</StyledDuelbotsContainer>
      </div>
    </div>
  );
}
