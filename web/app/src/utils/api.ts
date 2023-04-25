import { ltrim, rtrim } from '@app/utils/strings';
import AppConfig from '@app/AppConfig';

/**
 * Composes the URI of an API resource with base API URI.
 */
export function apiUri(path: string): string {
  return [
    rtrim(AppConfig.API_BASE_URI, '/'),
    ltrim(path, '/'),
  ].join('/');
}

/**
 * Returns a function that fetches a JSON API resource.
 */
export function jsonApiQuery(path: string) {
  return () => fetch(apiUri(path)).then((res) => res.json());
}
