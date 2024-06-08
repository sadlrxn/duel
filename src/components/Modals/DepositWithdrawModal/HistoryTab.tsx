import React, { FC, useCallback, useState } from 'react';
import styled from 'styled-components';
import DataTable, {
  TableColumn,
  TableStyles
} from 'react-data-table-component';
import { Span, Text } from 'components/Text';
import { Box, Flex } from 'components/Box';
import Pagination from '../LogModal/pagination';
import { imageProxy } from 'config';

import dayjs from 'dayjs';
import Select from 'components/Select';
import { formatNumber } from 'utils/format';
import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { ReactComponent as UpIcon } from 'assets/imgs/icons/up.svg';
import { ReactComponent as DownIcon } from 'assets/imgs/icons/down.svg';
import useSWR from 'swr';
import api from 'utils/api';
import useModal from 'components/Modal/useModal';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';
import { SolscanIcon } from 'components/Icon';
import { Button } from 'components';
import { ClipLoader } from 'react-spinners';
import { convertBalanceToChip } from 'utils/balance';

interface HistoryData {
  nftDetail: {
    collectionImage: string;
    collectionName: string;
    image: string;
    mintAddress: string;
    name: string;
    price: number;
  }[];
  solDetail: {
    solAmount: number;
    finalBalance: number;
    usdAmount: number;
  };
  status: string;
  time: string;
  txId: string;
  type: string;
}

const StyledNft = styled.img`
  width: 28px;
  height: 28px;
  border: 2px solid #203044;
  border-radius: 8px;

  &:not(:first-child) {
    margin-left: -15px;
  }
`;

const columns: TableColumn<HistoryData>[] = [
  {
    name: 'Transaction Type',
    cell: row => {
      if (
        row.type === 'deposit_sol' ||
        row.type === 'deposit_bonk' ||
        row.type === 'deposit_usdc' ||
        row.type === 'deposit_nft'
      )
        return (
          <>
            <DownIcon />
            <Span ml="10px" color={'#4FFF8B'}>
              Deposit
            </Span>
          </>
        );
      else
        return (
          <>
            <UpIcon />
            <Span ml="10px" color={'#C78FFF'}>
              Withdraw
            </Span>
          </>
        );
    },
    minWidth: '200px !important'
  },
  {
    name: 'Value',
    cell: row => {
      if (row.nftDetail.length === 0) {
        return (
          <Flex alignItems={'center'}>
            <CoinIcon />
            <Span ml="10px" color={'#fff'}>
              {formatNumber(convertBalanceToChip(row.solDetail.usdAmount))}
            </Span>
          </Flex>
        );
      } else {
        return row.nftDetail.map(value => (
          <StyledNft
            src={imageProxy(300) + value.image}
            alt="nft"
            key={value.mintAddress}
          />
        ));
      }
    },
    minWidth: '130px !important'
  },
  {
    name: 'Date',
    selector: row => dayjs(new Date(row.time)).format('MM/DD/YY HH:mm'),
    minWidth: '150px !important'
  },
  {
    name: 'Signature',
    cell: row => (
      <Flex alignItems={'center'}>
        <SolscanIcon />
        {/* <img src={PhantomImg} width="20px" height="20px" alt="phantom" /> */}
        <a
          href={`https://solscan.io/tx/${row.txId}?cluster=mainnet`}
          rel="noreferrer"
          target={'_blank'}
        >
          <Span ml="10px" color={'#fff'}>
            {row.txId.slice(0, 4)}...
            {row.txId.slice(-4)}
          </Span>
        </a>
      </Flex>
    ),
    minWidth: '160px !important',
    grow: 1
  },
  {
    name: 'Status',
    cell: row => {
      if (row.status === 'success')
        return <Span color={'#4FFF8B'}>Success</Span>;
      else if (row.status === 'pending')
        return <Span color={'#FFCD4B'}>Pending</Span>;
      else return <Span color={'#F24822'}>Failed</Span>;
    }
  }
];

