import { PageSpinner, Flex, Box } from 'components';
import { useAppSelector } from 'state';

import { CreateGame, List } from './components';
import TopBar from './components/Topbar';

export default function Coinflip() {
  const loading = useAppSelector(state => !state.coinflip.fetch);

  if (loading) return <PageSpinner />;

  return (
    <Box padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}>
      <TopBar />
      <Flex flexDirection={'column'} gap={70} mt="34px">
        <CreateGame />
        <List />
      </Flex>
    </Box>
  );
}
