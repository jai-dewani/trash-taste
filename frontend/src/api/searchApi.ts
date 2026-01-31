import axios from 'axios';
import type { SearchResponse } from '../types';

// When using Vite proxy, we don't need a base URL
// The proxy will forward /api requests to the backend
const API_BASE_URL = import.meta.env.VITE_API_URL || '';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

export const searchApi = {
  search: async (query: string, limit: number = 50): Promise<SearchResponse> => {
    const response = await apiClient.get<SearchResponse>('/api/search', {
      params: { q: query, limit },
    });
    return response.data;
  },
};
