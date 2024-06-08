import { DependencyList, EffectCallback, useEffect } from "react";
import useSWR from "swr";
import { FAST_INTERVAL, SLOW_INTERVAL } from "config/constants";

export function useFastRefreshEffect(
  effect: EffectCallback,
  deps?: DependencyList
) {
  const { data = 0 } = useSWR([FAST_INTERVAL]);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(effect.bind(null, data), [data, ...(deps || [])]);
}

export function useSlowRefreshEffect(
  effect: EffectCallback,
  deps?: DependencyList
) {
  const { data = 0 } = useSWR([SLOW_INTERVAL]);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(effect.bind(null, data), [data, ...(deps || [])]);
}
