import React, { useMemo } from 'react';
import styled from 'styled-components';

import { Box, Flex } from 'components/Box';
import { Modal, ModalProps } from 'components/Modal';
import { Span, Text } from 'components/Text';
import { Button } from 'components/Button';
import { BlockIcon } from 'components/Icon';

import { useAppDispatch, useAppSelector } from 'state';
import { setCountryCode } from 'state/user/actions';
import { countrycodes } from 'config';

export interface JurisdictionsModalProps extends ModalProps {}

export default function JurisdictionsModal({
  onDismiss,
  ...props
}: JurisdictionsModalProps) {
  const dispatch = useAppDispatch();
  const code = useAppSelector(state => state.user.code);

  const countryName = useMemo(() => {
    if (code) {
      const country = countrycodes.find(
        c => c.code.toLowerCase() === code.toLowerCase()
      );
      return country ? country.name : '';
    }
    return '';
  }, [code]);

  return (
    <Modal
      onDismiss={() => {
        dispatch(setCountryCode(''));
        onDismiss && onDismiss();
      }}
      {...props}
    >
      <Container>
        <Box
          background="linear-gradient(180deg, #132031 0%, #1A293C 100%)"
          px="45px"
          py="40px"
          borderRadius="20px"
        >
          <Flex gap={18} mb="48px" alignItems="center">
            <BlockIcon />
            <Text
              fontSize="20px"
              fontWeight={600}
              color="white"
              lineHeight="24px"
              letterSpacing="0.18em"
            >
              Prohibited Jurisdictions
            </Text>
          </Flex>
          <Text
            fontSize="15px"
            fontWeight={400}
            lineHeight="18px"
            color="#B2D1FF"
            mb="46px"
          >
            <p>
              Based on your IP address, it seems you are currently located
              within{' '}
              <Span fontWeight={700} fontStyle="italic">
                {countryName},
              </Span>{' '}
              a prohibited jurisdiction.
            </p>
            <p>
              Duel is an online gambling company that is only available for a
              limited number of jurisdictions. Users should consult{' '}
              <Span fontWeight={700}>this list</Span> to determine if they are
              prohibited from playing on the platform.
            </p>
            {/* <p>
              We apologize for any inconvenience this may cause. You are still welcome to browse the site, but please be aware that you will not be able to place any bets.
            </p> */}
          </Text>
          {/* <Flex width="100%" justifyContent="center" alignItems="center">
            <Button variant="secondary" background="#1A5032" px="30px" py="12px" fontSize="14px" lineHeight="18px" fontWeight={600} color="success">
              See List of Prohibited Jurisdictions
            </Button>
          </Flex> */}
        </Box>
      </Container>
    </Modal>
  );
}

const Container = styled(Box)`
  background: linear-gradient(180deg, #6a7f9e 0%, rgba(106, 127, 158, 0) 100%);
  padding: 1px;
  overflow: hidden;
  border-radius: 20px;
  min-width: 350px;
  max-width: 630px;
  width: 100%;
`;
