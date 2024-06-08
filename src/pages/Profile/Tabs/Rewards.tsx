import React, { useCallback, useMemo } from 'react';
import { useSWRConfig } from 'swr';
import { Text } from 'components/Text';
import { Box, Button, Flex } from 'components';
import styled from 'styled-components';
import { CoinIcon } from 'components';
import { api } from 'services';
import { useCustomSWR } from 'hooks';
import { toast } from 'react-toastify';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';
import useSWR from 'swr';
import { useAppDispatch, useAppSelector } from 'state';
import { updateBalance } from 'state/user/actions';

export default function Rewards() {
  const dispatch = useAppDispatch();
  const { config } = useAppSelector(state => state.user);

  const { mutate } = useSWRConfig();

  const { data: dailyRewards } = useSWR(`/daily-race/rewards`, async arg =>
    api.get(arg).then(res => res.data)
  );

  const { data: weeklyRewards } = useSWR('/weekly-raffle/rewards', async arg =>
    api.get(arg).then(res => res.data)
  );

  const { data } = useCustomSWR({
    key: 'rewards',
    route: '/rewards',
    method: 'get'
  });

  const dailyWagered = useMemo(() => {
    if (!dailyRewards) return 0;
    return dailyRewards.wagered;
  }, [dailyRewards]);

  const rewards = useMemo(() => {
    let rewards = { rakeback: 0, rakebackTotal: 0 };
    if (data) {
      rewards.rakeback = data.rewards.rakeback;
      rewards.rakebackTotal = data.rewards['rakeback-total'];
    }
    return rewards;
  }, [data]);

  const totalDailyRewards = useMemo(() => {
    let sum = 0;

    if (!dailyRewards) return 0;

    dailyRewards.rewards.forEach((item: any) => {
      sum += item.prize;
    });

    return sum;
  }, [dailyRewards]);

  const totalWeeklyRewards = useMemo(() => {
    let sum = 0;

    if (!weeklyRewards) return 0;

    weeklyRewards.rewards.forEach((item: any) => {
      sum += item.prize;
    });

    return sum;
  }, [weeklyRewards]);

  const handleClaimRakeback = useCallback(async () => {
    await api
      .post('/rewards/rakeback')
      .then(res => {
        dispatch(updateBalance({ type: 1, usdAmount: res.data }));
      })
      .catch((error: any) => {
        if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        }
      })
      .finally(() => {
        mutate('rewards');
      });
  }, [mutate, dispatch]);

  const handleClaimDailyRewards = async () => {
    if (!dailyRewards) return;

    await api
      .post('/daily-race/claim', {
        ids: dailyRewards.rewards.map((item: any) => item.id)
      })
      .then(res => {
        dispatch(updateBalance({ type: 1, usdAmount: res.data }));
      })
      .catch((error: any) => {
        if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        }
      })
      .finally(() => {
        mutate('/daily-race/rewards');
      });
  };

  const handleClaimWeeklyRewards = async () => {
    if (!weeklyRewards) return;

    await api
      .post('/weekly-raffle/claim', {
        ids: weeklyRewards.rewards.map((item: any) => item.id)
      })
      .then(res => {
        dispatch(updateBalance({ type: 1, usdAmount: res.data }));
      })
      .catch((error: any) => {
        if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        }
      })
      .finally(() => {
        mutate('/weekly-raffle/rewards');
      });
  };

  return (
    <div className="container">
      <div className="box">
        <Text color="white" fontSize={'25px'} mb="10px" fontWeight={500}>
          Rewards
        </Text>

        <StyledFlex>
          <Box width={['auto', 'auto', 'auto', '500px']}>
            <Text fontSize={'18px'} color="#D0DAEB" fontWeight={600}>
              RAKEBACK
            </Text>

            <Text fontSize={'14px'} color="#B9D2FD">
              Get {config.rakebackRate}% house edge back from every bet against
              the house.
            </Text>
          </Box>

          <Box>
            <Text
              fontSize={'14px'}
              color="#D0DAEB"
              fontWeight={600}
              letterSpacing="0.08em"
            >
              RAKEBACK
            </Text>

            <Text
              fontSize={'19px'}
              color="#D0DAEB"
              textAlign={['start', 'start', 'start', 'center']}
            >
              +{config.rakebackRate}%
            </Text>
          </Box>

          <StyledButton
            width={'150px'}
            onClick={handleClaimRakeback}
            disabled={rewards.rakeback < convertChipToBalance(0.01)}
          >
            Claim
            <CoinIcon />
            {convertBalanceToChip(rewards.rakeback)}
          </StyledButton>
        </StyledFlex>

        <StyledFlex>
          <Box width={['auto', 'auto', 'auto', '500px']}>
            <Text fontSize={'18px'} color="#D0DAEB" fontWeight={600}>
              DAILY RACE
            </Text>

            <Text fontSize={'14px'} color="#B9D2FD">
              Duel rewards players with the the highest daily wagers with CHIPS.
            </Text>
          </Box>

          <Box>
            <Text
              fontSize={'14px'}
              color="#D0DAEB"
              fontWeight={600}
              letterSpacing="0.08em"
            >
              WAGERED
            </Text>

            <Flex alignItems="center" gap={10} mt="5px">
              <CoinIcon />
              <Text color="#D0DAEB" fontWeight={600}>
                {convertBalanceToChip(dailyWagered)}
              </Text>
            </Flex>
          </Box>

          <StyledButton
            width={'150px'}
            onClick={handleClaimDailyRewards}
            disabled={totalDailyRewards < 10 ? true : false}
          >
            Claim
            <CoinIcon />
            {convertBalanceToChip(totalDailyRewards)}
          </StyledButton>
        </StyledFlex>

        <StyledFlex>
          <Box width={['auto', 'auto', 'auto', '500px']}>
            <Text fontSize={'18px'} color="#D0DAEB" fontWeight={600}>
              WEEKLY RAFFLE
            </Text>

            <Text fontSize={'14px'} color="#B9D2FD">
              Wager to earn raffle tickets for a chance to with the Grand Prize!
            </Text>
          </Box>

          <Box>
            <Text
              fontSize={'14px'}
              color="#D0DAEB"
              fontWeight={600}
              letterSpacing="0.08em"
            >
              MY TICKETS
            </Text>

            <Text
              color="#D0DAEB"
              fontWeight={600}
              mt="5px"
              textAlign={'center'}
            >
              {weeklyRewards ? weeklyRewards.tickets : '0'}
            </Text>
          </Box>

          <StyledButton
            width={'150px'}
            onClick={handleClaimWeeklyRewards}
            disabled={totalWeeklyRewards < 10 ? true : false}
          >
            Claim
            <CoinIcon />
            {convertBalanceToChip(totalWeeklyRewards)}
          </StyledButton>
        </StyledFlex>
      </div>
    </div>
  );
}

