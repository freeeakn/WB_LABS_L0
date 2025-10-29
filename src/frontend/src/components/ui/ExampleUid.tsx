import { memo } from "react";

export const ExampleUid = memo(() => {
  return (
    <div className="text-center text-base-content/50 mt-8">
      <p className="text-lg">Пример UID:</p>
      <code className="bg-base-300 px-3 py-1 rounded">
        b563feb7b2b84b6test
      </code>
    </div>
  );
});

ExampleUid.displayName = "ExampleUid";