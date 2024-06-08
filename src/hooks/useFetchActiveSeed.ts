import useSWR from 'swr';
import api from 'utils/api';

const fetchActiveSeed = async () => {
  const { data } = await api.get(`/seed/get-active-seed`);
  return data;
};

export default function useFetchActiveSeed() {
  return useSWR(`Active Seed`, () => fetchActiveSeed());
}
