import React, { useEffect, useState } from 'react';

import { Flex } from 'components/Box';
import { Modal, ModalProps } from 'components/Modal';
import { Text } from 'components/Text';
import styled from 'styled-components';

interface AuthorizeModalProps extends ModalProps {
  setAuthorized: React.Dispatch<React.SetStateAction<boolean>>;
}

export const PASSWORD = process.env.REACT_APP_STAGE === 'dev' && 'DUEL2023';

export default function AuthorizeModal({
  setAuthorized,
  ...props
}: AuthorizeModalProps) {
  const [value, setValue] = useState('');

  useEffect(() => {
    if (value === PASSWORD) {
      localStorage.setItem('duel', PASSWORD);
      setAuthorized(true);
      props.onDismiss && props.onDismiss();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value]);

  return (
    <Modal {...props} hideCloseButton>
      <Flex
        px={'30px'}
        py="35px"
        border="2px solid #7389a9"
        borderRadius={'17px'}
        background="linear-gradient(180deg, #132031 0%, #1A293D 100%)"
        flexDirection="column"
        justifyContent="space-around"
        gap={16}
      >
        <Text color="#8192AA">
          To access duel.win please enter the launch codes
        </Text>
        <StyledInput
          type="password"
          value={value}
          onChange={e => setValue(e.target.value)}
        />
      </Flex>
    </Modal>
  );
}

const StyledInput = styled.input`
  font-family: 'Inter';
  font-style: normal;
  font-weight: 400;
  font-size: 14px;
  line-height: 17px;
  color: #8192aa;
  flex-grow: 1;
  border: none;
  outline: none;
  background: ${({ theme }) => theme.colors.secondary};
  border-radius: 10px;
  padding: 5px 20px;

  &::placeholder {
    color: #8192aa;
  }

  ::-webkit-outer-spin-button,
  ::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
  -moz-appearance: textfield;
`;
