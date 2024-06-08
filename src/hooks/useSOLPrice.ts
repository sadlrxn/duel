import useSWR from "swr";
import api from "utils/api";

interface fetchSolPriceResponse {
  price: number;
}

const useSOLPrice = () => {
  const { data: solPrice } = useSWR<fetchSolPriceResponse>(
    ["fetch-sol-price"],
    async () => api.get("/sol-price").then((res) => res.data),
    { refreshInterval: 10000 }
  );

  return {
    solPrice: solPrice?.price,
  };
};

export default useSOLPrice;
