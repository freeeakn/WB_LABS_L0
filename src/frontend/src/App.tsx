import { useState, useMemo } from "react";
import { useOrderQuery } from "./queries/useOrderQuery";
import { OrderSearch } from "./components/ui/OrderSearch";
import { LoadingSpinner } from "./components/ui/LoadingSpinner";
import { ErrorAlert } from "./components/ui/ErrorAlert";
import { OrderCard } from "./components/ui/OrderCard";
import { ExampleUid } from "./components/ui/ExampleUid";
import { useAutoAnimate } from "@formkit/auto-animate/react";

export default function App() {
  const [orderUid, setOrderUid] = useState<string>("");

  const { data: order, error, isLoading } = useOrderQuery(orderUid);
  const [ref] = useAutoAnimate<HTMLDivElement>();

  const orderCard = useMemo(() => {
    return order ? <OrderCard key={order.order_uid} order={order} /> : null;
  }, [order]);

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="w-full max-w-4xl">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-primary">Order Lookup</h1>
          <p className="text-base-content/70 mt-2">Введите UID заказа</p>
        </div>

        <OrderSearch onSearch={setOrderUid} />

        <div ref={ref} className="space-y-4">
          {isLoading && <LoadingSpinner />}
          {error && <ErrorAlert />}
          {orderCard}
          {!order && !isLoading && !error && <ExampleUid />}
        </div>
      </div>
    </div>
  );
}
