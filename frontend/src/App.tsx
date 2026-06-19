import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { useAuthStore } from './store/auth.store';

// Layouts
import MainLayout from './components/layout/MainLayout';
import AuthLayout from './components/layout/AuthLayout';

// Pages - Auth
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';

// Pages - Services
import TrackingPage from './pages/tracking/TrackingPage';
import OrderPage from './pages/order/OrderPage';
import DispatchPage from './pages/dispatch/DispatchPage';
import PricingPage from './pages/pricing/PricingPage';
import SettlementPage from './pages/settlement/SettlementPage';
import EpodPage from './pages/epod/EpodPage';
import WarehousePage from './pages/warehouse/WarehousePage';

// Guard: redirect to login if not authenticated
function PrivateRoute({ children, allowedRoles }: { children: React.ReactNode, allowedRoles?: string[] }) {
  const { isAuthenticated, user, isLoading } = useAuthStore();
  
  if (isLoading) return <div>Loading...</div>;
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (allowedRoles && user && !allowedRoles.includes(user.role)) {
    return <Navigate to="/tracking" replace />; // Fallback to a safe route
  }
  
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
          <Route path="/register" element={<RegisterPage />} />
        </Route>

        {/* Public routes — accessible tanpa login */}
        <Route element={<MainLayout />}>
          <Route path="/tracking" element={<TrackingPage />} />
        </Route>

        {/* Protected routes */}
        <Route element={<PrivateRoute><MainLayout /></PrivateRoute>}>
          <Route path="/" element={<Navigate to="/tracking" replace />} />
          
          {/* Accessible by pelanggan & admin */}
          <Route path="/orders" element={<PrivateRoute allowedRoles={['admin', 'pelanggan']}><OrderPage /></PrivateRoute>} />
          <Route path="/pricing" element={<PrivateRoute allowedRoles={['admin', 'pelanggan']}><PricingPage /></PrivateRoute>} />
          
          {/* Accessible by kurir & admin */}
          <Route path="/epod" element={<PrivateRoute allowedRoles={['admin', 'kurir']}><EpodPage /></PrivateRoute>} />
          <Route path="/settlement" element={<PrivateRoute allowedRoles={['admin', 'kurir']}><SettlementPage /></PrivateRoute>} />
          
          {/* Accessible by admin only */}
          <Route path="/dispatch" element={<PrivateRoute allowedRoles={['admin']}><DispatchPage /></PrivateRoute>} />
          <Route path="/warehouse" element={<PrivateRoute allowedRoles={['admin']}><WarehousePage /></PrivateRoute>} />
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
