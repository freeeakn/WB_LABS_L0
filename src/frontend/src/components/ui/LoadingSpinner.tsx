import { memo } from "react";

export const LoadingSpinner = memo(() => {
  return (
    <div className="flex justify-center">
      <span className="loading loading-spinner loading-lg"></span>
    </div>
  );
});

LoadingSpinner.displayName = "LoadingSpinner";