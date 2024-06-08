import React, { useMemo, useState } from 'react';
import styled from 'styled-components';
import useSWR, { useSWRConfig } from 'swr';

import { Flex, Box, Grid, Text, Button, Avatar, useModal } from 'components';
import Coin from 'components/Icon/Coin';
import { InputBox } from 'components/InputBox';

import { useAppDispatch } from 'state';
import { updateBalance } from 'state/user/actions';

import api from 'utils/api';
import toast from 'utils/toast';
import { convertBalanceToChip } from 'utils/balance';
import { formatNumber, formatUserName } from 'utils/format';

import ReferralRow from '../components/ReferralRow';
import ConfirmActivateModal from '../components/ConfirmActivateModal';

export interface ReferralData {
  activeAffiliate: {
    code: string;
    ownerAvatar: string;
    ownerName: string;
    ownerId: number;
    rate: number;
  };
  created: {
    code: string;
    reward: number;
    totalEarned: number;
    totalWagered: number;
    userCnt: number;
    rate: number;
  }[];
}

const StyledFlex = styled(Flex)`
  flex-direction: column;
  gap: 15px;
  align-items: center;
  div,
  button {
    width: 100%;
  }
  ${({ theme }) => theme.mediaQueries.md} {
    flex-direction: row;
    gap: 30px;

    div,
    button {
      width: auto;
    }
  }
`;

