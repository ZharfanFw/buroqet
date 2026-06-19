import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

// Layouts
import PublicLayout from './components/layout/PublicLayout';
import AuthLayout from './components/layout/AuthLayout';
import CustomerLayout from './components/layout/CustomerLayout';
import CourierLayout from './components/layout/CourierLayout';
import ProtectedRoute from './routes/ProtectedRoute';

// Public Pages
import LandingPage from './pages/public/LandingPage';
import TrackingPage from './pages/public/TrackingPage';
import PricingPage from './pages/public/PricingPage';
import SiteIndexPage from './pages/public/SiteIndexPage';

// Auth Pages
import CustomerLogin from './pages/auth/CustomerLogin';
import CustomerRegister from './pages/auth/CustomerRegister';
import CourierLogin from './pages/auth/CourierLogin';
import CourierRegister from './pages/auth/CourierRegister';

// Customer Pages
import DashboardCustomer from './pages/customer/DashboardCustomer';
import CustomerOrderList from './pages/customer/CustomerOrderList';
import CreateOrderPage from './pages/customer/CreateOrderPage';
import ProfilePage from './pages/customer/ProfilePage';
import DashboardCourier from './pages/courier/DashboardCourier';

// Other Service Pages (Mocked/Reused for now)
import OrderPage from './pages/order/OrderPage';
import DispatchPage from './pages/dispatch/DispatchPage';
import SettlementPage from './pages/settlement/SettlementPage';
import EpodPage from './pages/epod/EpodPage';
import WarehousePage from './pages/warehouse/WarehousePage';

export default function App() {
  return (
    <Router>
      <Routes>
        {/* PUBLIC DOMAIN */}
        <Route element={<PublicLayout />}>
          <Route path="/" element={<LandingPage />} />
          <Route path="/tracking" element={<TrackingPage />} />
          <Route path="/pricing" element={<PricingPage />} />
          <Route path="/services" element={<div className="p-8 text-center">Services Page (Coming Soon)</div>} />
          <Route path="/about" element={<div className="p-8 text-center">About Us (Coming Soon)</div>} />
          <Route path="/contact" element={<div className="p-8 text-center">Contact (Coming Soon)</div>} />
          <Route path="/sitemap" element={<SiteIndexPage />} />
        </Route>

        {/* AUTH DOMAIN */}
        <Route element={<AuthLayout />}>
          <Route path="/customer/login" element={<CustomerLogin />} />
          <Route path="/customer/register" element={<CustomerRegister />} />
          <Route path="/courier/login" element={<CourierLogin />} />
          <Route path="/courier/register" element={<CourierRegister />} />
        </Route>

        {/* CUSTOMER DOMAIN */}
        <Route element={<ProtectedRoute allowedRoles={['pelanggan']}><CustomerLayout /></ProtectedRoute>}>
          <Route path="/customer/dashboard" element={<DashboardCustomer />} />
          <Route path="/customer/orders" element={<CustomerOrderList />} />
          <Route path="/customer/orders/create" element={<CreateOrderPage />} />
          <Route path="/customer/profile" element={<ProfilePage />} />
        </Route>

        {/* COURIER DOMAIN */}
        <Route element={<ProtectedRoute allowedRoles={['kurir']}><CourierLayout /></ProtectedRoute>}>
          <Route path="/courier/dashboard" element={<DashboardCourier />} />
          <Route path="/courier/assignment" element={<div className="p-8 text-center">Assignment Pickup/Delivery</div>} />
          <Route path="/courier/epod" element={<EpodPage />} />
          <Route path="/courier/history" element={<div className="p-8 text-center">Riwayat Pengiriman</div>} />
          <Route path="/courier/profile" element={<div className="p-8 text-center">Profil Courier</div>} />
        </Route>

        </Route>

        {/* CATCH ALL */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}
