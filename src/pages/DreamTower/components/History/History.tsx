import { DreamtowerHistoryRound } from 'api/types/dreamtower';
import { Avatar, BoxProps, Flex, Chip, Text, useModal } from 'components';
import dayjs from 'dayjs';
import { FC, memo } from 'react';
import { convertBalanceToChip } from 'utils/balance';
import { formatUserName } from 'utils/format';
import FairnessModal from '../Modal/Fairness';
import { HistoryItemContainer } from './History.styles';
import VerifyButton from './VerifyButton';

interface DreamtowerHistoryProps extends BoxProps {
  round: DreamtowerHistoryRound;
  selected: boolean;
  onClick: () => void;
}

const HistoryItem: FC<DreamtowerHistoryProps> = ({
  round,
  selected,
  onClick
}) => {
  const [onFairnessModal] = useModal(
    <FairnessModal roundId={round.roundId} />,
    true
  );
  return (
    <HistoryItemContainer selected={selected} onClick={onClick}>
      <Flex
        flexDirection="row"
        alignItems="center"
        justifyContent="space-between"
        gap={9}
        width="100%"
      >
        <Flex flexDirection="row" alignItems="center" gap={13.5} width="40%">
          <Avatar
            userId={round.user.id}
            name={round.user.name}
            image={round.user.avatar}
            border="none"
            borderRadius="8px"
            padding="0px"
            size="30px"
          />
          <Text fontSize={'14px'} fontWeight={500} color="textWhite">
            {formatUserName(round.user.name)}
          </Text>
        </Flex>
        <Flex flexDirection="column" alignItems="start" width="15%">
          <Text
            fontSize={'14px'}
            fontWeight={500}
            color="textWhite"
            opacity={0.5}
          >
            BET
          </Text>
          <Chip
            color="#FFFFFF"
            chipType={round.paidBalanceType}
            price={Number(convertBalanceToChip(round.betAmount)).toFixed(2)}
            fontWeight={600}
            fontSize="14px"
          />
        </Flex>
        <Flex flexDirection="column" alignItems="center" width="25%" gap={3}>
          <Text
            fontSize={'14px'}
            fontWeight={500}
            color="textWhite"
            opacity={0.5}
          >
            MULTIPLIER
          </Text>
          <Text fontSize={'14px'} fontWeight={500} color="textWhite">
            {round.multiplier.toFixed(2)}
          </Text>
        </Flex>
        <Flex flexDirection="column" alignItems="end">
          <Text
            fontSize={'14px'}
            fontWeight={500}
            color="textWhite"
            opacity={0.5}
          >
            PAYOUT
          </Text>
          <Chip
            color="#FFFFFF"
            chipType={round.paidBalanceType}
            price={Number(convertBalanceToChip(round.profit)).toFixed(2)}
            fontWeight={600}
            fontSize="14px"
          />
        </Flex>
      </Flex>
      <Flex
        flexDirection={'row'}
        justifyContent="space-between"
        gap={9}
        width="100%"
      >
        <Flex
          flexDirection="row"
          alignItems="center"
          justifyContent="space-between"
          width="30%"
        >
          <Text
            fontSize={'14px'}
            fontWeight={500}
            color="textWhite"
            opacity={0.5}
          >
            GAME ID
          </Text>
          <Text fontSize={'14px'} fontWeight={500} color="textWhite">
            {round.roundId}
          </Text>
        </Flex>
        <Flex
          flexDirection="row"
          alignItems="center"
          justifyContent="space-between"
          width="40%"
        >
          <Text
            fontSize={'14px'}
            fontWeight={500}
            color="textWhite"
            opacity={0.5}
          >
            DATE
          </Text>
          <Text fontSize={'14px'} fontWeight={500} color="textWhite">
            {dayjs(round.time).format('MMM DD, hh:mm A')}
          </Text>
        </Flex>
        <Flex flexDirection="row" justifyContent="end" width="10%">
          <VerifyButton onClick={onFairnessModal} />
        </Flex>
      </Flex>
    </HistoryItemContainer>
  );
};

export default memo(HistoryItem);
