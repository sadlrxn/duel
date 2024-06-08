import useSWR from 'swr';
import api from 'utils/api';

interface fetchSolPriceResponse {
  price: number;
}

const useTokenPrices = () => {
  const { data: prices } = useSWR<fetchSolPriceResponse>(
    `token-prices`,
    async () => api.get(`/token-prices`).then(res => res.data),
    { refreshInterval: 60 * 1000 }
  );

  return {
    prices
  };
};

export default useTokenPrices;