const customStyles: TableStyles = {
  table: {
    style: {
      fontFamily: 'Inter',

      backgroundColor: 'transparent'
    }
  },
  headRow: {
    style: {
      fontFamily: 'Inter',
      fontSize: '16px',
      border: 'none',
      color: '#768BAD',
      backgroundColor: 'transparent'
    }
  },

  rows: {
    style: {
      border: 'none !important',
      //   borderBottomColor: "transparent !important",
      fontSize: '14px',
      color: '#FFF',
      backgroundColor: 'transparent'
      // padding: '0px 5px !important'
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
  },
  progress: {
    style: {
      backgroundColor: 'transparent',
      color: '#697E9C'
    }
  },

  headCells: {
    style: {
      padding: '0px'
    }
  },
  cells: {
    style: {
      padding: '0px'
    }
  }
};

const temp_data = {
  count: 1,
  history: [
    {
      time: '2022-12-26T05:23:56.145719Z',
      type: 'deposit_nft',
      status: 'success',
      solDetail: {
        solAmount: 0,
        usdAmount: 0
      },
      nftDetail: [
        {
          name: 'DUELBOTS #903',
          mintAddress: 'DYQagijm1jnEy4VCCShsYf3rMK5S5yD1WLRmWC4T9fNx',
          image:
            'https://arweave.net/OHNRSyYALXJqZseh_vzyop_pQxEUI3Ou2705aYqN9mQ',
          collectionName: 'Duelbots',
          collectionImage:
            'https://arweave.net/0jc-_IcbWwp8avEJNwD_oL3PYoQiAS3VY_lKCW8dGlo?ext=png',
          price: 118418
        }
      ],
      txId: '3PjBvi7XUS3Urejs7Uop55a98J9NMj4bzzCR9UvjFrQJYNdw4S2JzJp8W9GjQQLjCQTFPpC9aVbgB244m7HikDjG'
    },
    {
      time: '2023-01-07T07:57:21.262723Z',
      type: 'deposit_bonk',
      status: 'success',
      solDetail: {
        solAmount: 55363675000,
        usdAmount: 106
      },
      nftDetail: [],
      txId: '61smmJFtFn15jXLG7wjPkzmpJAUtZn2o2GNLfTrK2EeVAkgYsjC5JmLjkGvKED51VCummC13e28iVJZow59rW5bq'
    },
    {
      time: '2023-01-07T06:53:23.429736Z',
      type: 'withdraw_bonk',
      status: 'success',
      solDetail: {
        solAmount: 55363675000,
        usdAmount: 100
      },
      nftDetail: [],
      txId: '63PsXoo4uD8LeUyYy3pVCWjLBab2Nyb7cw8FWLBtoPEpknsUgZevbc1jfiU5NwQvLPMnS19Vf1FW6GETdnHjdoce'
    }
  ],

  offset: 0,
  total: 1
};

const StyledWrapFlex = styled(Flex)`
  width: calc(100vw - 60px);
  overflow: auto;
  ${({ theme }) => theme.mediaQueries.md} {
    width: auto;
  }
`;

const StyledFlex = styled(Flex)`
  margin-bottom: 20px;
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
  }
`;

const History: FC = () => {
  const [pageIndex, setPageIndex] = useState(1);
  const [, onDismiss] = useModal(<></>, false);
  const [filter, setFilter] = useState(0);
  console.log(filter);

  // The API URL includes the page index, which is a React state.
  // filter = 1 depoist, 2 withdraw
  // `/pay/history?offset=${(pageIndex - 1) * 5}&count=${5}&filter=${filter}`,
  const { data } = useSWR(
    `/pay/history?offset=${0}&count=${pageIndex * 5}&filter=${filter}`,
    async arg => api.get(arg).then(res => res.data)
  );

  const handlePageChange = (page: number) => {
    setPageIndex(page);
  };

  const getOptions = useCallback(() => {
    let options: { label: string; value: string }[] = [
      { label: 'All Transactions', value: 'all' },
      { label: 'Deposit', value: 'deposit' },
      { label: 'Withdraw', value: 'withdraw' }
    ];

    return options;
  }, []);

  const handleShowMore = () => {
    if (!data) return;
    if (pageIndex < Math.ceil(data.total / 5)) setPageIndex(pageIndex + 1);
    return;
  };

  return (
    <div className="container">
      <div className="box">
        <StyledFlex justifyContent={'space-between'}>
          <Text
            fontSize={'20px'}
            fontWeight={600}
            color="white"
            letterSpacing={'0.18em'}
          >
            WALLET HISTORY
          </Text>

          <Flex gap={30} alignItems="center">
            <Select
              background="#192637"
              hoverBackground="#03060933"
              color="#B2D1FF"
              options={getOptions()}
              onChange={(selectedOption: any) => {
                if (selectedOption.value === 'all') setFilter(0);
                else if (selectedOption.value === 'deposit') setFilter(1);
                else setFilter(2);
              }}
              // components={{
              //   Menu: CustomMenu,
              // }}
            />

            <CloseIcon color="#96A8C2" onClick={onDismiss} cursor="pointer" />
          </Flex>
        </StyledFlex>

        <Box>
          <StyledWrapFlex flexDirection={'column-reverse'}>
            <DataTable
              columns={columns}
              data={data ? data.history : []}
              // pagination
              progressPending={data === undefined}
              // paginationServer
              paginationTotalRows={data && data.total}
              onChangePage={handlePageChange}
              customStyles={customStyles}
              // noTableHead
              // paginationComponent={Pagination}
              // paginationPerPage={5}
              style={{ minWidth: '600px' }}
            />
          </StyledWrapFlex>

          <Button
            variant="secondary"
            outlined
            scale="sm"
            width={153}
            background="linear-gradient(180deg, #070B10 0%, rgba(7, 11, 16, 0.3) 100%)"
            color="#FFFFFF"
            borderColor="chipSecondary"
            marginX="auto"
            marginTop={20}
            onClick={data === undefined ? undefined : handleShowMore}
          >
            {data === undefined ? (
              <ClipLoader size={20} color="#fff" />
            ) : (
              'SHOW MORE'
            )}
          </Button>
        </Box>
      </div>
    </div>
  );
};

export default History;
