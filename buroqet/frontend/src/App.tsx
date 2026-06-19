import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { useAuthStore } from './store/auth.store';

// Layouts
import MainLayout from './components/layout/MainLayout';
import AuthLayout from './components/layout/AuthLayout';

// Pages - Auth
import LoginPage from './pages/auth/LoginPage';

// Pages - Services
import TrackingPage from './pages/tracking/TrackingPage';
import OrderPage from './pages/order/OrderPage';
import DispatchPage from './pages/dispatch/DispatchPage';
import PricingPage from './pages/pricing/PricingPage';
import SettlementPage from './pages/settlement/SettlementPage';
import EpodPage from './pages/epod/EpodPage';
import WarehousePage from './pages/warehouse/WarehousePage';

// Guard: redirect to login if not authenticated
function PrivateRoute({ children }: { children: React.ReactNode }) {
  // const { isAuthenticated } = useAuthStore();
  // return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
  
  // Bypassed for UI testing
  return <>{children}</>;
}

function App() {
  const { fetchMe, isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (isAuthenticated) {
      fetchMe();
    }
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        {/* Auth routes */}
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<LoginPage />} />
        </Route>

        {/* Public routes — accessible tanpa login */}
        <Route element={<MainLayout />}>
          <Route path="/tracking" element={<TrackingPage />} />
        </Route>

        {/* Protected routes */}
        <Route
          element={
            <PrivateRoute>
              <MainLayout />
            </PrivateRoute>
          }
        >
          <Route path="/" element={<Navigate to="/tracking" replace />} />
          <Route path="/orders" element={<OrderPage />} />
          <Route path="/dispatch" element={<DispatchPage />} />
          <Route path="/pricing" element={<PricingPage />} />
          <Route path="/settlement" element={<SettlementPage />} />
          <Route path="/epod" element={<EpodPage />} />
          <Route path="/warehouse" element={<WarehousePage />} />
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
