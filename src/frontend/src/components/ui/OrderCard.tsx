import { memo } from "react";
import { OrderItem } from "./OrderItem";
import type { Order } from "../../types";
import { useAutoAnimate } from "@formkit/auto-animate/react";

export const OrderCard = memo(({ order }: { order: Order }) => {
  const [parent] = useAutoAnimate()
  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <div className="flex justify-between items-start mb-4">
          <h2 className="card-title text-2xl">
            Заказ: <span className="font-mono">{order.order_uid}</span>
          </h2>
          <div className="badge badge-success">Активен</div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h3 className="font-semibold text-lg mb-2">Доставка</h3>
            <div className="space-y-1 text-sm">
              <p><strong>Получатель:</strong> {order.delivery.name}</p>
              <p><strong>Телефон:</strong> {order.delivery.phone}</p>
              <p><strong>Адрес:</strong> {`${order.delivery.zip}, ${order.delivery.city}, ${order.delivery.address}`}</p>
              <p><strong>Email:</strong> {order.delivery.email}</p>
            </div>
          </div>

          <div>
            <h3 className="font-semibold text-lg mb-2">Оплата</h3>
            <div className="space-y-1 text-sm">
              <p><strong>Сумма:</strong> {order.payment.amount} {order.payment.currency}</p>
              <p><strong>Доставка:</strong> {order.payment.delivery_cost} {order.payment.currency}</p>
              <p><strong>Товары:</strong> {order.payment.goods_total} {order.payment.currency}</p>
              <p><strong>Провайдер:</strong> {order.payment.provider}</p>
            </div>
          </div>
        </div>

        <div className="mt-6">
          <h3 className="font-semibold text-lg mb-3">Товары</h3>
          <div className="overflow-x-auto">
            <table className="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Название</th>
                  <th>Бренд</th>
                  <th>Размер</th>
                  <th>Цена</th>
                  <th>Скидка</th>
                  <th>Итого</th>
                </tr>
              </thead>
              <tbody ref={parent}>
                {order.items.map((item) => (
                  <OrderItem key={item.chrt_id} item={item} />
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className="mt-6 text-sm text-base-content/60">
          <p>
            <strong>Трек-номер:</strong> {order.track_number} |{" "}
            <strong>Создан:</strong> {new Date(order.date_created).toLocaleString()}
          </p>
        </div>
      </div>
    </div>
  );
});

OrderCard.displayName = "OrderCard";