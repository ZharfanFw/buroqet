// BASE_URL is only the base (e.g. "/api"), do NOT include service prefix here.
// The service prefix (e.g. "/warehouse", "/settlement") is added per-endpoint below.
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
    base: `${BASE_URL}/dispatch`,
    fleet: `${BASE_URL}/dispatch/fleet`,
    assign: (id: string) => `${BASE_URL}/dispatch/${id}/assign`,
  },
  pricing: {
    base: `${BASE_URL}/pricing`,
    calculate: `${BASE_URL}/pricing/calculate`,
  },
  // Settlement: requests go to /api/settlement/v1/...
  // Vite proxy catches /api/settlement → strips to /api → forwards to port 8085
  // Backend (port 8085) receives /api/v1/... ✓
  settlement: {
    base: `${BASE_URL}/settlement/v1/commissions`,
    list: `${BASE_URL}/settlement/v1/commissions`,
    earnings: (id: string) => `${BASE_URL}/settlement/v1/couriers/${id}/earnings`,
  },
  epod: {
    base: `${BASE_URL}/epod`,
    byId: (id: string) => `${BASE_URL}/epod/${id}`,
    upload: `${BASE_URL}/epod/upload`,
  },
  // Warehouse: requests go to /api/warehouse/v1/...
  // Vite proxy catches /api/warehouse → strips to /api → forwards to port 8087
  // Backend (port 8087) receives /api/v1/... ✓
  warehouse: {
    base: `${BASE_URL}/warehouse/v1/warehouse`,
    packages: `${BASE_URL}/warehouse/v1/packages`,
    inbound: `${BASE_URL}/warehouse/v1/inbound`,
    dispatch: `${BASE_URL}/warehouse/v1/dispatch`,
  },
} as const;

export default BASE_URL;
