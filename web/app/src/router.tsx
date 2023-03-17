import { createBrowserRouter, generatePath, Params } from 'react-router-dom';
import Layout from '@app/pages/Layout';
import RecipeList from '@app/pages/RecipeList';
import ProjetConfiguration from '@app/pages/ProjetConfiguration';
import React from 'react';
import Error404 from '@app/pages/Error404';

export const Routes = {
  recipes: '/list',
  init: '/init/:recipeName',
};

export const router = createBrowserRouter([
  {
    path: Routes.recipes,
    element: <Layout><RecipeList /></Layout>,
  },
  {
    path: Routes.init,
    element: <Layout><ProjetConfiguration /></Layout>,
  },
  {
    path: '*',
    element: <Error404 />,
  },
]);

/**
 * Gets the path for a page.
 */
function path(route: string, parameters: Params|undefined = undefined): string {
  if (!parameters) {
    return route;
  }

  return generatePath(route, parameters);
}

/**
 * Alias exposed outside to resolve the path of page and use it for links: <Link to={route(Home)} />…</Link>
 */
export const route = path;
