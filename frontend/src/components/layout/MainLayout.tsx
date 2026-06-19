import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';
import './MainLayout.css';

const navItems = [
  { path: '/tracking',   label: 'Tracking',    icon: '📦' },
  { path: '/orders',     label: 'Orders',       icon: '🧾' },
  { path: '/dispatch',   label: 'Dispatch',     icon: '🚚' },
  { path: '/warehouse',  label: 'Warehouse',    icon: '🏭' },
  { path: '/pricing',    label: 'Pricing',      icon: '💰' },
  { path: '/settlement', label: 'Settlement',   icon: '💳' },
  { path: '/epod',       label: 'ePOD',         icon: '✍️' },
];

export default function MainLayout() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="layout">
      {/* Sidebar */}
      <aside className="sidebar">
        <div className="sidebar-brand">
          <span className="brand-logo">🚀</span>
          <span className="brand-name">Buroqet</span>
        </div>

        <nav className="sidebar-nav">
          {navItems
            .filter((item) => {
              if (!user) return false;
              if (user.role === 'admin') return true; // admin sees all
              if (user.role === 'pelanggan') {
                return ['/tracking', '/orders', '/pricing'].includes(item.path);
              }
              if (user.role === 'kurir') {
                return ['/tracking', '/epod', '/settlement'].includes(item.path);
              }
              return false;
            })
            .map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
              >
                <span className="nav-icon">{item.icon}</span>
                <span className="nav-label">{item.label}</span>
              </NavLink>
          ))}
        </nav>

        <div className="sidebar-footer">
          <div className="user-info">
            <div className="user-avatar">{user?.name?.[0]?.toUpperCase() ?? 'U'}</div>
            <div className="user-details">
              <span className="user-name">{user?.name ?? 'User'}</span>
              <span className="user-role">{user?.role ?? '-'}</span>
            </div>
          </div>
          <button className="logout-btn" onClick={handleLogout} title="Logout">
            ⏏
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
