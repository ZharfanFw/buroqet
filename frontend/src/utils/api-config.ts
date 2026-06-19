const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

export const API_ENDPOINTS = {
  auth: {
    login: `${BASE_URL}/auth/login`,
    register: `${BASE_URL}/auth/register`,
    logout: `${BASE_URL}/auth/logout`,
    me: `${BASE_URL}/auth/me`,
    refresh: `${BASE_URL}/auth/refresh`,
  },
  tracking: {
    base: `${BASE_URL}/tracking`,
    byAwb: (awb: string) => `${BASE_URL}/tracking/${awb}`,
    events: (awb: string) => `${BASE_URL}/tracking/${awb}/events`,
  },
  order: {
    base: `${BASE_URL}/orders`,
    byId: (id: string) => `${BASE_URL}/orders/${id}`,
    status: (id: string) => `${BASE_URL}/orders/${id}/status`,
  },
  dispatch: {
    base: `${BASE_URL}/dispatch/v1/dispatch`,
    assign: `${BASE_URL}/dispatch/v1/dispatch/assign`,
    startDelivery: `${BASE_URL}/dispatch/v1/dispatch/start-delivery`,
  },
  pricing: {
    base: `${BASE_URL}/pricing`,
    calculate: `${BASE_URL}/pricing/calculate`,
  },
  settlement: {
    base: `${BASE_URL}/settlement/api/v1`,
    commissions: `${BASE_URL}/settlement/api/v1/commissions`,
    earnings: (courierId: string) => `${BASE_URL}/settlement/api/v1/couriers/${courierId}/earnings`,
  },
  epod: {
    base: `${BASE_URL}/epod`,
    byId: (id: string) => `${BASE_URL}/epod/${id}`,
    upload: `${BASE_URL}/epod/upload`,
  },
  warehouse: {
    base: `${BASE_URL}/warehouse/api/v1`,
    packages: `${BASE_URL}/warehouse/api/v1/packages`,
    inbound: `${BASE_URL}/warehouse/api/v1/inbound`,
    dispatch: `${BASE_URL}/warehouse/api/v1/dispatch`,
  },
} as const;

export default BASE_URL;
