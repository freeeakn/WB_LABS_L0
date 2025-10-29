import { memo, useCallback, useState } from "react";

interface Props {
  onSearch: (uid: string) => void;
}

export const OrderSearch = memo(({ onSearch }: Props) => {
  const [input, setInput] = useState("");

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      const uid = input.trim();
      if (uid) {
        onSearch(uid);
        setInput("");
      }
    },
    [input, onSearch]
  );

  return (
    <form onSubmit={handleSubmit} className="flex gap-2 mb-8">
      <input
        type="text"
        placeholder="b563feb7b2b84b6test"
        className="input input-bordered input-lg flex-1"
        value={input}
        onChange={(e) => setInput(e.target.value)}
      />
      <button type="submit" className="btn btn-primary btn-lg">
        Найти
      </button>
    </form>
  );
});

OrderSearch.displayName = "OrderSearch";