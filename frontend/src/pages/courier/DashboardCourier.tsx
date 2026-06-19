import { useNavigate } from 'react-router-dom';

export default function CourierDashboard() {
  const navigate = useNavigate();

  return (
    <div className="animate-fade-in space-y-8">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-extrabold text-slate-800">Agent Dashboard</h2>
          <p className="text-slate-500 font-medium mt-1">Tugas dan rute pengiriman Anda hari ini</p>
        </div>
        <div className="flex gap-3">
          <button className="px-4 py-2 bg-white border border-[#e4ece3] text-slate-600 font-bold rounded-xl hover:bg-slate-50 shadow-sm transition-all duration-300">
            Scanner Mode
          </button>
          <button className="px-4 py-2 bg-[#64965a] hover:bg-[#53804a] text-white font-bold rounded-xl shadow-md hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 flex items-center gap-2">
            <span>🚀</span> Start Route
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-3xl p-6 shadow-sm border border-[#e4ece3] hover:shadow-lg hover:-translate-y-1 transition-all duration-300 group">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 rounded-2xl bg-indigo-50 text-indigo-500 flex items-center justify-center text-2xl group-hover:bg-indigo-500 group-hover:text-white transition-colors duration-300">
              📋
            </div>
          </div>
          <p className="text-slate-500 font-bold text-sm uppercase tracking-wider">Pending Pickups</p>
          <h3 className="text-3xl font-extrabold text-slate-800 mt-1">5</h3>
        </div>

        <div className="bg-white rounded-3xl p-6 shadow-sm border border-[#e4ece3] hover:shadow-lg hover:-translate-y-1 transition-all duration-300 group">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 rounded-2xl bg-orange-50 text-orange-500 flex items-center justify-center text-2xl group-hover:bg-orange-500 group-hover:text-white transition-colors duration-300">
              🛵
            </div>
          </div>
          <p className="text-slate-500 font-bold text-sm uppercase tracking-wider">Active Deliveries</p>
          <h3 className="text-3xl font-extrabold text-slate-800 mt-1">12</h3>
        </div>

        <div className="bg-white rounded-3xl p-6 shadow-sm border border-[#e4ece3] hover:shadow-lg hover:-translate-y-1 transition-all duration-300 group">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 rounded-2xl bg-[#64965a]/10 text-[#64965a] flex items-center justify-center text-2xl group-hover:bg-[#64965a] group-hover:text-white transition-colors duration-300">
              ✅
            </div>
          </div>
          <p className="text-slate-500 font-bold text-sm uppercase tracking-wider">Completed Today</p>
          <h3 className="text-3xl font-extrabold text-slate-800 mt-1">8</h3>
        </div>
      </div>

      {/* Main Content Area: Routes */}
      <div className="grid grid-cols-1 gap-8">
        <div className="bg-white rounded-3xl p-8 shadow-sm border border-[#e4ece3] transition-all hover:shadow-lg">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h3 className="text-xl font-extrabold text-slate-800">Current Assignment Route</h3>
              <p className="text-sm font-medium text-slate-500 mt-1">Optimized sequence for your active deliveries</p>
            </div>
            <button className="px-4 py-2 bg-slate-50 border border-[#e4ece3] text-slate-600 font-bold rounded-xl hover:bg-slate-100 transition-colors text-sm">
              Refresh Route
            </button>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-[#e4ece3]">
                  <th className="py-4 font-bold text-slate-400 uppercase text-xs tracking-wider">Seq</th>
                  <th className="py-4 font-bold text-slate-400 uppercase text-xs tracking-wider">AWB Number</th>
                  <th className="py-4 font-bold text-slate-400 uppercase text-xs tracking-wider">Recipient / Address</th>
                  <th className="py-4 font-bold text-slate-400 uppercase text-xs tracking-wider">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#e4ece3]">
                <tr className="hover:bg-slate-50 transition-colors">
                  <td className="py-4">
                    <div className="w-8 h-8 rounded-full bg-[#64965a] text-white flex items-center justify-center font-extrabold shadow-md">
                      1
                    </div>
                  </td>
                  <td className="py-4 font-bold text-slate-700">BQ-2024-JKT-102</td>
                  <td className="py-4">
                    <div className="font-bold text-slate-800">Ahmad Fauzi</div>
                    <div className="text-sm font-medium text-slate-500 mt-0.5">Jl. Mawar No. 15, Jaksel</div>
                  </td>
                  <td className="py-4">
                    <button 
                      className="px-4 py-2 bg-[#64965a] hover:bg-[#53804a] text-white font-bold rounded-xl shadow-sm hover:shadow-md transition-all text-sm"
                      onClick={() => navigate('/courier/epod')}
                    >
                      Submit ePOD
                    </button>
                  </td>
                </tr>
                <tr className="hover:bg-slate-50 transition-colors">
                  <td className="py-4">
                    <div className="w-8 h-8 rounded-full bg-slate-100 text-slate-500 flex items-center justify-center font-extrabold border border-[#e4ece3]">
                      2
                    </div>
                  </td>
                  <td className="py-4 font-bold text-slate-700">BQ-2024-JKT-088</td>
                  <td className="py-4">
                    <div className="font-bold text-slate-800">Siti Aminah</div>
                    <div className="text-sm font-medium text-slate-500 mt-0.5">Jl. Melati No. 4, Jaksel</div>
                  </td>
                  <td className="py-4">
                    <button className="px-4 py-2 bg-white border border-[#e4ece3] text-slate-600 font-bold rounded-xl hover:bg-slate-50 transition-all text-sm">
                      Directions
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
