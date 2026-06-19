export default function DashboardCustomer() {
  return (
    <div className="animate-fade-in">
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-slate-800">Ringkasan Akun</h2>
        <p className="text-slate-500">Kelola pesanan dan lacak pengiriman Anda.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-4 hover:shadow-md transition-shadow">
          <div className="w-12 h-12 rounded-full bg-indigo-50 text-indigo-600 flex items-center justify-center text-xl">
            📦
          </div>
          <div>
            <p className="text-sm text-slate-500 font-medium">Total Pesanan</p>
            <p className="text-2xl font-bold text-slate-800">12</p>
          </div>
        </div>
        
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-4 hover:shadow-md transition-shadow">
          <div className="w-12 h-12 rounded-full bg-amber-50 text-amber-600 flex items-center justify-center text-xl">
            🚚
          </div>
          <div>
            <p className="text-sm text-slate-500 font-medium">Dalam Perjalanan</p>
            <p className="text-2xl font-bold text-slate-800">2</p>
          </div>
        </div>

        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-4 hover:shadow-md transition-shadow">
          <div className="w-12 h-12 rounded-full bg-emerald-50 text-emerald-600 flex items-center justify-center text-xl">
            ✅
          </div>
          <div>
            <p className="text-sm text-slate-500 font-medium">Selesai</p>
            <p className="text-2xl font-bold text-slate-800">10</p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
        <div className="p-6 border-b border-slate-100 flex justify-between items-center">
          <h3 className="text-lg font-bold text-slate-800">Pesanan Aktif</h3>
          <button className="text-indigo-600 text-sm font-semibold hover:text-indigo-700">Lihat Semua</button>
        </div>
        <div className="p-6 text-center text-slate-500 py-12">
          Belum ada pesanan yang sedang aktif saat ini.
        </div>
      </div>
    </div>
  );
}
