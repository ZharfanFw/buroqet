import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';

const navItems = [
  { path: '/customer/dashboard', label: 'Dashboard', icon: '🏠' },
  { path: '/customer/orders/create', label: 'Buat Order', icon: '➕' },
  { path: '/customer/orders', label: 'Order Saya', icon: '🧾' },
  { path: '/customer/profile', label: 'Profil', icon: '👤' },
];

export default function CustomerLayout() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/customer/login');
  };

  return (
    <div className="flex h-screen bg-[#f7faf7] overflow-hidden">
      {/* Sidebar */}
      <aside className="w-72 bg-white border-r border-[#e4ece3] flex flex-col shadow-sm">
        <div className="h-20 flex items-center px-8 border-b border-[#e4ece3]">
          <span className="text-3xl mr-3 text-[#64965a]">🚀</span>
          <span className="text-2xl font-extrabold text-slate-800 tracking-tight">Buroqet</span>
        </div>

        <div className="p-6 pb-2">
          <p className="text-xs font-bold text-slate-400 uppercase tracking-widest mb-4">Customer Area</p>
        </div>

        <nav className="flex-1 px-4 space-y-2 overflow-y-auto custom-scrollbar">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `flex items-center gap-4 px-4 py-3.5 text-sm font-bold rounded-2xl transition-all duration-300 ${
                  isActive
                    ? 'bg-[#64965a] text-white shadow-md transform scale-[1.02]'
                    : 'text-slate-500 hover:bg-[#64965a]/10 hover:text-[#64965a]'
                }`
              }
            >
              <span className="text-xl">{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="p-6 border-t border-[#e4ece3] bg-slate-50/50">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#64965a] to-[#79ae6f] text-white flex items-center justify-center font-extrabold text-lg shadow-md">
              {user?.name?.[0]?.toUpperCase() || 'C'}
            </div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-bold text-slate-800 truncate">{user?.name}</p>
              <p className="text-xs text-slate-500 font-medium truncate">{user?.email}</p>
            </div>
            <button onClick={handleLogout} className="text-slate-400 hover:text-red-500 p-2 hover:bg-red-50 rounded-xl transition-colors" title="Logout">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col h-screen overflow-hidden">
        {/* Topbar */}
        <header className="h-20 bg-white border-b border-[#e4ece3] flex items-center justify-between px-10 z-10 shadow-sm">
          <div className="flex items-center gap-4">
            <h1 className="text-xl font-extrabold text-slate-800">Customer Portal</h1>
            <div className="hidden md:flex items-center gap-2 text-sm text-slate-400 font-medium">
              <span>/</span>
              <span className="text-[#64965a]">Home</span>
            </div>
          </div>

          <div className="flex items-center gap-6">
            <div className="relative hidden md:block">
              <span className="absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
              </span>
              <input type="text" placeholder="Track package..." className="py-2.5 pl-11 pr-4 bg-slate-50 border border-[#e4ece3] rounded-2xl text-sm focus:outline-none focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] transition-shadow w-64" />
            </div>

            <button className="text-slate-400 hover:text-[#64965a] relative p-2 hover:bg-slate-50 rounded-xl transition-colors">
              <span className="absolute top-1 right-1 w-2.5 h-2.5 bg-red-500 rounded-full border-2 border-white"></span>
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path></svg>
            </button>
          </div>
        </header>

        <div className="flex-1 overflow-y-auto p-10 bg-[#f7faf7]">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
