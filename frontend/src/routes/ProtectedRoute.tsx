import { Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '../store/auth.store';

interface ProtectedRouteProps {
  children: React.ReactNode;
  allowedRoles: ('admin' | 'kurir' | 'pelanggan')[];
}

export default function ProtectedRoute({ children, allowedRoles }: ProtectedRouteProps) {
  const { isAuthenticated, user } = useAuthStore();
  const location = useLocation();

  if (!isAuthenticated || !user) {
    // Determine which login page to redirect to based on the requested path
    if (location.pathname.startsWith('/admin')) {
      return <Navigate to="/admin/login" replace />;
    }
    if (location.pathname.startsWith('/courier')) {
      return <Navigate to="/courier/login" replace />;
    }
    return <Navigate to="/customer/login" replace />;
  }

  // Check if user has required role
  if (!allowedRoles.includes(user.role as any)) {
    // Redirect to their respective dashboard if they try to access wrong role area
    if (user.role === 'admin') return <Navigate to="/admin/dashboard" replace />;
    if (user.role === 'kurir') return <Navigate to="/courier/dashboard" replace />;
    return <Navigate to="/customer/dashboard" replace />;
  }

  return <>{children}</>;
}
