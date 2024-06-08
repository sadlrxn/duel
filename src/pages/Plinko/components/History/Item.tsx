import { useMemo } from 'react';
import styled from 'styled-components';

import { Flex, Span } from 'components';

// import { FairnessModal } from '../Modal';

interface ItemProps {
  multiplier?: number;
  roundId?: number;
}

const HistoryItem = styled(Flex)`
  cursor: pointer;
  &:hover {
    -webkit-filter: brightness(130%);
    filter: brightness(130%);
  }
`;

export default function Item({ multiplier = 1 /*, roundId = 0*/ }: ItemProps) {
  // const [onFairnessModal] = useModal(<FairnessModal roundId={roundId} />, true);
  const [color, background] = useMemo(() => {
    let color = '#FA00FF';
    let background =
      'linear-gradient(90deg, rgba(219, 0, 255, 0.26) 0.01%, rgba(128, 0, 255, 0.26) 100%)';
    if (multiplier < 2) {
      color = '#00C2FF';
      background =
        'linear-gradient(90deg, rgba(62, 162, 255, 0.26) 0%, rgba(0, 209, 255, 0.26) 0.01%, rgba(0, 133, 255, 0.26) 100%)';
    } else if (multiplier < 10) {
      color = '#4FFF8B';
      background =
        'linear-gradient(90deg, rgba(79, 255, 139, 0.26) 0.01%, rgba(0, 160, 83, 0.26) 100%)';
    }
    return [color, background];
  }, [multiplier]);

  return (
    <HistoryItem
      height="30px"
      alignItems="center"
      justifyContent="center"
      overflow="hidden hidden"
      border={`1px solid ${color}`}
      borderRadius="5px"
      background={background}
      // onClick={onFairnessModal}
    >
      <Span
        color="white"
        fontWeight={600}
        fontSize="12px"
        lineHeight={1}
        padding="7.5px 10px"
      >
        {multiplier.toFixed(2)}x
      </Span>
    </HistoryItem>
  );
}
