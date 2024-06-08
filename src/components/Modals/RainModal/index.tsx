import React from 'react';
import styled, { css } from 'styled-components';

import { ReactComponent as RainIcon } from 'assets/imgs/icons/rain.svg';

import { ChatUser } from 'api/types/chat';
import Avatar from 'components/Avatar';
import { Modal, ModalProps } from 'components/Modal';
import { Flex, Grid } from 'components/Box';
import { Span } from 'components/Text';
import { Chip } from 'components/Chip';
import { Badge } from 'components/Badge';
import { formatUserName } from 'utils/format';
import { useAppSelector } from 'state';
import { convertBalanceToChip } from 'utils/balance';

export interface RainModalProps extends ModalProps {
  rainer?: ChatUser;
  duelers?: ChatUser[];
  amount?: number;
}

export default function RainModal({
  rainer = {
    id: 0,
    name: '',
    avatar: ''
  },
  amount = 0,
  duelers = [],
  ...props
}: RainModalProps) {
  const { id: userId } = useAppSelector(state => state.user);

  return (
    <StyledModal {...props}>
      <Container>
        <Flex gap={18} alignItems="center">
          <RainIcon width="22px" height="19px" />
          <Span fontWeight={600} fontSize="20px" color="white">
            {"IT'S RAINING CHIPS"}
          </Span>
        </Flex>

        <Flex
          flexWrap="wrap"
          fontWeight={400}
          alignItems="center"
          gap={5}
          width="100%"
          mt="26px"
        >
          <Span fontWeight={700}>{formatUserName(rainer.name)}</Span> has tipped{' '}
          {duelers.length} Duelers{' '}
          <Badge
            variant="secondary"
            fontSize={14}
            lineHeight="17px"
            px="5px"
            py="1px"
          >
            <Chip price={convertBalanceToChip(amount) / duelers.length} />
          </Badge>{' '}
          each, a total of{' '}
          <Badge
            variant="secondary"
            fontSize={14}
            lineHeight="17px"
            px="5px"
            py="1px"
          >
            <Chip price={convertBalanceToChip(amount)} />
          </Badge>
          <Span>.</Span> Make sure to give them a warm thank you in chat!
        </Flex>

        <UserContainer>
          {duelers.map(dueler => {
            return (
              <DuelerItem key={dueler.id} isUser={dueler.id === userId}>
                <Flex gap={12} alignItems="center">
                  <Avatar
                    userId={dueler.id}
                    name={dueler.name}
                    image={dueler.avatar}
                    border="none"
                    borderRadius="100%"
                    padding="0px"
                    size="28px"
                  />
                  <Span
                    fontWeight={700}
                    fontSize="14px"
                    color={dueler.id === userId ? 'success' : 'white'}
                  >
                    {formatUserName(dueler.name)}
                  </Span>
                </Flex>
                <Chip
                  price={convertBalanceToChip(amount) / duelers.length}
                  fontSize="14px"
                />
              </DuelerItem>
            );
          })}
        </UserContainer>
      </Container>
    </StyledModal>
  );
}

const DuelerItem = styled(Flex)<{ isUser: boolean }>`
  background: #182738;
  border-radius: 8px;
  align-items: center;
  justify-content: space-between;
  padding: 7px 13px 7px 15px;

  ${({ isUser }) =>
    isUser
      ? css`
          border: 1px solid #4fff8b;
          box-shadow: 0px 0px 10px rgba(79, 255, 139, 0.25);
        `
      : css``}
`;

const UserContainer = styled(Grid)`
  grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
  background: #0f1a26;
  border: 2px solid #0f1a26;
  border-radius: 13px;
  overflow: hidden auto;
  padding: 16px 15px 14px;
  margin-top: 30px;

  grid-row-gap: 10px;
  grid-column-gap: 15px;
  min-height: 76px;
`;

const Container = styled(Flex)`
  flex-direction: column;
  padding: 20px 30px 22px 25px;
  border-radius: 0px;
  background: linear-gradient(180deg, #132031 0%, #1a293c 100%);
  color: ${({ theme }) => theme.colors.textWhite};

  ${({ theme }) => theme.mediaQueries.md} {
    padding: 40px 50px 42px 45px;
    border-radius: 20px;
  }
  width: 100%;
  height: 100%;
`;

const StyledModal = styled(Modal)`
  padding: 1px;
  background: linear-gradient(
    180deg,
    #6a7f9e 0%,
    rgba(106, 127, 158, 0) 107.51%
  );

  width: 100vw;
  height: 100vh;
  overflow: hidden auto;

  border-radius: 0px;

  ${({ theme }) => theme.mediaQueries.md} {
    max-width: 620px;
    max-height: 490px;
    height: min-content;
    border-radius: 20px;
  }
`;
