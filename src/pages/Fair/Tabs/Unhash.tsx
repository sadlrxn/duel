import React, {
  useState,
  useRef,
  useEffect,
  useMemo,
  useCallback
} from 'react';
import styled from 'styled-components';

import { Text, Flex, Button, Input } from 'components';
import { ReactComponent as CopyIcon } from 'assets/imgs/icons/copy.svg';
import {
  CopyButton,
  InputBox
} from 'pages/Jackpot/components/Modal/Fairness/styles';
import copy from 'copy-to-clipboard';
import { toast } from 'react-toastify';
import { api } from 'services';

interface UnhashProps {
  seedHash?: string;
}

const Unhash: React.FC<UnhashProps> = ({ seedHash }) => {
  const [seed, setSeed] = useState('');
  const [hash, setHash] = useState('');

  const handleClick = useCallback(async () => {
    try {
      const { data } = await api.get(`/seed/unhash?hash=${hash}`);
      setSeed(data.seed);
    } catch (error) {
      toast.error(
        'The current server seed is still in use, please change your client seed to unhash the server seed.'
      );
      return;
    }
  }, [hash]);
  return (
    <Flex flexDirection="column" gap={25}>
      <Flex flexDirection="column" gap={10}>
        <Text fontSize="18px" lineHeight="22px" fontWeight={500} color="white">
          UNHASH SERVER SEED
        </Text>
        <Text fontSize="14px" lineHeight="20px" color="#B9D2FD">
          Everytime you generate a new client seed a new hashed server seed is
          paired with it. You can only unhash a server seed once you generate a
          new client seed.
        </Text>
      </Flex>
      <Flex flexDirection="column" gap={8}>
        <Text fontSize="16px" lineHeight="19px" color="#768BAD">
          Server Seed (Hashed)
        </Text>
        <Flex flexDirection="row" justifyContent="space-between" gap={22}>
          <StyledInput
            value={hash}
            onChange={e => {
              setHash(e.target.value);
            }}
          ></StyledInput>
          <UnhashButton onClick={handleClick}>UNHASH</UnhashButton>
        </Flex>
      </Flex>
      <Flex flexDirection="column" gap={8}>
        <Text fontSize="16px" lineHeight="19px" color="#768BAD">
          Server Seed
        </Text>
        <InputBox readOnly={true}>
          <Input readOnly={true} value={seed} type="text" />
          <CopyButton
            onClick={() => {
              copy(seed);
              toast.success('Copy to clipboard success.');
            }}
          >
            <CopyIcon width={12} height={12} />
          </CopyButton>
        </InputBox>
      </Flex>
    </Flex>
  );
};

export const StyledInput = styled(Input)`
  padding: 8px 8px 8px 20px;
  background: rgba(3, 6, 9, 0.6);
  border-radius: 11px;
  width: 100%;
  font-size: 16px;
  font-weight: 400;
  line-height: 19px;
  color: #ffffff;
  height: 44px;
`;

export const UnhashButton = styled(Button)`
  border: 2px solid ${({ theme }) => theme.colors.success};
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0.3) 100%);
  border-radius: 6.75px;

  font-size: 14px;
  font-weight: 600;
  line-height: 17px;
  letter-spacing: 0.16em;
  color: white;

  padding: 9px 13px;
  gap: 10px;

  text-transform: uppercase;
`;

export default React.memo(Unhash);
