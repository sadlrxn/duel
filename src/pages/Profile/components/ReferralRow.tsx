import React, { FC } from 'react';
import { Badge, Button, Chip, Flex, TrashIcon, useModal } from 'components';
import ConfirmDeleteModal from './ConfirmDeleteModal';
import copy from 'copy-to-clipboard';
import toast from 'utils/toast';
import AffiliateUserModal from './AffiliateUserModal';
import { convertBalanceToChip } from 'utils/balance';

export default function ReferralRow({
  data: { code, reward, totalEarned, totalWagered, userCnt, rate },
  onClaim,
  onDelete
}: {
  data: {
    code: string;
    reward: number;
    totalEarned: number;
    totalWagered: number;
    userCnt: number;
    rate: number;
  };
  onClaim: (codes: string[]) => void;
  onDelete: (code: string) => void;
}) {
  const [onPresentCodeUsers] = useModal(<AffiliateUserModal code={code} />);
  const [onPresentConfirm] = useModal(
    <ConfirmDeleteModal onDelete={() => onDelete(code)} />
  );

  return (
    <tr key={code} onClick={onPresentCodeUsers}>
      <td>{code}</td>
      <td>{userCnt}</td>
      <td>{rate}%</td>
      <td>
        <Chip price={convertBalanceToChip(totalWagered)} />
      </td>
      <td>
        <Chip price={convertBalanceToChip(totalEarned)} />
      </td>
      <td>
        <Button
          background={'#1A5032'}
          borderRadius="5px"
          p="5px 10px"
          disabled={reward < 10 ? true : false}
          onClick={() => onClaim([code])}
        >
          <Chip price={convertBalanceToChip(reward)} color="#4FFF8B" />
        </Button>
      </td>
      <td>
        <Flex
          alignItems="center"
          justifyContent="space-between"
          width="100%"
          gap={30}
        >
          <Button
            color="#768BAD"
            borderRadius={'5px'}
            background="#242F42"
            p="10px 10px"
            fontWeight={600}
            onClick={() => {
              copy(`https://duel.win/r/${code}`);
              toast.success('Copy to clipboard success.');
            }}
          >
            Copy Link
          </Button>
          <Button
            color="#768BAD"
            borderRadius={'5px'}
            background="#242F42"
            p="10px 10px"
            fontWeight={600}
            onClick={onPresentConfirm}
          >
            <TrashIcon />
          </Button>
        </Flex>
      </td>
    </tr>
  );
}
