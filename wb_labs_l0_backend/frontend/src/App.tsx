import { useState } from "react";
import type { Order } from "./types";

function App() {
  const [id, setId] = useState<string>("");
  const [order, setOrder] = useState<Order | null>(null);
  const [err, setErr] = useState<string>("");

  const fetchOrder = async (): Promise<void> => {
    setErr("");
    setOrder(null);

    if (!id) {
      setErr("Введите id заказа");
      return;
    }

    try {
      const res = await fetch(`/orders/${encodeURIComponent(id)}`);
      if (!res.ok) {
        const txt = await res.text();
        throw new Error(`${res.status} ${txt}`);
      }

      const data: Order = await res.json();
      setOrder(data);
    } catch (e) {
      const message = e instanceof Error ? e.message : String(e);
      setErr(message);
    }
  };

  return (
    <div style={{ padding: 20, fontFamily: "Arial, sans-serif" }}>
      <h2>Order viewer</h2>
      <div style={{ marginBottom: 10 }}>
        <input
          value={id}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setId(e.target.value)}
          placeholder="order_uid"
          style={{ padding: 8, width: 320 }}
        />
        <button onClick={fetchOrder} style={{ marginLeft: 8, padding: 8 }}>
          Fetch
        </button>
      </div>
      {err && <div style={{ color: "red" }}>{err}</div>}
      {order && (
        <pre style={{ background: "#f6f6f6", padding: 12 }}>
          {JSON.stringify(order, null, 2)}
        </pre>
      )}
    </div>
  );
}

export default App;
