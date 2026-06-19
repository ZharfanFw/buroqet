import { Outlet } from 'react-router-dom';
import './AuthLayout.css';

export default function AuthLayout() {
  return (
    <div className="auth-layout">
      <div className="auth-card">
        <div className="auth-brand">
          <span className="auth-logo">🚀</span>
          <h1 className="auth-title">Buroqet</h1>
          <p className="auth-subtitle">Logistics Management System</p>
        </div>
        <Outlet />
      </div>
    </div>
  );
}
