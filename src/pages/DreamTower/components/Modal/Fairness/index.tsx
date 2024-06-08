import { useState, useMemo } from 'react';
import styled from 'styled-components';

import {
  Box,
  Flex,
  Span,
  FairnessIcon,
  ModalProps,
  Button,
  Input
} from 'components';
import { useFetchRoundInfo } from 'hooks';
import { DreamtowerFairData, initialFair } from 'api/types/dreamtower';

import {
  Container,
  TitleContainer,
  Title,
  Description,
  StyledModal
} from './styles';
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

  const { data, error } = useFetchRoundInfo('dreamtower', round);

  const { gameData } = useMemo(() => {
    const isError = error ? true : false;
    const isLoading = !isError && !data ? true : false;
    const gameData =
      !isError && !isLoading ? (data as DreamtowerFairData) : initialFair;
    return { gameData };
  }, [data, error]);

  return (
    <StyledModal {...props} onDismiss={onDismiss}>
      <Box
        p={['40px 20px', '40px 30px', '40px 40px', '40px 50px']}
        background="linear-gradient(180deg, #132031 0%, #1a293c 100%)"
        borderRadius="20px"
      >
        <Container>
          <TitleContainer>
            <FairnessIcon size={22} color="#4FFF8B" />
            DREAMTOWER FAIRNESS
          </TitleContainer>

          <Box overflow="auto" pl="2px" pr="10px">
            <Box>
              <Title>Game Details</Title>
              <Flex
                justifyContent="space-between"
                flexWrap="wrap"
                gap={14}
                mt="10px"
              >
                <Description>
                  Change the <Span fontWeight={700}>Game ID</Span> to verify
                  outcomes for other games.
                </Description>
                <Flex gap={6} width="max-content" flexDirection="column">
                  <Detail
                    title="Game ID"
                    text={round}
                    type="number"
                    setText={setRound}
                    enableCopy
                  />
                  <Span
                    color="#FF623B"
                    fontWeight={400}
                    fontSize="12px"
                    lineHeight="15px"
                    opacity={error ? 1 : 0}
                  >
                    The game with this ID doesn’t exist
                  </Span>
                </Flex>
              </Flex>
            </Box>

            <FairData roundId={+round} onDismiss={onDismiss} game={gameData} />
          </Box>
        </Container>
      </Box>
    </StyledModal>
  );
}

export const StyledInput = styled(Input)`
  padding: 8px 8px 8px 20px;
  background: rgba(3, 6, 9, 0.6);
  border-radius: 11px;
  width: 80%;
  font-size: 16px;
  font-weight: 400;
  line-height: 19px;
  color: #ffffff;
  height: 44px;
`;

export const StyledButton = styled(Button)`
  border: 2px solid ${({ theme }) => theme.colors.success};
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0.3) 100%);
  border-radius: 6.75px;

  font-size: 14px;
  font-weight: 600;
  line-height: 17px;
  letter-spacing: 0.16em;
  color: white;
  width: 20%;
  white-sapce: nowrap;

  padding-left: 13px;
  padding-right: 13px;
  gap: 10px;

  text-transform: uppercase;
`;
