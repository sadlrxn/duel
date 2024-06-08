import React from 'react';
import {
  LazyLoadImage,
  LazyLoadImageProps
} from 'react-lazy-load-image-component';
import styled from 'styled-components';
import { Box, Flex } from 'components/Box';
import { Button, Chip, Span, Text } from 'components';
import DuelBots from 'components/Icon/DuelBots';
import { ReactComponent as PlusIcon } from 'assets/imgs/icons/plus.svg';
import { ReactComponent as CheckedIcon } from 'assets/imgs/icons/checked.svg';
import api from 'utils/api';
import { imageProxy } from 'config';
import useStaking from '../hooks/useStaking';
import { convertBalanceToChip } from 'utils/balance';

const StyledLazyLoadImage = styled(LazyLoadImage)<LazyLoadImageProps>`
  border-radius: 6px;
  height: 100%;
  width: 100%;
  margin: auto;
  object-fit: cover;
`;

export interface DuelCardProps {
  mintAddress: string;
  name?: string;
  image: string;
  staked?: boolean;
  totalEarned?: number;
  stakingReward?: number;
  selected?: boolean;
  onClick?: () => void;
}

const StyledDuelCard = styled(Flex)<{ selected: boolean }>`
  background: #131c28;

  border-radius: 14px;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  align-items: center;
  width: 100%;
  cursor: pointer;

  ${({ selected }) =>
    selected
      ? `
    box-shadow: 0 0 0 5px #4fff8b33;
    border: 1.5px solid #4FFF8B;
  `
      : `border: 1.5px solid #1d2735;`}
`;

const StyledPlusIcon = styled(PlusIcon)`
  position: absolute;
  top: 13px;
  right: 13px;
`;

const StyledCheckedIconIcon = styled(CheckedIcon)`
  position: absolute;
  top: 13px;
  right: 13px;
`;

export default function DuelCard({
  name,
  mintAddress,
  image,
  staked = false,
  totalEarned = 0,
  stakingReward = 0,
  selected = false,
  onClick
}: DuelCardProps) {
  const { stakeDuelBots, unStakeDuelBots, claimRewards } = useStaking();

  return (
    <StyledDuelCard selected={selected} onClick={onClick}>
      <Box position={'relative'} height={225} width="100%" mb="10px">
        <StyledLazyLoadImage
          width={200}
          height={200}
          src={imageProxy(300) + image}
          alt={name}
        />
        {selected ? <StyledCheckedIconIcon /> : <StyledPlusIcon />}
      </Box>
      <Box p="8px 10px" borderRadius={'6px'} background="#1A2534" width="100%">
        <Flex alignItems="center" justifyContent={'space-between'}>
          <DuelBots />
          <Span
            color="white"
            fontSize={'12px'}
            fontWeight={600}
            textAlign="center"
          >
            {name?.slice(9)}
          </Span>
        </Flex>
        {!staked ? (
          <Button
            background={'#1A5032'}
            color="success"
            borderRadius="5px"
            fontSize={'14px'}
            fontWeight={700}
            width="100%"
            p="8px 10px"
            mt="5px"
            onClick={() => stakeDuelBots([mintAddress])}
          >
            Stake
          </Button>
        ) : (
          <Flex justifyContent={'space-between'} alignItems="end" mt="5px">
            <Box>
              <Text
                textTransform="uppercase"
                color={'#96A8C2'}
                fontSize="10px"
                fontWeight={600}
              >
                Lifetime earnings
              </Text>

              <Chip price={convertBalanceToChip(totalEarned).toFixed(2)} />
            </Box>

            <Button
              background={'#1A5032'}
              color="success"
              borderRadius="5px"
              fontSize={'14px'}
              fontWeight={700}
              p="5px 8px"
              onClick={() => {
                claimRewards([mintAddress]);
              }}
              disabled={
                convertBalanceToChip(stakingReward) < 0.1 ? true : false
              }
            >
              {convertBalanceToChip(stakingReward) < 0.1 ? (
                'Claimed'
              ) : (
                <Chip price={convertBalanceToChip(stakingReward).toFixed(2)} />
              )}
            </Button>
          </Flex>
        )}
      </Box>
    </StyledDuelCard>
  );
}
