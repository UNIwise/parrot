import {
  QueryClient,
  QueryClientProvider
} from "@tanstack/react-query";
import { FC, ReactNode } from "react";

// Create react-query client
const queryClient = new QueryClient();

const ReactQueryClientProvider: FC<{ children: ReactNode }> = ({
  children,
}) => {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export { ReactQueryClientProvider };
