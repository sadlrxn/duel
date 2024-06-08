import { Box, Flex } from 'components/Box';
import { formatNumber } from 'utils/format';

export const Icon = ({ img, size = 24 }: { img: string; size?: number }) => {
  return (
    <>
      <img src={img} alt="" width={size} height={size} />
    </>
  );
};

export const SelectItem = (
  label: string,
  img: string,
  size: number = 24,
  balance: number
) => {
  return (
    <Flex justifyContent="space-between">
      <Flex alignItems="center" gap={10}>
        <Icon img={img} size={size} />
        {label.slice(0, 3) === 'SOL' ? 'SOLANA' : label}
      </Flex>
      {`${formatNumber(balance)} ${label}`}
    </Flex>
  );
};
