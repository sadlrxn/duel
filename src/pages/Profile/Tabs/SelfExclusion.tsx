import React, { useCallback, useState } from 'react';
import { Text } from 'components/Text';
import { Box, Button, Flex, Grid } from 'components';
import Checkbox from 'components/Checkbox';
import toast from 'utils/toast';
import { api } from 'services';
import { useAppDispatch } from 'state';
import { logoutUser, requestLogin } from 'state/user/actions';
import { useWallet } from '@solana/wallet-adapter-react';

export default function SelfExclusion() {
  const [date, setDate] = useState(0);
  const dispatch = useAppDispatch();
  const { disconnect } = useWallet();

  const handleChange = (e: any) => {
    setDate(Number(e.target.value));
  };

  const logout = useCallback(async () => {
    try {
      dispatch(requestLogin(false));
      await api.get('/user/logout');
    } catch (error) {
      console.info(error);
    } finally {
      disconnect();

      dispatch(logoutUser());
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleExclusion = async () => {
    if (date === 0) {
      toast.warn('check days!');
      return;
    }

    try {
      await api.post('/user/self-exclude', { days: date });
      logout();
    } catch (error) {
      toast.error('failed to self-exclude');
    }
  };

  return (
    <div className="container">
      <div className="box">
        <Text fontSize={'25px'} color="#fff" fontWeight={500}>
          Self Exclusion
        </Text>

        <Text color="#B9D2FD" fontSize="14px" mt="20px">
          For the majority of people, gambling is an enjoyable leisure and
          entertainment activity. But for some, gambling can have negative
          impacts. At Duel, we want to be an industry leader in providing a safe
          environment for our customers. We actively encourage and promote
          responsible gambling practices and provide tools to assist our
          customers in maintaining control of their gambling.
        </Text>

        <Box borderRadius={'10px'} p="20px" background="#121A25" mt="30px">
          <Text
            fontSize={'18px'}
            color="#D0DAEB"
            fontWeight={600}
            textTransform="uppercase"
            letterSpacing={'0.08em'}
          >
            Platform Access
          </Text>

          <Text color="#B9D2FD" fontSize="14px" mt="15px" maxWidth={'650px'}>
            This tool assists you in restricting access to all of Duel’s
            features for your selected time period. This action is irreversible,
            you will regain access once the time resets.
          </Text>

          <Flex
            mt="20px"
            alignItems={'center'}
            flexDirection={['column', 'column', 'column', 'row']}
            gap={10}
          >
            <Grid
              flex={1}
              gap={10}
              width={['100%', '100%', '100%', 'auto']}
              gridTemplateColumns={'repeat(auto-fit, minmax(120px, 1fr))'}
            >
              <Checkbox
                name="30"
                value={30}
                label="30 Days"
                checked={date === 30}
                onChange={handleChange}
              />
              <Checkbox
                name="60"
                value={60}
                label="60 Days"
                checked={date === 60}
                onChange={handleChange}
              />
              <Checkbox
                name="90"
                value={90}
                label="90 Days"
                checked={date === 90}
                onChange={handleChange}
              />
              <Checkbox
                name="365"
                value={365}
                label="365 Days"
                checked={date === 365}
                onChange={handleChange}
              />
            </Grid>

            <Button
              color="#4FFF8B"
              background="#1A5032"
              borderRadius={'5px'}
              p="10px 30px"
              width={['100%', '100%', '100%', 'auto']}
              onClick={handleExclusion}
            >
              Restrict Access
            </Button>
          </Flex>
        </Box>
      </div>
    </div>
  );
}
