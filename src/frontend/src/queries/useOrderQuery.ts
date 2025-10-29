import { useQuery } from "@tanstack/react-query";
import { orderApi } from "../services/api";
import type { Order } from "../types";

export const useOrderQuery = (orderUid: string) => {
  return useQuery<Order, Error>({
    queryKey: ["order", orderUid],
    queryFn: () => orderApi.getByUid(orderUid),
    enabled: !!orderUid,
    retry: 1,
    staleTime: 5 * 60 * 1000,
  });
};