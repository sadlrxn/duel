import React from 'react';
import styled from 'styled-components';
import { Box, Button, Flex, Text } from 'components';
import { Modal, ModalProps, useModal } from 'components/Modal';
import DuelBotSVG from './DuelBotSVG';
import { useAppSelector } from 'state';
import api from 'utils/api';
import useStaking from '../hooks/useStaking';

const StyledBox = styled(Box)`
  background: linear-gradient(180deg, #132031 0%, #1a293c 100%);
  padding: 45px;
  border-radius: 20px;
`;

export default function UnstakeModal({ tabIndex, ...props }: ModalProps) {
  const { stakeDuelBots, unStakeDuelBots } = useStaking();
  const { selectedBots } = useAppSelector(state => state.staking);
  const [, hide] = useModal(<></>);

  const handleStaking = async () => {
    await unStakeDuelBots(selectedBots.map(bot => bot.mintAddress));
    hide();
  };

  return (
    <Modal {...props}>
      <StyledBox>
        <Flex alignItems={'center'} gap={15}>
          <DuelBotSVG />
          <Text
            textTransform="uppercase"
            fontSize={'20px'}
            fontWeight={600}
            color="white"
            letterSpacing={'0.18em'}
          >
            Unstake duelbots?
          </Text>
        </Flex>

        <Text
          fontSize="15px"
          color={'#B2D1FF'}
          maxWidth="530px"
          lineHeight={'18px'}
          mt="30px"
        >
          Are you sure you want to unstake your duelbot(s)? All unclaimed
          staking CHIP rewards will automatically be claimed and added to your
          wallet.
        </Text>

        <Flex gap={20} mt="30px" justifyContent={'flex-end'}>
          <Button
            background={'#242F42'}
            color="#768BAD"
            borderRadius="5px"
            fontSize={'14px'}
            fontWeight={600}
            p="12px 20px"
            onClick={hide}
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
            onClick={handleStaking}
          >
            Unstake Selected
          </Button>
        </Flex>
      </StyledBox>
    </Modal>
  );
}
