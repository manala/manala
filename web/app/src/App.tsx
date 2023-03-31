import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Layout from '@app/pages/Layout';
import Error404 from '@app/pages/Error404';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout>todo homepage</Layout>,
  },
  {
    path: '*',
    element: <Error404 />,
  },
]);

export default function App() {
  return <RouterProvider router={router} />;
}
