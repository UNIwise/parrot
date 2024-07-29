import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { AxiosError } from "axios";
import { StatusCodes } from "http-status-codes";
import { FC, ReactNode } from "react";

interface GeneralAPIResponse {
  error?: {
    code: number;
    message: string;
    errors: { reason: string; message: string }[];
  };
}

// Create react-query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 15, // 15 minutes
      gcTime: 1000 * 60 * 15, // 15 minutes
      notifyOnChangeProps: ["data", "isLoading", "error"],
      retry: (retryCount, err) => {
        // Get status from error object
        let status: number | undefined = undefined;

        // General API error setup
        if (Object.prototype.hasOwnProperty.call(err, "code")) {
          status = (err as GeneralAPIResponse).error?.code;
        }

        // Normal Axios error setup
        if (Object.prototype.hasOwnProperty.call(err, "response")) {
          status = (err as AxiosError).response?.status;
        }

        // Do not retry if status is Unauthorized or Forbidden
        if (
          status !== undefined &&
          [StatusCodes.UNAUTHORIZED, StatusCodes.FORBIDDEN].includes(status)
        ) {
          return false;
        }

        // Default 3 tries (initial + 2 retries)
        return retryCount < 2;
      },
    },
  },
  queryCache: new QueryCache({
    onError: (error) => console.error(`Something went wrong: ${error.message}`),
  }),
});

const ReactQueryClientProvider: FC<{ children: ReactNode }> = ({
  children,
}) => {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export { ReactQueryClientProvider };
