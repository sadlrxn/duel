import { useCallback } from 'react';
import styled from 'styled-components';
// import { Options } from "react-select";
import Skeleton from 'react-loading-skeleton';

import TipModal from '../TipModal';
import { Modal, ModalProps, useModal } from 'components/Modal';
import { Box, Flex, Grid } from 'components/Box';
import { Button } from 'components/Button';
import { Text, Span } from 'components/Text';
import Avatar from 'components/Avatar';
// import Select from "components/Select";

import { ReactComponent as DeathIcon } from 'assets/imgs/icons/death.svg';
import { ReactComponent as AngledSwordIcon } from 'assets/imgs/icons/angled-sword.svg';
import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { formatNumber, formatUserName } from 'utils/format';
import { useAppDispatch, useAppSelector } from 'state';
import { CopyIcon } from 'components/Icon';
import { ReactComponent as MuteIcon } from 'assets/imgs/icons/mute.svg';
import { ReactComponent as UnmuteIcon } from 'assets/imgs/icons/unmute.svg';
import { ReactComponent as BanIcon } from 'assets/imgs/icons/ban.svg';
import { ReactComponent as UnbanIcon } from 'assets/imgs/icons/handshake.svg';
import copy from 'copy-to-clipboard';
import { toast } from 'utils/toast';
import { sendMessage } from 'state/socket';
import useSWR, { useSWRConfig } from 'swr';
import { fetchUserInfo, FetchUserInfoResponse } from 'services';
import { AxiosError } from 'axios';
import { convertBalanceToChip } from 'utils/balance';

// const STATISTICS_OPTIONS: Options<{ label: string; value: string }> = [
//   { label: "All", value: "all" },
//   { label: "Coin Flip", value: "coinflip" },
//   { label: "Jackpot", value: "jackpot" },
// ];

interface ProfileModalProps extends ModalProps {
  userId?: number;
  avatar?: string;
  name?: string;
}

