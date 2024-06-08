import React, { FC } from 'react';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex } from 'components/Box';
import { Span, Text } from 'components/Text';
import { Button } from 'components/Button';
import Icon from './Icon';
import styled from 'styled-components';
import { Avatar, Chip } from 'components';
import useSWR from 'swr';
import { api } from 'services';
import { formatNumber, formatSecond2Day, formatUserName } from 'utils/format';
import { ClipLoader } from 'react-spinners';
import { convertBalanceToChip } from 'utils/balance';

const StyledTable = styled.table`
  width: 100%;
  font-family: 'Inter';
  font-size: 14px;
  color: #b2d1ff;
  border-collapse: separate;
  border-spacing: 0px 10px;

  th {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 17px;

    letter-spacing: 0.17em;
    text-transform: uppercase;
    color: #b2d1ff;
    text-align: left;
    padding: 8px 10px;
  }

  tbody tr {
    background-color: #182738;
    cursor: pointer;
    &:hover {
      border: 1px solid #49f884;
    }

    &:hover td:first-child {
      border-top-left-radius: 8px;
      border-bottom-left-radius: 8px;
    }

    &:hover td:last-child {
      border-top-right-radius: 8px;
      border-bottom-right-radius: 8px;
    }
  }

  td {
    padding: 8px 10px;
  }

  tbody tr {
    &:hover {
      background: #263449;
    }
  }
  tr td:first-child {
    border-top-left-radius: 8px;
    border-bottom-left-radius: 8px;
  }
  tr td:last-child {
    border-top-right-radius: 8px;
    border-bottom-right-radius: 8px;
  }
`;

interface AffiliateUserModalProps extends ModalProps {
  code: string;
}
const AffiliateUserModal: FC<AffiliateUserModalProps> = ({
  code,
  ...props
}) => {
  const { data: info } = useSWR(`/affiliate/code-detail`, async arg =>
    api.get(`${arg}?code=${code}`).then(res => res.data)
  );

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
            Affiliate code users
          </Text>
        </Flex>
        {info ? (
          <Box
            background={'#0F1A26'}
            borderRadius="13px"
            p="10px"
            mt="25px"
            overflowX={'auto'}
          >
            <StyledTable>
              <thead>
                <tr>
                  <th>Users</th>
                  <th>Days Active</th>
                  <th>Wagered</th>
                </tr>
              </thead>
              <tbody>
                {info.users.map((item: any) => (
                  <tr key={item.name}>
                    <td>
                      <Flex gap={10} alignItems="center">
                        <Avatar
                          userId={item.id}
                          image={item.avatar}
                          size="28px"
                          border="none"
                        />
                        <Text color="white" fontWeight={500}>
                          {formatUserName(item.name)}
                        </Text>
                      </Flex>
                    </td>
                    <td>{formatSecond2Day(item.lifetime)}</td>
                    <td>
                      <Chip price={convertBalanceToChip(item.wagered)} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </StyledTable>
          </Box>
        ) : (
          <Box mx={'auto'} mt="30px">
            <ClipLoader color="#fff" size={40} />
          </Box>
        )}
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

export default AffiliateUserModal;
