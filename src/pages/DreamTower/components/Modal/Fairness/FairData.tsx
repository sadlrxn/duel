import React from 'react';
import { Link } from 'react-router-dom';
import { Box, Flex, Grid, Span } from 'components';
import { DreamtowerFairData } from 'api/types/dreamtower';

import { Description, Divider, GrayButton, Title } from './styles';

import Detail from './Detail';
import Tower from './Tower';

interface FairDataProps {
  roundId?: number;
  onDismiss?: any;
  game: DreamtowerFairData;
}

const FairData: React.FC<FairDataProps> = ({ onDismiss, game }) => {
  return (
    <>
      <Flex gap={14} mt="10px">
        <Detail title="Player" text={game.user.name} enableCopy readOnly />
        <Detail
          title="Difficulty"
          text={game.difficulty.level}
          enableCopy
          readOnly
        />
      </Flex>
      <Divider />
      <Box mb="17px" gap={50}>
        <Title>Seed Pair</Title>
        <Flex
          justifyContent="space-between"
          flexWrap="wrap"
          gap={14}
          mt="10px"
          mb="17px"
        >
          <Description>
            Game results are generated from a combinations of the client seed
            (the player playing the game) and the server seed (hashed as long as
            the client seed is active). Every time you update your client seed a
            new server seed is generated, hashed and paired to create a seed
            pair - your previous client seed becomes expired and the server seed
            unhashed. To verify the outcome of a previous game the player must
            generate a new client seed,{' '}
            <Span
              fontWeight={700}
              color="#b9d2fd"
              onClick={() => {
                onDismiss && onDismiss();
              }}
            >
              <Link to="/fair" state={{ room: 'seed' }}>
                click here
              </Link>
            </Span>{' '}
            to generate a new client seed.
          </Description>
        </Flex>
        {/* <Flex flexDirection="column" gap={8}>
                <Text fontSize="16px" lineHeight="19px" color="#768BAD">
                  Client Seed
                </Text>
                <Flex
                  flexDirection="row"
                  justifyContent="space-between"
                  gap={22}
                >
                  <StyledInput
                    value={clientSeed}
                    onChange={e => {
                      setClientSeed(e.target.value);
                    }}
                  ></StyledInput>
                  <StyledButton onClick={handleClick}>Change Seed</StyledButton>
                </Flex>
              </Flex> */}
        <Detail
          title="Client Seed"
          text={game.clientSeed}
          status={game.serverSeed ? 'Expired' : 'Active'}
          enableCopy
          readOnly
        />
        <Detail
          title="Server Seed (Hashed)"
          text={game.serverSeedHash}
          enableCopy
          readOnly
        />
        <Detail
          title="Server Seed (UnHashed)"
          text={game.serverSeed ?? ''}
          placeholder="Generate a new client seed to see unhashed server seed."
          enableCopy
          readOnly
        />
        <Grid
          gridTemplateColumns="repeat(auto-fill, minmax(280px, 1fr))"
          gridColumnGap="30px"
          gridRowGap="17px"
        >
          <Detail
            title="Nonce"
            text={game.nonce.toString()}
            enableCopy
            readOnly
          />
          <Detail
            title="Total Bets Made With Pair"
            text={game.seedNonce.toString()}
            placeholder="Game In Progress..."
            enableCopy
            readOnly
          />
        </Grid>
      </Box>
      <Flex justifyContent="center" pt="50px" mb="17px">
        <Tower tower={game.tower} blocksInRow={game.difficulty.blocksInRow} />
      </Flex>

      <Flex gap={30} flexWrap="wrap" pt="18px">
        <GrayButton width={['100%', '100%', '192px']} onClick={onDismiss}>
          <Link to="/fair">Provably Fair</Link>
        </GrayButton>
        <GrayButton
          width={['100%', '100%', '216px']}
          onClick={onDismiss}
          disabled={!game.serverSeed}
        >
          <Link to="/fair" state={{ gameType: 'dreamtower', gameData: game }}>
            Advanced Verification
          </Link>
        </GrayButton>
      </Flex>
    </>
  );
};

export default React.memo(FairData);
