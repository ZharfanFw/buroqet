import { create } from 'zustand';
import apiClient from '../services/api-client';
import { API_ENDPOINTS } from '../utils/api-config';

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string, role: string) => Promise<void>;
  logout: () => void;
  fetchMe: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: !!localStorage.getItem('access_token'),
  isLoading: false,
  error: null,

  login: async (email, password) => {
    set({ isLoading: true, error: null });
    try {
      const { data } = await apiClient.post(API_ENDPOINTS.auth.login, { email, password });
      localStorage.setItem('access_token', data.access_token);
      if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token);
      set({ user: data.user, isAuthenticated: true, isLoading: false });
    } catch (err: unknown) {
      const msg = (err as any)?.response?.data?.error || (err as any)?.response?.data?.message || 'Login gagal';
      set({ error: msg, isLoading: false });
    }
  },

  register: async (name, email, password, role) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.post(API_ENDPOINTS.auth.register, { name, email, password, role });
      set({ isLoading: false });
    } catch (err: unknown) {
      const msg = (err as any)?.response?.data?.error || (err as any)?.response?.data?.message || 'Registrasi gagal';
      set({ error: msg, isLoading: false });
      throw new Error(msg); // throw error agar component bisa tangkap
    }
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    set({ user: null, isAuthenticated: false });
  },

  fetchMe: async () => {
    set({ isLoading: true });
    try {
      const { data } = await apiClient.get(API_ENDPOINTS.auth.me);
      set({ user: data, isAuthenticated: true, isLoading: false });
    } catch {
      set({ isLoading: false, isAuthenticated: false });
    }
  },
}));
