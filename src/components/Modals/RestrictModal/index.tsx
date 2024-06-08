import React, { FC } from 'react';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex, Text } from 'components';
import FlipClockCountdown from '@leenguyen/react-flip-clock-countdown';
import '@leenguyen/react-flip-clock-countdown/dist/index.css';
import styled from 'styled-components';
import { useAppSelector } from 'state';

const RestrictModal: FC<ModalProps> = ({ ...props }) => {
  const { selfExcludeTime } = useAppSelector(state => state.user);
  return (
    <Modal {...props}>
      <Container>
        <Text color="#E8F0FF" fontSize={'25px'} fontWeight={600}>
          Restricted Access
        </Text>

        <Text color={'#B9D2FD'} maxWidth="500px" mt="20px">
          You are currently restricted from all features on the Duel platform.
          You will regain access to Duel when the time is up.
        </Text>

        <Box mx="auto" mt="50px">
          <FlipClockCountdown
            to={new Date(new Date().getTime() + selfExcludeTime * 1000)}
            labels={['DAYS', 'HOURS', 'MINUTES', 'SECONDS']}
            labelStyle={{
              fontSize: 10,
              fontWeight: 500,
              textTransform: 'uppercase',
              fontFamily: 'Termina',
              paddingTop: '5px'
            }}
            digitBlockStyle={{
              width: 35,
              height: 55,
              fontSize: 40,
              fontWeight: 700,
              background: '#D9D9D9',
              color: '#333333'
            }}
            dividerStyle={{ color: 'black', height: 1 }}
            separatorStyle={{ color: 'white', size: '8px' }}
          />
        </Box>
      </Container>
    </Modal>
  );
};

const Container = styled(Flex)`
  flex-direction: column;
  flex: 1;
  background: linear-gradient(180deg, #132031 0%, #1a293c 100%);

  padding: 40px 21px;

  min-width: 350px;

  ${({ theme }) => theme.mediaQueries.md} {
    border: 2px solid #43546c;
    border-radius: 15px;

    padding: 40px 40px;
  }

  overflow: hidden auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
`;

export default RestrictModal;
