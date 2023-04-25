import React from 'react';
import { RouterProvider } from 'react-router-dom';
import { router } from '@app/router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // https://tanstack.com/query/v4/docs/react/guides/window-focus-refetching
      refetchOnWindowFocus: false,
      // https://tanstack.com/query/v4/docs/react/guides/query-retries
      retry: false,
    },
  },
});

export default function App() {
  return <QueryClientProvider client={queryClient}>
    <RouterProvider router={router} />
  </QueryClientProvider>
  ;
}
