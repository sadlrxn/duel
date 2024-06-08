import useSWR from 'swr';

import { api } from 'services';

type Method = 'get' | 'post';

interface CustomSWRProps {
  key: string;
  route: string;
  params?: any;
  method?: Method;
  onError?: any;
}

const fetchCustomData = async (
  route: string,
  method: Method,
  params?: any
  // onError?: any
) => {
  const { data } = await api.request({
    method,
    url: route,
    params
  });
  return data;
};

export default function useCustomSWR({
  key,
  route,
  params,
  method = 'get'
}: // onError
CustomSWRProps) {
  return useSWR(key, () =>
    fetchCustomData(route, method, params /*, onError*/)
  );
}
