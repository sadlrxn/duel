import { Box, Flex } from 'components/Box';
import { Span, Label, Text } from 'components/Text';
import { FlexFooter } from './styles';

import { ReactComponent as DiscordIcon } from 'assets/imgs/icons/discord.svg';
import { ReactComponent as TwitterIcon } from 'assets/imgs/icons/twitter.svg';
import { useAppSelector } from 'state';
import styled from 'styled-components';
import { DotIcon } from 'components/Icon';
import Logo from 'components/Icon/Logo';

const StyledTextLink = styled(Text)`
  color: #b2d1ff;
  font-weight: 600;
  font-size: 14px;
  &:hover {
    color: #4fff8b;
    /* text-decoration: underline; */
  }
`;

const StyledSpanLink = styled(Span)`
  &:hover {
    color: #4fff8b;
  }
`;

const StyledBox = styled(Box)`
  text-align: center;
  ${({ theme }) => theme.mediaQueries.md} {
    text-align: inherit;
  }
`;
const StyledFlex = styled(Flex)`
  justify-content: center;
  text-align: center;
  ${({ theme }) => theme.mediaQueries.md} {
    justify-content: inherit;
    text-align: inherit;
  }
`;

export default function Footer() {
  const { connected } = useAppSelector(state => state.socket);

  return (
    <FlexFooter>
      <StyledBox>
        <Logo />

        <Text color={'white'} mt="5px">
          Duel Games Corp © 2023
        </Text>

        <Text color={'#768BAD'} fontSize="10px" maxWidth={'400px'} mt="5px">
          duel.win is operated by Duelana Games B.V. registered under No. 161199
          at, Johan Van Walbeeckplein 2, Curacao. This website is licensed and
          regulated by Curaçao eGaming under Curaçao license No. 1668 JAZ issued
          by Curaçao eGaming.
        </Text>

        <Text color={'#768BAD'} fontSize="10px" maxWidth={'400px'} mt="10px">
          In order to register for this website, the user is required to accept
          the{' '}
          <a
            href="https://docs.duel.win/more-info/terms-and-conditions"
            target="_blank"
            rel="noreferrer"
          >
            <StyledSpanLink color={'#B7D1FB'}>
              General Terms and Conditions
            </StyledSpanLink>
          </a>
          . In the event the General Terms and Conditions are updated, existing
          users may choose to discontinue using the products and services before
          the said update shall become effective, which is a minimum of two
          weeks after it has been announced.
        </Text>

        <StyledFlex gap={20} my="10px" alignItems={'center'}>
          <a href="https://discord.gg/duel" rel="noreferrer" target={'_blank'}>
            <DiscordIcon />
          </a>

          <a
            href="https://twitter.com/DuelCasino"
            rel="noreferrer"
            target={'_blank'}
          >
            <TwitterIcon />
          </a>
        </StyledFlex>

        <StyledFlex alignItems={'center'} gap={5}>
          <Text color={'#768BAD'} fontSize="14px" fontWeight={500}>
            Server Status
          </Text>
          <DotIcon color={connected ? '#4FFF8B' : '#FF4F4F'} />
          <Label
            color={connected ? 'success' : 'warning'}
            fontSize="14px"
            fontWeight={600}
          >
            {connected ? 'Connected' : 'Disconnected'}
          </Label>
        </StyledFlex>
      </StyledBox>

      <StyledFlex flexDirection={'column'} gap={5}>
        <Text color={'#768BAD'} fontSize="12px" fontWeight={700}>
          LEGAL
        </Text>
        <a
          href="https://docs.duel.win/more-info/terms-and-conditions"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Terms & Conditions</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/privacy-policy"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Privacy Policy</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/privacy-policy"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Privacy & Management of Personal Data</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/kyc-policies"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>KYC Policies</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/aml-policy"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>AML</StyledTextLink>
        </a>
      </StyledFlex>

      <StyledFlex flexDirection={'column'} gap={5}>
        <Text color={'#768BAD'} fontSize="12px" fontWeight={700}>
          RESOURCES
        </Text>

        <a
          href="https://docs.duel.win/more-info/self-exclusion"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Self-Exclusion</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/self-exclusion"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Responsible Gambling</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/duels-provable-fairness"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Fairness & RNG Testing Methods</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/accounts-payouts-and-bonuses"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Account, Payouts, and Bonuses</StyledTextLink>
        </a>

        <a
          href="https://docs.duel.win/more-info/dispute-resolution"
          target="_blank"
          rel="noreferrer"
        >
          <StyledTextLink>Dispute Resolution</StyledTextLink>
        </a>
      </StyledFlex>
    </FlexFooter>
  );
}
