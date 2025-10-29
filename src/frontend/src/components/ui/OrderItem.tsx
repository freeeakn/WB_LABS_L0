import { memo } from "react";
import type { Order } from "../../types";

type OrderItemProps = {
  item: Order["items"][0];
};

export const OrderItem = memo(({ item }: OrderItemProps) => {
  return (
    <tr>
      <td>{item.name}</td>
      <td>{item.brand}</td>
      <td>{item.size}</td>
      <td>{item.price}</td>
      <td>{item.sale}%</td>
      <td className="font-semibold">{item.total_price}</td>
    </tr>
  );
});

OrderItem.displayName = "OrderItem";