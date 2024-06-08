import React from 'react';
import styled from 'styled-components';

import { Box, Flex, Span } from 'components';

const MultiplierItem = ({
  multiplier,
  show
}: {
  multiplier: number;
  show: boolean;
}) => {
  return (
    <Flex
      gap={10}
      height="1px"
      maxHeight="1px"
      alignItems="center"
      width="60px"
      minWidth="60px"
    >
      <Box
        height="1px"
        width="9px"
        background="#D9D9D9"
        style={{ opacity: show ? 1 : 0.5 }}
      />
      {show && <MultiSpan>{multiplier.toFixed(2)}x</MultiSpan>}
    </Flex>
  );
};

interface MultiplierRulerProps {
  min?: number;
  max?: number;
}

export default function MultiplierRuler({
  min = 1,
  max = 2
}: MultiplierRulerProps) {
  return (
    <ContainerWrapper>
      <Container>
        {Array(21)
          .fill(1)
          .map((_, index) => {
            return (
              <MultiplierItem
                key={`crash_multiplier_item_${index}`}
                show={index % 5 === 0}
                multiplier={min + ((max - min) / 20) * index}
              />
            );
          })}
      </Container>
    </ContainerWrapper>
  );
}

const MultiSpan = styled(Span)`
  font-weight: 700;
  font-size: 12px;
  line-height: 1.3;
  color: white;

  .width_700 & {
    font-size: 15px;
  }
`;

const Container = styled(Flex)`
  flex-direction: column-reverse;
  justify-content: space-between;
  align-items: flex-end;
  width: 100%;
  height: 90%;
`;

const ContainerWrapper = styled(Flex)`
  width: 100%;
  height: 100%;
  align-items: flex-end;
`;
