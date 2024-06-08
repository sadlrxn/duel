import { useState, useMemo } from 'react';
import styled from 'styled-components';
import dayjs from 'dayjs';

import { Box, Flex, Span, FairnessIcon, ModalProps } from 'components';

import {
  Container,
  TitleContainer,
  Title,
  Description,
  StyledModal
} from './styles';
import { useFetchRoundInfo } from 'hooks';
import { CrashFairData, initialFair } from 'api/types/crash';

import Detail from './Detail';
import FairData from './FairData';

interface FairnessModalProps extends ModalProps {
  roundId?: number;
}

export default function FairnessModal({
  roundId = 0,
  onDismiss,
  ...props
}: FairnessModalProps) {
  const [round, setRound] = useState(roundId.toString());

  const { data, error } = useFetchRoundInfo('crash', round);

  const { gameData } = useMemo(() => {
    const isError = error ? true : false;
    const isLoading = !isError && !data ? true : false;
    const gameData =
      !isError && !isLoading ? (data as CrashFairData) : initialFair;
    return { gameData };
  }, [data, error]);

  const formatedDate = dayjs(gameData.date).format('MM/DD/YYYY hh:mm A');

  return (
    <StyledModal {...props} onDismiss={onDismiss}>
      <Box
        p={['40px 20px', '40px 20px', '40px 30px', '40px 30px', '40px 50px']}
        background="linear-gradient(180deg, #132031 0%, #1a293c 100%)"
        borderRadius={'20px'}
      >
        <TitleContainer>
          <FairnessIcon size={22} color="#4FFF8B" />
          CRASH FAIRNESS
        </TitleContainer>
        <Container>
          <IDWrapper
            justifyContent="space-between"
            flexWrap="wrap"
            gap={14}
            mb="17px"
            pr="10px"
            pl="2px"
          >
            <Box>
              <Title mb="10px">Game Details</Title>
              <Description>
                Change the <Span fontWeight={700}>Game ID</Span> to verify
                outcomes for other games.
              </Description>
            </Box>

            <Detail
              title="Game ID"
              text={round}
              type="number"
              setText={setRound}
              enableCopy
            />
          </IDWrapper>

          <Flex
            justifyContent="space-between"
            flexWrap="wrap"
            gap={1}
            pr="10px"
            pl="2px"
          >
            <Detail
              title="Round Seed"
              text={gameData.seed}
              type="text"
              readOnly
              enableCopy
            />
            <Detail
              title="Crashed At"
              text={String(gameData.outcome) + 'x'}
              type="text"
              readOnly
              enableCopy
            />
            <Detail
              title="Date"
              text={formatedDate}
              type="text"
              readOnly
              enableCopy
            />
          </Flex>

          <FairData roundId={+round} onDismiss={onDismiss} game={gameData} />
        </Container>
      </Box>
    </StyledModal>
  );
}

const IDWrapper = styled(Flex)`
  & > div:first-child {
    order: 3;
  }

  ${({ theme }) => theme.mediaQueries.md} {
    & > div:first-child {
      order: 0;
    }
  }
`;
