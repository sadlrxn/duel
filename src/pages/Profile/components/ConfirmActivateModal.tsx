import React, { FC } from 'react';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex } from 'components/Box';
import { Span, Text } from 'components/Text';
import { Button } from 'components/Button';
import Icon from './Icon';
import styled from 'styled-components';

interface ConfirmActivateModalProps extends ModalProps {
  code: string;
  onActivate: () => Promise<void>;
}
const ConfirmActivateModal: FC<ConfirmActivateModalProps> = ({
  code,
  onActivate,
  ...props
}) => {
  const handleActivate = async () => {
    await onActivate();
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
            Referral Code Activation
          </Text>
        </Flex>

        <Box maxWidth={'500px'}>
          <Text color={'#768BAD'} maxWidth="720px" mt="30px">
            You are about to activate referral code&nbsp;
            <Span fontWeight={600}>{code}</Span>. When activated you will
            receive a 5% increase your Rakeback for 24 hours.
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
            background={'#1A5032'}
            borderRadius="5px"
            fontSize={'16px'}
            fontWeight={600}
            color="#4FFF8B"
            onClick={handleActivate}
          >
            Activate Code
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

export default ConfirmActivateModal;
