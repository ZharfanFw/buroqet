import { useAuthStore } from '../../store/auth.store';
import { useNavigate } from 'react-router-dom';

export default function ProfilePage() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="max-w-4xl mx-auto pb-12 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-3xl font-extrabold text-slate-800 mb-2">👤 Profil Saya</h1>
        <p className="text-slate-600">Kelola informasi akun dan preferensi Anda di sini.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* PROFILE CARD */}
        <div className="md:col-span-1">
          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 text-center">
            <div className="w-24 h-24 mx-auto bg-gradient-to-br from-[#64965a] to-[#79ae6f] rounded-full flex items-center justify-center text-white text-3xl font-bold shadow-lg mb-4">
              {user?.name?.charAt(0).toUpperCase() || 'U'}
            </div>
            <h2 className="text-xl font-extrabold text-slate-800 mb-1">{user?.name || 'User Buroqet'}</h2>
            <p className="text-slate-500 text-sm font-medium mb-4">{user?.email || 'user@buroqet.com'}</p>
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-[#64965a]/10 text-[#64965a] rounded-full text-xs font-bold uppercase tracking-wider mb-6">
              <span className="w-2 h-2 rounded-full bg-[#64965a]"></span>
              {user?.role || 'Pelanggan'}
            </div>
            <div className="border-t border-slate-100 pt-6">
              <button 
                onClick={handleLogout}
                className="w-full px-4 py-2 text-red-600 font-bold hover:bg-red-50 rounded-xl transition-colors"
              >
                Logout / Keluar
              </button>
            </div>
          </div>
        </div>

        {/* DETAILS SECTION */}
        <div className="md:col-span-2 space-y-6">
          {/* PERSONAL INFO */}
          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 md:p-8">
            <h3 className="text-lg font-bold text-slate-800 mb-6 flex items-center gap-2">
              <span className="text-[#64965a]">📋</span> Informasi Personal
            </h3>
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">Nama Lengkap</label>
                  <div className="font-semibold text-slate-800 bg-slate-50 px-4 py-3 rounded-xl border border-slate-100">
                    {user?.name || '-'}
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">Alamat Email</label>
                  <div className="font-semibold text-slate-800 bg-slate-50 px-4 py-3 rounded-xl border border-slate-100">
                    {user?.email || '-'}
                  </div>
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">Nomor Handphone</label>
                  <div className="font-semibold text-slate-800 bg-slate-50 px-4 py-3 rounded-xl border border-slate-100 flex justify-between items-center">
                    <span>+62 8XX-XXXX-XXXX</span>
                    <button className="text-xs text-[#64965a] font-bold hover:underline">Edit</button>
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">Kata Sandi</label>
                  <div className="font-semibold text-slate-800 bg-slate-50 px-4 py-3 rounded-xl border border-slate-100 flex justify-between items-center">
                    <span>••••••••</span>
                    <button className="text-xs text-[#64965a] font-bold hover:underline">Ubah</button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* ALAMAT TERSIMPAN */}
          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 md:p-8">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-lg font-bold text-slate-800 flex items-center gap-2">
                <span className="text-[#64965a]">🏠</span> Alamat Tersimpan
              </h3>
              <button className="text-sm font-bold text-[#64965a] hover:underline">+ Tambah Alamat</button>
            </div>
            
            <div className="space-y-4">
              <div className="border border-[#64965a] bg-[#64965a]/5 p-4 rounded-2xl relative">
                <div className="absolute top-4 right-4 flex gap-2">
                  <span className="px-2 py-1 bg-[#64965a] text-white text-[10px] font-bold rounded uppercase tracking-wider">Utama</span>
                </div>
                <h4 className="font-bold text-slate-800 mb-1">Rumah (Jakarta)</h4>
                <p className="text-sm text-slate-600 mb-2">Jl. Buroqet Logistic No. 99, Kebayoran Baru, Jakarta Selatan, 12110</p>
                <p className="text-xs text-slate-500 font-medium">Budi Santoso - 081234567890</p>
              </div>
              
              <div className="border border-slate-200 bg-slate-50 p-4 rounded-2xl hover:border-slate-300 transition-colors cursor-pointer">
                <h4 className="font-bold text-slate-800 mb-1">Kantor (Bandung)</h4>
                <p className="text-sm text-slate-600 mb-2">Gedung Buroqet Tower Lt. 5, Jl. Asia Afrika, Bandung, 40111</p>
                <p className="text-xs text-slate-500 font-medium">Budi Santoso - 081234567890</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
