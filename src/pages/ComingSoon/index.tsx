import React from 'react';
import styled from 'styled-components';
import { AuthorizeModal, Box, Button, Flex, Text, useModal } from 'components';
import FlipClockCountdown from '@leenguyen/react-flip-clock-countdown';
import { TypeAnimation } from 'react-type-animation';
import rocket from 'assets/imgs/comingsoon/rocket.png';
import planet from 'assets/imgs/comingsoon/planet.png';
import Logo from 'components/Icon/Logo';
import Plus from 'components/Icon/Plus';
import Discord from 'components/Icon/Discord';
import Symbol from 'components/Icon/Symbol';
import '@leenguyen/react-flip-clock-countdown/dist/index.css';
import Twitter from 'components/Icon/Twitter';

const Container = styled(Flex)`
  height: 100vh;
  align-items: center;
  justify-content: center;
  background-color: #000;
  background-image: url(${planet});
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center center;
`;

const LogoBox = styled(Box)`
  position: absolute;
  top: 2.5rem;
  left: 2.5rem;
`;

const RightPlusBox = styled(Box)`
  position: absolute;
  top: 2.5rem;
  right: 2.5rem;
`;

const TextBox = styled(Box)`
  position: absolute;
  bottom: 80px;
  left: 85px;
  display: none;

  ${({ theme }) => theme.mediaQueries.md} {
    display: block;
  }
`;

const BottomPlusBox = styled(Box)`
  position: absolute;
  bottom: 2.5rem;
  left: 2.5rem;
`;

const SymbolBox = styled(Box)`
  position: absolute;
  bottom: 2.5rem;
  right: 2.5rem;
`;

export default function ComingSoon({
  setAuthorized
}: {
  setAuthorized: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [onPresentModal] = useModal(
    <AuthorizeModal setAuthorized={setAuthorized} />,
    true
  );

  const handleAuthorize = () => {
    onPresentModal();
  };
  return (
    <Container padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}>
      <LogoBox>
        <Logo />
      </LogoBox>
      <RightPlusBox>
        <Plus />
      </RightPlusBox>
      <Flex
        justifyContent={'center'}
        alignItems="center"
        flexDirection={'column'}
        gap={12}
      >
        <img src={rocket} alt="" />

        <Text
          color={'white'}
          textTransform="uppercase"
          fontSize={'12px'}
          fontFamily="Termina"
        >
          The countown has begun
        </Text>

        <Box mt={'20px'}>
          <FlipClockCountdown
            to={'2022-12-16T20:00:00+00:00'}
            labels={['DAYS', 'HOURS', 'MINUTES', 'SECONDS']}
            labelStyle={{
              fontSize: 10,
              fontWeight: 500,
              textTransform: 'uppercase',
              fontFamily: 'Termina',
              paddingTop: '5px'
            }}
            digitBlockStyle={{
              width: 35,
              height: 55,
              fontSize: 40,
              fontWeight: 700,
              background: '#D9D9D9',
              color: '#333333'
            }}
            dividerStyle={{ color: 'black', height: 1 }}
            separatorStyle={{ color: 'white', size: '8px' }}
          />
        </Box>

        <Flex alignItems={'center'} justifyContent="center" gap={40} mt="20px">
          <a
            href="https://discord.gg/duel"
            rel="noreferrer"
            target={'_blank'}
            className="text-white hover:text-[#4FFF8B]"
          >
            <Discord />
          </a>
          <a
            href="https://twitter.com/DuelCasino"
            rel="noreferrer"
            target={'_blank'}
            className="text-white hover:text-[#4FFF8B]"
          >
            <Twitter />
          </a>
        </Flex>
        <Button p="10px" background="#D9D9D9" onClick={handleAuthorize}>
          Input Launch Codes
        </Button>
      </Flex>

      <TextBox>
        <Text color={'#2E2E2E'} lineHeight="1em" width={'384px'}>
          DUEL.WIN LAUNCH SEQUENCE Systems
          <TypeAnimation
            sequence={[
              'starting up...', // Types 'One'
              1000, // Waits 1s
              ''
            ]}
            wrapper="p"
            cursor={true}
            repeat={Infinity}
          />
          <br />
          <TypeAnimation
            sequence={[
              'Loading propelent...', // Types 'One'
              1000, // Waits 1s
              ''
            ]}
            wrapper="p"
            cursor={true}
            repeat={Infinity}
          />
          [100%] Propelent loaded [✓] <br /> <br />
          <TypeAnimation
            sequence={[
              'Confirming launch...', // Types 'One'
              1000, // Waits 1s
              ''
            ]}
            wrapper="p"
            cursor={true}
            repeat={Infinity}
          />
          Launch Confirmed [✓] <br /> <br />
          Bots have taken over the countdown. [✓] <br />
          ALL SYSTEMS GO. [✓]
        </Text>
      </TextBox>
      <BottomPlusBox>
        <Plus />
      </BottomPlusBox>

      <SymbolBox>
        <Symbol />
      </SymbolBox>
    </Container>
  );
}