export default function ProfileModal({
  userId,
  avatar,
  name,
  ...props
}: ProfileModalProps) {
  const {
    id,
    name: userName,
    walletAddress,
    role: userRole
  } = useAppSelector(state => state.user);
  const { data, error } = useSWR<FetchUserInfoResponse, AxiosError>(
    `User Info: id:${userId} name:${name}`,
    () => fetchUserInfo({ userId, userName: name })
  );
  // const { data, error } = useUserInfo({ userId, userName: name });
  const { mutate } = useSWRConfig();

  const dispatch = useAppDispatch();

  const [onTipModal] = useModal(
    <TipModal userId={userId} name={name} avatar={avatar} />,
    true
  );

  const valueItem = useCallback(
    (val?: number, star: string = '****', showOperator: boolean = false) => {
      return (
        <>
          {error || !data ? (
            <Skeleton height={24} width={40} />
          ) : !data.statistics ? (
            star
          ) : (
            ((showOperator === true && val! > 0 && '+') || '') +
            formatNumber(val)
          )}
        </>
      );
    },
    [data, error]
  );

  if (data?.info.name.toLowerCase() === 'hidden' || error)
    return (
      <Modal {...props}>
        <Flex
          pt={['64px', '64px', '71px']}
          px={['33px', '33px', '42px']}
          pb="66px"
          background="linear-gradient(180deg, #132031 0%, #1A293C 100%)"
          borderRadius="20px"
          border="2px solid #43546c"
          maxWidth="530px"
          width="95vw"
        >
          <Text
            fontSize="24px"
            fontWeight={600}
            lineHeight="29px"
            color="white"
            textAlign="center"
          >
            {error
              ? 'User not found.'
              : 'This profile is hidden and is not publicly available.'}
          </Text>
        </Flex>
      </Modal>
    );

  const handleBan = () => {
    let message = `/${data?.info.banned ? 'un' : ''}ban ${data?.info.name}`;
    dispatch(
      sendMessage({
        type: 'message',
        room: 'chat',
        content: JSON.stringify({
          message: message.trim()
        })
      })
    );
    mutate(`User Info: id:${userId} name:${data?.info.name}`);
  };

  const handleMute = () => {
    let message = `/${data?.info.muted ? 'un' : ''}mute ${data?.info.name}`;
    dispatch(
      sendMessage({
        type: 'message',
        room: 'chat',
        content: JSON.stringify({
          message: message.trim()
        })
      })
    );
    mutate(`User Info: id:${userId} name:${data?.info.name}`);
  };

  return (
    <Modal {...props}>
      <Container>
        <Flex
          flexDirection={['column', 'column', 'column', 'row']}
          justifyContent={'space-between'}
          alignItems={'center'}
          gap={16}
        >
          <Flex alignItems="center" gap={20}>
            <Avatar image={data?.info.avatar} size="80px" border="none" />

            <Box>
              <Flex gap={10}>
                <Text
                  color={'#D0DAEB'}
                  fontSize="20px"
                  fontWeight={700}
                  lineHeight="28px"
                  letterSpacing={'0.08em'}
                >
                  {formatUserName(
                    userId === id ? userName : data?.info.name || ''
                  )}
                </Text>

                {id !== 0 &&
                  id !== userId &&
                  (userRole === 'admin' || userRole === 'moderator') && (
                    <>
                      <Button
                        borderRadius={'5px'}
                        background="#304869"
                        size={30}
                        onClick={handleBan}
                      >
                        {data?.info.banned ? (
                          <UnbanIcon color="#768BAD" />
                        ) : (
                          <BanIcon color="#768BAD" />
                        )}
                      </Button>

                      <Button
                        borderRadius={'5px'}
                        background="#304869"
                        size={30}
                        onClick={handleMute}
                      >
                        {data?.info.muted ? (
                          <UnmuteIcon color="#768BAD" />
                        ) : (
                          <MuteIcon color="#768BAD" />
                        )}
                      </Button>
                    </>
                  )}
              </Flex>

              {(userRole === 'admin' || userRole === 'moderator') && (
                <Flex alignItems={'center'} gap={10}>
                  <Text
                    color={'#556273'}
                    fontSize="13px"
                    fontWeight={600}
                    lineHeight="28px"
                  >
                    {formatUserName(
                      userId === id
                        ? walletAddress
                        : data?.info.walletAddress || ''
                    )}
                  </Text>
                  <Button
                    color="#768BAD"
                    borderRadius={'5px'}
                    background="#242F42"
                    p="5px"
                    fontWeight={600}
                    onClick={() => {
                      copy(
                        userId === id
                          ? walletAddress
                          : data?.info.walletAddress || ''
                      );
                      toast.success('Copied!');
                    }}
                  >
                    <CopyIcon />
                  </Button>
                </Flex>
              )}
            </Box>
          </Flex>
          {id !== 0 && id !== userId && (
            <TipButton px="20px" py="5px" onClick={onTipModal}>
              <CoinIcon />
              Send Tip
            </TipButton>
          )}
        </Flex>

        <Text color="white" fontSize={'30px'} fontWeight={500} mt="30px">
          Statistics
        </Text>

        <Grid
          background={'#121A25BF'}
          p="20px"
          gridTemplateColumns={'repeat(auto-fit, minmax(120px, 1fr))'}
          borderRadius="10px"
          mt="20px"
        >
          <Box>
            <Text color={'#96A8C2'} fontWeight={600} letterSpacing={'2px'}>
              TOTAL GAMES
            </Text>
            <Text
              color={'white'}
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {valueItem(data?.statistics?.total_rounds ?? 0)}
            </Text>
          </Box>

          <Box>
            <Text color={'#96A8C2'} fontWeight={600} letterSpacing={'2px'}>
              GAMES WON
            </Text>
            <Text
              color="white"
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {valueItem(data?.statistics?.winned_rounds ?? 0)}
            </Text>
          </Box>

          <Box>
            <Text color="#96A8C2" fontWeight={600} letterSpacing={'2px'}>
              GAMES LOST
            </Text>
            <Text
              color="white"
              fontSize={['1.66em', '1.5em', '1.5em']}
              fontWeight={700}
            >
              {valueItem(data?.statistics?.lost_rounds ?? 0)}
            </Text>
          </Box>

          <Box>
            <Text color={'#96A8C2'} fontWeight={600} letterSpacing={'2px'}>
              WIN RATIO
            </Text>
            <Text
              color="success"
              fontWeight={700}
              fontSize={['1.66em', '1.5em', '1.5em']}
              display="flex"
            >
              {valueItem(
                data
                  ? (data.statistics?.total_rounds ?? 0) === 0
                    ? 0
                    : +(
                        ((data.statistics?.winned_rounds ?? 0) /
                          (data.statistics?.total_rounds ?? 100)) *
                        100
                      ).toFixed(2)
                  : 0,
                '**'
              )}
              %
            </Text>
          </Box>
        </Grid>

        <Divider />

        <Grid
          gap={20}
          my="20px"
          gridTemplateColumns={'repeat(auto-fit, minmax(250px, 1fr))'}
          color="#96A8C2"
        >
          <Box background={'#121A25BF'} p="20px" borderRadius={'10px'}>
            <Text color={'#96A8C2'} fontWeight={600} letterSpacing={'2px'}>
              BEST & WORST STREAKS
            </Text>
            <Flex alignItems={'center'}>
              <AngledSwordIcon width={'19px'} />
              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {valueItem(data?.statistics?.best_streaks ?? 0, '*')}
              </Span>
              <Span
                color="#4F617B"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={500}
                mx="15px"
              >
                /
              </Span>

              <DeathIcon />
              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {valueItem(data?.statistics?.worst_streaks ?? 0, '*')}
              </Span>
            </Flex>
          </Box>
          <Box background={'#121A25BF'} p="20px" borderRadius={'10px'}>
            <Text color={'#96A8C2'} fontWeight={600} letterSpacing={'2px'}>
              TOTAL WAGERED
            </Text>

            <Flex alignItems={'center'}>
              <CoinIcon />

              <Span
                color="white"
                fontSize={['1.66em', '1.5em', '1.5em']}
                fontWeight={700}
                ml="10px"
              >
                {valueItem(
                  data
                    ? convertBalanceToChip(data.statistics?.total_wagered ?? 0)
                    : 0
                )}
              </Span>
            </Flex>
          </Box>
        </Grid>

        {/* <Grid
          gap={20}
          gridTemplateColumns={"repeat(auto-fit, minmax(250px, 1fr))"}
          color="#96A8C2"
        >
          <Detail2>
            <Text fontWeight={600} letterSpacing={"2px"}>
              MAX PROFIT
            </Text>
            <Flex alignItems={"center"}>
              <CoinIcon />

              <Span
                color="white"
                fontSize={["1.66em", "1.5em", "1.5em"]}
                fontWeight={700}
                ml="10px"
              >
                {valueItem(data ? (data.statistics?.max_profit ?? 0) / 100 : 0)}
              </Span>
            </Flex>
          </Detail2>
          <Detail2>
            <Text fontWeight={600} letterSpacing={"2px"}>
              TOTAL PROFIT
            </Text>

            <Flex alignItems={"center"}>
              <CoinIcon />

              <Span
                color={
                  data && (data!.statistics?.total_profit ?? 0) < 0
                    ? "warning"
                    : "success"
                }
                fontSize={["1.66em", "1.5em", "1.5em"]}
                fontWeight={700}
                ml="10px"
              >
                {valueItem(
                  data ? (data.statistics?.total_profit ?? 0) / 100 : 0,
                  "****",
                  true
                )}
              </Span>
            </Flex>
          </Detail2>
        </Grid> */}
      </Container>
    </Modal>
  );
}

const Divider = styled(Box)`
  height: 2px;
  margin-top: 20px;
  background-color: #354d6d;

  ${({ theme }) => theme.mediaQueries.md} {
    display: none;
  }
`;

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

const TipButton = styled(Button)`
  display: flex;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.chip};
  background: #242f42;
  padding: 9px 10px;
`;

TipButton.defaultProps = { variant: 'secondary', type: 'button' };
