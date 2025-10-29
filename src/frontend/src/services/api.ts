import axios from "axios";
import type { Order } from "../types";


export const orderApi = {
  getByUid: async (uid: string): Promise<Order> => {
    const response = await axios.get(`/order/${uid}`);
    return response.data;
  },
};