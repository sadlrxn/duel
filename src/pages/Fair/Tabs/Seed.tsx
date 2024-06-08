import React, { useState, useCallback, useEffect } from 'react';
import styled from 'styled-components';
import { toast } from 'react-toastify';
import { ClipLoader } from 'react-spinners';

import { Box, Text, Flex, Button, Span } from 'components';
import { useFetchActiveSeed } from 'hooks';
import { api } from 'services';

import Detail, { TableItem } from './Detail';

interface History {
  client: string;
  server: string;
  serverHash: string;
  nonce: number;
}

interface SeedProps {
  seedHash?: string;
}

const Seed: React.FC<SeedProps> = () => {
  const { data, error } = useFetchActiveSeed();

  const [seed, setSeed] = useState('');
  const [hash, setHash] = useState('');
  const [newSeed, setNewSeed] = useState('');
  const [isUnhashing, setIsUnhashing] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isFetchHistory, setIsFetchHistory] = useState(false);
  const [history, setHistory] = useState<History[]>([]);
  const [activeSeed, setActiveSeed] = useState({
    client: '',
    server: '',
    bets: 0
  });

  useEffect(() => {
    let seed = {
      client: '',
      server: '',
      bets: 0
    };
    if (!error && data) {
      const { clientSeed: client, serverSeedHash: server, nonce: bets } = data;
      seed = { client: client ?? '', server: server ?? '', bets: bets ?? 0 };
    }
    setActiveSeed(seed);
  }, [error, data]);

  const handleFetchHistory = useCallback(async () => {
    setIsFetchHistory(true);
    try {
      const { data } = await api.get('/seed/history', {
        params: {
          offset: history.length,
          count: 3
        }
      });

      const newHistory = data
        //@ts-ignore
        .map(s => {
          const { clientSeed, serverSeed, serverSeedHash, nonce } = s;
          return {
            client: clientSeed ?? '',
            server: serverSeed ?? '',
            serverHash: serverSeedHash ?? '',
            nonce: nonce ?? 0
          };
        })
        .filter((h: History) => {
          const his = history.find(item => item.server === h.server);
          if (his) return false;
          return true;
        });
      setHistory([...history, ...newHistory]);
    } catch {}
    setIsFetchHistory(false);
  }, [history]);

  const handleUpdating = useCallback(async () => {
    if (newSeed === '') {
      toast.error('New client seed is empty.');
      return;
    }
    setIsUpdating(true);
    try {
      const { data } = await api.post('/seed/rotate', {
        serverSeedHash: activeSeed.server,
        clientSeed: newSeed
      });
      const { clientSeed: client, serverSeedHash: server, nonce: bets } = data;
      setActiveSeed({
        client: client ?? '',
        server: server ?? '',
        bets: bets ?? 0
      });
      toast.success('The client seed was updated successfully.');
    } catch (err: any) {
      if (err.response.status === 429) {
        toast.error(err.response.data.message);
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      }
    }
    setIsUpdating(false);
  }, [newSeed]);

  const handleClick = useCallback(async () => {
    setIsUnhashing(true);
    try {
      const { data } = await api.get(`/seed/unhash?hash=${hash}`);
      setSeed(data.seed);
    } catch (error) {
      toast.error(
        'The current server seed is still in use, please change your client seed to unhash the server seed.'
      );
    }
    setIsUnhashing(false);
  }, [hash]);

  useEffect(() => {
    if (history.length === 0 && !isFetchHistory) handleFetchHistory();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div>
      <Flex flexDirection="column" gap={20} fontSize="20px" lineHeight="1.25em">
        <Text color="#D0DAEB" fontWeight={600} fontSize="18px">
          ACTIVE SEED PAIR
        </Text>

        <Text color="#B9D2FD" fontSize="14px" mt="-10px" mb="10px">
          Your <Span fontStyle="italic">client seed</Span> is like your lucky
          key, game outcomes depend on your{' '}
          <Span fontStyle="italic">client seed</Span> and{' '}
          <Span fontStyle="italic">server seed</Span> (seed pair). Every time
          you update your <Span fontStyle="italic">client seed</Span> a new{' '}
          <Span fontStyle="italic">server seed</Span> is generated, hashed and
          paired to create a <Span fontStyle="italic">seed pair</Span> - your
          previous <Span fontStyle="italic">client seed</Span> becomes expired
          and the server seed unhashed. Allowing you to validate the outcomes of
          all passed and future games.
        </Text>

        <Detail
          title="New Client Seed"
          buttonText="UPDATE SEED"
          buttonClick={handleUpdating}
          isLoading={isUpdating}
          text={newSeed}
          setText={setNewSeed}
          mb="10px"
        />

        <Detail
          title="Active Client Seed"
          text={activeSeed.client}
          readOnly
          enableCopy
        />

        <Detail
          title="Active Server Seed (Hashed)"
          text={activeSeed.server}
          readOnly
          enableCopy
        />

        <Box maxWidth={200}>
          <Detail
            title="Total Bets with Pair"
            text={activeSeed.bets.toString()}
            readOnly
            enableCopy
          />
        </Box>

        <TableWrapper>
          <Table>
            <thead>
              <tr>
                <th>Client Seed</th>
                <th>Server Seed (Hashed)</th>
                <th>Server Seed (Unhashed)</th>
                <th>Total Bets with Pair</th>
              </tr>
            </thead>
            <tbody>
              {history.map((h, i) => {
                return (
                  <tr key={i}>
                    <td>
                      <TableItem text={h.client} enableCopy />
                    </td>
                    <td>
                      <TableItem text={h.serverHash} enableCopy />
                    </td>
                    <td>
                      <TableItem text={h.server} enableCopy />
                    </td>
                    <td>
                      <TableItem text={h.nonce.toString()} />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </Table>
        </TableWrapper>

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
          onClick={isFetchHistory ? undefined : handleFetchHistory}
        >
          {isFetchHistory ? <ClipLoader color="#fff" size={20} /> : 'SHOW MORE'}
        </Button>

        <Text color="#D0DAEB" fontWeight={600} fontSize="18px" mt="75px">
          UNHASH SERVER SEED
        </Text>

        <Text color="#B9D2FD" fontSize="14px" mt="-10px" mb="10px">
          You can unhash anyone’s server seed as long as their seed pair is
          expired. To unhash your current server seed update your active client
          seed with a new seed.
        </Text>

        <Detail
          title="Server Seed (Hashed)"
          text={hash}
          setText={setHash}
          buttonText="UNHASH"
          isLoading={isUnhashing}
          buttonClick={handleClick}
        />
        <Detail title="Server Seed" text={seed} readOnly enableCopy />
      </Flex>
    </div>
  );
};

export default React.memo(Seed);

const TableWrapper = styled(Box)`
  margin-top: 50px;
  max-height: 600px;
  overflow: auto;
  &::-webkit-scrollbar-track {
    margin-top: 38px;
    margin-bottom: 10px;
  }
  &::-webkit-scrollbar-track-piece {
    margin-right: 10px;
  }
`;

const Table = styled.table`
  position: relative;
  width: 100%;
  min-width: 700px;

  font-size: 16px;
  font-weight: 400;
  line-height: 19px;
  color: #ffffff;

  border-collapse: separate;
  border-spacing: 0;

  padding-bottom: 10px;

  th,
  td {
    text-align: left;
    vertical-align: middle;
    padding: 9px 20px;
    width: 26%;
  }

  th {
    position: sticky;
    top: 0;
    color: #768bad;
    background: #0b141e;
  }

  td {
    background: #121c2a;
  }

  tr td:last-child,
  th:last-child {
    text-align: right;
  }

  tr:first-child {
    td:first-child {
      border-top-left-radius: 10px;
    }
    td:last-child {
      border-top-right-radius: 10px;
    }
  }

  tr:last-child {
    td:first-child {
      border-bottom-left-radius: 10px;
    }
    td:last-child {
      border-bottom-right-radius: 10px;
    }
  }
`;