const StyledFlex = styled(Flex)`
  flex-direction: column;
  gap: 15px;
  background: #121a25;
  align-items: normal;

  justify-content: space-between;
  padding: 20px;
  border-radius: 10px;
  margin-bottom: 20px;

  .width_800 & {
    flex-direction: row;
    align-items: center;
  }
`;

const StyledButton = styled(Button)`
  background: #1a5032;
  font-size: 14px;
  font-weight: 600;
  color: #4fff8b;
  border-radius: 5px;
  padding: 10px 0px;
`;

// const StyledTable = styled.table`
//   width: 100%;
//   font-family: 'Inter';
//   font-size: 14px;
//   color: #b2d1ff;
//   border-collapse: separate;
//   border-spacing: 0px 10px;

//   th {
//     font-family: 'Inter';
//     font-style: normal;
//     font-weight: 500;
//     font-size: 14px;
//     line-height: 17px;

//     letter-spacing: 0.17em;
//     text-transform: uppercase;
//     color: #b2d1ff;
//     text-align: left;
//     padding: 8px 20px;
//   }

//   tbody tr {
//     background-color: #182738;
//     cursor: pointer;
//     &:hover {
//       border: 1px solid #49f884;
//     }

//     &:hover td:first-child {
//       border-top-left-radius: 8px;
//       border-bottom-left-radius: 8px;
//     }

//     &:hover td:last-child {
//       border-top-right-radius: 8px;
//       border-bottom-right-radius: 8px;
//     }
//   }

//   td {
//     padding: 8px 20px;
//   }

//   tr td:first-child {
//     border-top-left-radius: 8px;
//     border-bottom-left-radius: 8px;
//   }
//   tr td:last-child {
//     border-top-right-radius: 8px;
//     border-bottom-right-radius: 8px;
//   }
// `;
