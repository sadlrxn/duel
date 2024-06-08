import React, { FC } from 'react';
import styled from 'styled-components';
import DataTable, {
  TableColumn,
  TableStyles
} from 'react-data-table-component';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { Span } from 'components/Text';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex } from 'components/Box';
import pagination from './pagination';
import { useAppSelector } from 'state';
import { Log } from 'state/log/actions';
import dayjs from 'dayjs';
import PhantomImg from 'assets/imgs/icons/phantom.png';
import { formatNumber } from 'utils/format';
import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { ReactComponent as UpIcon } from 'assets/imgs/icons/up.svg';
import { ReactComponent as DownIcon } from 'assets/imgs/icons/down.svg';
import { imageProxy } from 'config';

const StyledNft = styled.img`
  width: 28px;
  height: 28px;
  border: 2px solid #203044;
  border-radius: 8px;

  &:not(:first-child) {
    margin-left: -15px;
  }
`;

const columns: TableColumn<Log>[] = [
  {
    name: 'type',
    cell: row => {
      if (row.type === 'Deposit')
        return (
          <>
            <DownIcon />
            <Span ml="10px" color={'#4FFF8B'}>
              {row.type}
            </Span>
          </>
        );
      else
        return (
          <>
            <UpIcon />
            <Span ml="10px" color={'#C78FFF'}>
              {row.type}
            </Span>
          </>
        );
    },
    minWidth: '140px !important'
  },
  {
    name: 'data',
    cell: row => {
      if (row.data && typeof row.data === 'number') {
        return (
          <Flex alignItems={'center'}>
            <CoinIcon />
            <Span ml="8px" color={'#FFF6CA'}>
              {formatNumber(row.data)}
            </Span>
          </Flex>
        );
      } else if (row.data && typeof row.data === 'object') {
        return row.data.map(value => (
          <StyledNft src={imageProxy(300) + value} alt="nft" key={value} />
        ));
      } else return;
    },
    minWidth: '100px !important'
  },
  {
    name: 'time',
    selector: row => dayjs(new Date(row.time)).format('MM/DD/YY HH:mm'),
    minWidth: '140px !important'
  },
  {
    name: 'signature',
    cell: row =>
      row.signature && (
        <a
          href={`https://solscan.io/tx/${row.signature}`}
          rel="noreferrer"
          target={'_blank'}
        >
          <Flex alignItems={'center'}>
            <img src={PhantomImg} width="20px" height="20px" alt="phantom" />
            <CopyToClipboard text={row.signature}>
              <Span ml="8px">
                {row.signature.slice(0, 4)}...
                {row.signature.slice(-4)}
              </Span>
            </CopyToClipboard>
          </Flex>
        </a>
      ),
    minWidth: '145px !important',
    grow: 1
  },
  {
    name: 'status',
    cell: row => (
      <Span
        color={
          row.status === 'Success'
            ? '#4FFF8B'
            : row.status === 'Failed'
            ? '#F24822'
            : '#FFCD4B'
        }
      >
        {row.status}
      </Span>
    )
  }
];

const customStyles: TableStyles = {
  table: {
    style: {
      fontFamily: 'Inter',

      backgroundColor: 'transparent'
    }
  },
  rows: {
    style: {
      border: 'none !important',
      //   borderBottomColor: "transparent !important",
      fontSize: '14px',
      color: '#697E9C',
      backgroundColor: 'transparent',
      padding: '0px 5px !important'
    }
  },
  pagination: {
    style: {
      backgroundColor: 'transparent',
      color: '#697E9C',
      fontFamily: 'Inter',
      fontSize: '14px',
      borderBottom: '1px solid #2A3D57'
    }
  },

  noData: {
    style: {
      backgroundColor: 'transparent',
      color: '#697E9C'
    }
  }
};

const LogModal: FC<ModalProps> = ({ ...props }) => {
  const { logs } = useAppSelector(state => state.log);
  return (
    <Modal {...props}>
      <Box
        px={'30px'}
        pt="35px"
        border="2px solid #7389a9"
        borderRadius={'17px'}
        background="linear-gradient(180deg, #132031 0%, #1A293D 100%)"
        minHeight={'400px'}
        width="700px"
      >
        <Flex flexDirection={'column-reverse'}>
          <DataTable
            columns={columns}
            data={logs}
            pagination
            customStyles={customStyles}
            noTableHead
            paginationComponent={pagination}
          />
        </Flex>
      </Box>
    </Modal>
  );
};

export default LogModal;