export default function Referral() {
  const dispatch = useAppDispatch();

  const { mutate } = useSWRConfig();
  const [formData, setFormData] = useState({
    forActive: '',
    forCreate: ''
  });
  const { data: referralCodes } = useSWR(`/affiliate/my-codes`, async arg =>
    api.get<ReferralData>(arg).then(res => res.data)
  );

  const info = useMemo(() => {
    let totalWagered = 0;
    let totalClaimed = 0;
    let availableToClaim = 0;

    if (referralCodes)
      referralCodes.created.forEach((item: any) => {
        totalWagered += item.totalWagered;
        totalClaimed += item.totalEarned;
        availableToClaim += item.reward;
      });

    totalWagered = convertBalanceToChip(totalWagered);
    totalClaimed = convertBalanceToChip(totalClaimed);
    availableToClaim = convertBalanceToChip(availableToClaim);

    return { totalWagered, totalClaimed, availableToClaim };
  }, [referralCodes]);

  const handleInputChange = (e: any) => {
    const input = e.target.value;
    if (e.target.name === 'forCreate') {
      setFormData({ ...formData, forCreate: input });
    } else {
      setFormData({ ...formData, forActive: input });
    }
  };

  const createReferralCode = (e: any) => {
    e.preventDefault();
    if (formData.forCreate === '') {
      toast.warn('input code!');
      return;
    }

    api
      .post('/affiliate/create', { codes: [formData.forCreate] })
      .then(() => {
        toast.success(
          'Referral code created! Your code is now ready to be shared.'
        );
      })
      .catch(err => {
        if (err.response.status === 406) {
          if (err.response.data.errorCode === 13011)
            toast.error('This referral code already exists.');
          else if (err.response.data.errorCode === 13012)
            toast.error('Failed! There was an error creating your code.');
          else if (err.response.data.errorCode === 13013)
            toast.error('Exceed the referral code limit.');
          else if (err.response.data.errorCode === 13014)
            toast.error('Should wager more than 2500 chips to create code.');
          else if (err.response.data.errorCode === 13015)
            toast.error('Code should not contain space');
          else if (err.response.data.errorCode === 13016)
            toast.error('Code should not contain reserved word');
        } else if (err.response.status === 429) {
          toast.error(err.response.data.message);
        } else if (err.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error('Failed! There was an error creating your code.');
      })
      .finally(() => {
        mutate(`/affiliate/my-codes`);
        setFormData({ ...formData, forCreate: '' });
      });
  };

  const deleteReferralCode = (code: string) => {
    api
      .post('/affiliate/delete', { codes: [code] })
      .then(res => {
        dispatch(updateBalance({ type: 1, usdAmount: res.data.claimed }));
        toast.success('Referral code Deleted!');
      })
      .catch(err => {
        if (err.response.status === 429) toast.error(err.response.data.message);
        else if (err.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error('Failed! There was an error deleting your code.');
      })
      .finally(() => {
        mutate(`/affiliate/my-codes`);
      });
  };

  const activateReferralCode = async () => {
    if (formData.forActive === '') {
      toast.warn('input code!');
      return;
    }
    try {
      await api.post('/affiliate/activate', { code: formData.forActive });
      toast.success('Referral code updated!');
    } catch (err: any) {
      if (err.response.status === 406) {
        if (err.response.data.errorCode === 13031)
          toast.error('You can’t activate your own code.');
        else if (err.response.data.errorCode === 13032)
          toast.error('This referral code doesn’t exist.');
        else if (err.response.data.errorCode === 13034)
          toast.error('You can’t activate code after 24 hours from sign up');
        else if (err.response.data.errorCode === 13033)
          toast.error('Failed! There was an error activating referral code.');
      } else if (err.response.status === 429) {
        toast.error(err.response.data.message);
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else
        toast.error('Failed! There was an error activating referral code.');
    } finally {
      setFormData({ ...formData, forActive: '' });
      mutate(`/affiliate/my-codes`);
    }
  };

  const claimReferralCode = (codes: string[]) => {
    api
      .post('/affiliate/claim', { codes })
      .then(res => {
        dispatch(updateBalance({ type: 1, usdAmount: res.data.claimed }));

        toast.success('Successfully claimed!');
      })
      .catch(err => {
        if (err.response.status === 429) toast.error(err.response.data.message);
        else
          toast.error('Failed! There was an error claiming affiliate rewards.');
      })
      .finally(() => {
        mutate(`/affiliate/my-codes`);
      });
  };

  const [onConfirmActivate] = useModal(
    <ConfirmActivateModal
      code={formData.forActive}
      onActivate={activateReferralCode}
    />
  );

  return (
    <div className="container">
      <div className="box">
        <Text fontSize={'25px'} color="#fff" fontWeight={500}>
          Redeem Code
        </Text>

        <Flex
          justifyContent={'space-between'}
          alignItems="center"
          background={'#121A25'}
          p="20px"
          mt="30px"
          gap={10}
          flexDirection={['column', 'column', 'column', 'row']}
          borderRadius={'10px'}
        >
          <form
            onSubmit={(e: any) => {
              e.preventDefault();
              if (formData.forActive === '') {
                toast.warn('input code!');
                return;
              }
              onConfirmActivate();
            }}
          >
            <Box>
              <Text
                textTransform="uppercase"
                fontSize={'18px'}
                fontWeight={800}
                color="#D0DAEB"
                letterSpacing={'0.08em'}
              >
                Activate code
              </Text>
              <Text fontSize={'14px'} color="#B9D2FD" mt="10px">
                Active someones code and get 5% increase in Rakeback for 24
                hours.
              </Text>
            </Box>

            <Text fontSize={'16px'} color="#768BAD" mt="25px">
              Referral Code
            </Text>

            <StyledFlex mt={'8px'}>
              <InputBox
                gap={20}
                padding={'15px 15px !important'}
                background="#090e14 !important"
                maxWidth={'350px'}
              >
                <input
                  placeholder="DUEL777"
                  name="forActive"
                  value={formData.forActive}
                  onChange={handleInputChange}
                  style={{ fontSize: '16px' }}
                />
              </InputBox>

              <Button
                color="#4FFF8B"
                borderRadius={'5px'}
                background="#1A5032"
                p="15px 30px"
                type="submit"
              >
                Activate Code
              </Button>
            </StyledFlex>
          </form>

          {referralCodes && referralCodes.activeAffiliate && (
            <Box background={'#090E14'} borderRadius="8px" p="15px">
              <Text
                textAlign={'center'}
                fontSize={'16px'}
                color="#D0DAEB"
                letterSpacing={'0.08em'}
              >
                Active Code
              </Text>
              <Text
                textAlign={'center'}
                fontSize={'18px'}
                fontWeight={600}
                color="#B9D2FD"
                mt="10px"
              >
                {referralCodes.activeAffiliate.code}
              </Text>

              <Flex
                alignItems={'center'}
                justifyContent="center"
                gap={10}
                mt="9px"
              >
                <Avatar
                  userId={referralCodes.activeAffiliate.ownerId}
                  image={referralCodes.activeAffiliate.ownerAvatar}
                  size="30px"
                  padding="0px"
                  border="none"
                />
                <Text
                  textAlign={'center'}
                  fontSize={'16px'}
                  fontWeight={500}
                  color="#B9D2FD"
                >
                  {formatUserName(referralCodes.activeAffiliate.ownerName)}
                </Text>
              </Flex>
            </Box>
          )}
        </Flex>

        <Text fontSize={'25px'} color="#fff" mt="30px" fontWeight={500}>
          Earn with Referrals
        </Text>

        <Box background={'#121A25'} p="20px" mt="30px" borderRadius={'10px'}>
          <form onSubmit={createReferralCode}>
            <Text
              textTransform="uppercase"
              fontSize={'18px'}
              fontWeight={800}
              color="#D0DAEB"
              letterSpacing={'0.08em'}
            >
              Create code
            </Text>

            <Text
              fontSize={'14px'}
              mt="10px"
              color="#B9D2FD"
              maxWidth={'500px'}
            >
              Earn 5% of house edge from Duelers that have your code active.
              Duelers using your code will receive a 5% increase in Rakeback for
              24 hours.
            </Text>

            <Text fontSize={'16px'} color="#768BAD" mt="25px">
              Your Referral Code
            </Text>

            <StyledFlex mt={'8px'}>
              <InputBox
                gap={20}
                padding={'15px 15px !important'}
                background="#090e14 !important"
                maxWidth={'350px'}
              >
                <input
                  placeholder="DUEL777"
                  name="forCreate"
                  value={formData.forCreate}
                  style={{ fontSize: '16px' }}
                  onChange={handleInputChange}
                />
              </InputBox>

              <Button
                color="#4FFF8B"
                borderRadius={'5px'}
                background="#1A5032"
                p="15px 30px"
                type="submit"
              >
                Create Code
              </Button>
            </StyledFlex>
          </form>
        </Box>

        <Container>
          <Box background={'#121A25'} borderRadius="13px" p={'19px 30px'}>
            <Text
              textTransform="uppercase"
              color={'#96A8C2'}
              fontSize="12px"
              fontWeight={600}
              letterSpacing="0.185em"
            >
              total wagered
            </Text>

            <Flex alignItems={'center'} gap={5} mt="10px">
              <Coin />
              <Text
                textTransform="uppercase"
                color={'#fff'}
                fontSize="20px"
                fontWeight={700}
              >
                {formatNumber(info.totalWagered)}
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
              total Earned
            </Text>
            <Flex alignItems={'center'} gap={5} mt="10px">
              <Coin />
              <Text
                textTransform="uppercase"
                color={'#fff'}
                fontSize="20px"
                fontWeight={700}
              >
                {formatNumber(info.totalClaimed)}
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

            <Flex justifyContent={'space-between'} alignItems="end">
              <Flex alignItems={'center'} gap={5} mt="10px">
                <Coin />
                <Text
                  textTransform="uppercase"
                  color={'#fff'}
                  fontSize="20px"
                  fontWeight={700}
                >
                  {formatNumber(info.availableToClaim)}
                </Text>
              </Flex>

              <Button
                color="#4FFF8B"
                borderRadius={'5px'}
                background="#1A5032"
                p="8px 20px"
                fontWeight={600}
                disabled={info.availableToClaim < 0.1 ? true : false}
                onClick={() => {
                  if (!referralCodes) return;
                  const codes = referralCodes.created.map(item => item.code);
                  claimReferralCode(codes);
                }}
              >
                Claim
              </Button>
            </Flex>
          </Box>
        </Container>

        {referralCodes && referralCodes.created.length > 0 && (
          <Box mt="30px">
            <Box
              background="#121A25"
              p="20px"
              borderRadius="13px"
              overflowX={'auto'}
            >
              <Text
                textTransform="uppercase"
                fontSize="18px"
                fontWeight={600}
                color="#D0DAEB"
                mb="5px"
              >
                Your referral codes
              </Text>
              <Text color="#B9D2FD" mb="15px" mt="5px">
                To collect referral earnings the code must have at least 2,500
                CHIPS wagered.
              </Text>
              <StyledTable>
                <thead>
                  <tr>
                    <th>Code</th>
                    <th>Users</th>
                    <th>Rate</th>
                    <th>Wagered</th>
                    <th>Earned </th>
                    <th>Claim</th>
                    <th>Share link</th>
                  </tr>
                </thead>
                <tbody>
                  {referralCodes?.created.map(item => (
                    <ReferralRow
                      key={item.code}
                      data={item}
                      onClaim={claimReferralCode}
                      onDelete={deleteReferralCode}
                    />
                  ))}
                </tbody>
              </StyledTable>

              {/* <Button
              variant="secondary"
              outlined
              scale="sm"
              width={153}
              background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
              color="#FFFFFF"
              borderColor="chipSecondary"
              marginX="auto"
              marginTop={10}
            >
              SHOW MORE
            </Button> */}
            </Box>
          </Box>
        )}
      </div>
    </div>
  );
}

const Container = styled(Grid)`
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  align-items: end;
  margin-top: 22px;
  gap: 22px;
`;

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
    padding: 8px 20px;

    :last-child {
      min-width: 210px;
      width: 210px;
    }
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
    padding: 8px 20px;
    :last-child {
      min-width: 210px;
      width: 210px;
    }
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
