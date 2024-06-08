import { useState } from 'react';

import { Box, Flex, Span, FairnessIcon, ModalProps } from 'components';

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
            JACKPOT FAIRNESS
          </TitleContainer>

          <Box overflow="auto" pl="2px" pr="10px">
            <Box mb="17px">
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
                <Detail
                  title="Game ID"
                  text={round}
                  type="number"
                  setText={setRound}
                  enableCopy
                />
              </Flex>
            </Box>

            <FairData roundId={+round} onDismiss={onDismiss} />
          </Box>
        </Container>
      </Box>
    </StyledModal>
  );
}
