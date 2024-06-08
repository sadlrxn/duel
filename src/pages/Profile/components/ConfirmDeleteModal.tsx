import React, { FC } from 'react';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex } from 'components/Box';
import { Text } from 'components/Text';
import { Button } from 'components/Button';
import Icon from './Icon';
import styled from 'styled-components';

interface ConfirmDeleteModalProps extends ModalProps {
  onDelete: () => void;
}
const ConfirmDeleteModal: FC<ConfirmDeleteModalProps> = ({
  onDelete,
  ...props
}) => {
  const handleDelete = async () => {
    onDelete();
    props.onDismiss!();
  };
  return (
    <Modal {...props}>
      <Container>
        <Flex gap={20} alignItems="center">
          <Icon />
          <Text
            fontSize={'20px'}
            fontWeight={600}
            letterSpacing="2px"
            color={'white'}
            textTransform="uppercase"
          >
            Delete Referral Code
          </Text>
        </Flex>

        <Box maxWidth={'500px'}>
          <Text color={'#768BAD'} maxWidth="720px" mt="30px">
            Are you sure you want to delete this referral code? The code will be
            deactivated for all users that are using this code. All unclaimed
            CHIP balances will automatically be claimed and added to your
            wallet.
          </Text>
        </Box>

        <Flex
          gap={20}
          flexDirection={['column', 'column', 'column', 'row']}
          justifyContent={'end'}
          mt="30px"
        >
          <Button
            p="12px 20px"
            background={'#242F42'}
            borderRadius="5px"
            fontSize={'16px'}
            fontWeight={600}
            color="#768BAD"
            onClick={props.onDismiss}
          >
            Cancel
          </Button>
          <Button
            p="12px 20px"
            background={'#501A1A'}
            borderRadius="5px"
            fontSize={'16px'}
            fontWeight={600}
            color="#FF5151"
            onClick={handleDelete}
          >
            Delete Code
          </Button>
        </Flex>
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
  width: 100vw;

  ${({ theme }) => theme.mediaQueries.md} {
    border: 2px solid #43546c;
    border-radius: 15px;
    width: 50vw;
    padding: 40px 40px;
  }

  overflow: hidden auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
`;

export default ConfirmDeleteModal;
