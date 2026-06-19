import { Outlet, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';

export default function PublicLayout() {
  const { user } = useAuthStore();

  return (
    <div className="min-h-screen flex flex-col bg-[#f7faf7]">
      {/* Navbar */}
      <nav className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-[#e4ece3] transition-all">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            <div className="flex items-center gap-2">
              <span className="text-3xl text-[#64965a]">🚀</span>
              <Link to="/" className="text-2xl font-extrabold text-slate-800 tracking-tight">Buroqet</Link>
            </div>

            <div className="hidden md:flex items-center space-x-8">
              <Link to="/" className="text-sm font-bold text-slate-500 hover:text-[#64965a] transition-colors">Home</Link>
              <a href="/#services" className="text-sm font-bold text-slate-500 hover:text-[#64965a] transition-colors">Layanan</a>
              <Link to="/tracking" className="text-sm font-bold text-slate-500 hover:text-[#64965a] transition-colors">Cek Resi</Link>
              <Link to="/pricing" className="text-sm font-bold text-slate-500 hover:text-[#64965a] transition-colors">Cek Ongkir</Link>
            </div>

            <div className="flex items-center space-x-4">
              {user ? (
                <Link
                  to={user.role === 'pelanggan' ? '/customer/dashboard' : user.role === 'admin' ? '/admin/dashboard' : '/courier/dashboard'}
                  className="px-5 py-2.5 bg-[#64965a]/10 text-[#64965a] rounded-xl text-sm font-bold hover:bg-[#64965a]/20 transition-colors"
                >
                  Dashboard
                </Link>
              ) : (
                <>
                  <Link to="/customer/login" className="text-sm font-bold text-[#64965a] hover:text-[#53804a] transition-colors">Masuk</Link>
                  <Link to="/customer/login" className="px-5 py-2.5 bg-[#64965a] text-white text-sm font-bold rounded-xl hover:bg-[#53804a] shadow-md hover:shadow-lg hover:-translate-y-0.5 transition-all">
                    Kirim Paket
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="flex-1">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-[#e4ece3] pt-16 pb-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-12">
            <div className="col-span-1 md:col-span-1">
              <div className="flex items-center gap-2 mb-6">
                <span className="text-3xl text-[#64965a]">🚀</span>
                <span className="text-2xl font-extrabold text-slate-800 tracking-tight">Buroqet</span>
              </div>
              <p className="text-slate-500 font-medium mb-6">Sistem Logistik Masa Depan untuk Indonesia yang lebih terhubung.</p>
            </div>

            <div>
              <h4 className="text-slate-800 font-extrabold mb-6">Perusahaan</h4>
              <ul className="space-y-4">
                <li><Link to="/about" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Tentang Kami</Link></li>
                <li><Link to="/contact" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Hubungi Kami</Link></li>
                <li><a href="#" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Karir</a></li>
              </ul>
            </div>

            <div>
              <h4 className="text-slate-800 font-extrabold mb-6">Layanan</h4>
              <ul className="space-y-4">
                <li><Link to="/pricing" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Cek Tarif</Link></li>
                <li><Link to="/tracking" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Lacak Paket</Link></li>
                <li><a href="#" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Buroqet API</a></li>
              </ul>
            </div>

            <div>
              <h4 className="text-slate-800 font-extrabold mb-6">Bantuan</h4>
              <ul className="space-y-4">
                <li><a href="#" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Syarat & Ketentuan</a></li>
                <li><a href="#" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Kebijakan Privasi</a></li>
                <li><a href="#" className="text-slate-500 hover:text-[#64965a] font-medium transition-colors">Pusat Bantuan</a></li>
              </ul>
            </div>
          </div>

          <div className="pt-8 border-t border-[#e4ece3] flex flex-col md:flex-row justify-between items-center gap-4">
            <p className="text-slate-400 text-sm font-medium">&copy; {new Date().getFullYear()} Buroqet Logistics. All rights reserved.</p>
            <div className="flex gap-4">
              <Link to="/admin/login" className="text-slate-300 hover:text-slate-500 text-xs font-medium">Admin Portal</Link>
              <Link to="/courier/login" className="text-slate-300 hover:text-slate-500 text-xs font-medium">Fleet Portal</Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
