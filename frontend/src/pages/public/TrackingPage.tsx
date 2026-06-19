import { useState } from 'react';
import TrackingResult from './TrackingResult';
import type { TrackingHistory } from '../../types';
import { API_ENDPOINTS } from '../../utils/api-config';

export default function TrackingPage() {
  const [awb, setAwb] = useState('');
  const [inputVal, setInputVal] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [data, setData] = useState<TrackingHistory | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = inputVal.trim();
    if (!trimmed) return;

    setLoading(true);
    setError('');
    setData(null);
    setAwb(trimmed);

    try {
      // Mock data for UI testing without backend
      if (trimmed === 'BQ-2024-BENTO-123' || trimmed === 'BQ-2024-JKT-001') {
        setTimeout(() => {
          setData({
            awb: trimmed,
            total: 4,
            events: [
              {
                id: '1', awb: trimmed, status: 'INBOUND',
                location: 'Gudang Sortir Jakarta', hub_id: 'HUB-JKT-01',
                description: 'Paket telah diterima di gudang asal',
                latitude: -6.2088, longitude: 106.8456,
                timestamp: new Date(Date.now() - 86400000 * 2).toISOString(),
                created_at: new Date(Date.now() - 86400000 * 2).toISOString(),
                source: 'WMS'
              },
              {
                id: '2', awb: trimmed, status: 'ON_TRANSIT',
                location: 'Perjalanan ke Bandung', hub_id: '',
                description: 'Paket diberangkatkan menuju kota tujuan',
                latitude: -6.5000, longitude: 107.2000,
                timestamp: new Date(Date.now() - 86400000 * 1.5).toISOString(),
                created_at: new Date(Date.now() - 86400000 * 1.5).toISOString(),
                source: 'Dispatch'
              },
              {
                id: '3', awb: trimmed, status: 'AT_HUB',
                location: 'Gudang Sortir Bandung', hub_id: 'HUB-BDG-01',
                description: 'Paket telah tiba di gudang tujuan',
                latitude: -6.9147, longitude: 107.6098,
                timestamp: new Date(Date.now() - 86400000 * 0.5).toISOString(),
                created_at: new Date(Date.now() - 86400000 * 0.5).toISOString(),
                source: 'WMS'
              },
              {
                id: '4', awb: trimmed, status: 'OUT_FOR_DELIVERY',
                location: 'Bandung', hub_id: 'HUB-BDG-01',
                description: 'Paket sedang dibawa kurir (Budi) menuju alamat penerima',
                latitude: -6.9200, longitude: 107.6100,
                timestamp: new Date().toISOString(),
                created_at: new Date().toISOString(),
                source: 'Dispatch'
              }
            ]
          });
          setLoading(false);
        }, 1200);
        return;
      }

      const res = await fetch(`${API_ENDPOINTS.tracking.byAwb(encodeURIComponent(trimmed))}/history`);

      if (res.status === 404) {
        setError(`Nomor resi "${trimmed}" tidak ditemukan.`);
        return;
      }
      if (!res.ok) {
        setError(`Server error: ${res.status}`);
        return;
      }

      const json: TrackingHistory = await res.json();
      setData(json);
    } catch {
      setError('Gagal terhubung ke server. Pastikan tracking-service sedang berjalan.');
    } finally {
      if (trimmed !== 'BQ-2024-BENTO-123' && trimmed !== 'BQ-2024-JKT-001') {
        setLoading(false);
      }
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col items-center py-12 px-4 sm:px-6 lg:px-8 font-sans">
      
      {/* Header Container - Centered */}
      <div className="w-full max-w-2xl text-center mb-10">
        <h1 className="text-4xl md:text-5xl font-extrabold text-slate-900 mb-4 tracking-tight">
          📦 Lacak Paket
        </h1>
        <p className="text-lg text-slate-600">
          Masukkan nomor resi untuk melihat status pengiriman secara real-time
        </p>
      </div>

      {/* Search Box Card - Centered with distinct border/background */}
      <div className="w-full max-w-2xl bg-white rounded-2xl shadow-sm border-2 border-slate-200 p-6 md:p-8 mb-8">
        <form className="flex flex-col sm:flex-row gap-4" onSubmit={handleSearch}>
          <div className="relative flex-grow">
            <span className="absolute left-4 top-1/2 transform -translate-y-1/2 text-xl text-slate-400">
              🔍
            </span>
            <input
              id="tracking-awb-input"
              type="text"
              className="w-full pl-12 pr-12 py-3.5 bg-slate-50 border-2 border-slate-200 text-slate-800 rounded-xl focus:bg-white focus:outline-none focus:border-[#79ae6f] focus:ring-4 focus:ring-[#79ae6f]/20 transition-all font-medium placeholder-slate-400"
              placeholder="Contoh: BQ-2024-JKT-001"
              value={inputVal}
              onChange={(e) => setInputVal(e.target.value)}
              autoComplete="off"
              spellCheck={false}
            />
            {inputVal && (
              <button
                type="button"
                className="absolute right-4 top-1/2 transform -translate-y-1/2 text-slate-400 hover:text-slate-600 hover:bg-slate-200 rounded-full w-6 h-6 flex items-center justify-center transition-colors"
                onClick={() => { setInputVal(''); setData(null); setError(''); }}
              >
                ✕
              </button>
            )}
          </div>
          <button
            id="tracking-search-btn"
            type="submit"
            className="flex-shrink-0 bg-[#64965a] hover:bg-[#53804a] text-white px-8 py-3.5 rounded-xl font-semibold transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center min-w-[140px] shadow-sm hover:shadow-md"
            disabled={loading || !inputVal.trim()}
          >
            {loading ? (
              <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            ) : (
              'Lacak'
            )}
          </button>
        </form>

        {/* Quick AWB examples */}
        <div className="mt-6 flex flex-wrap items-center gap-3">
          <span className="text-sm font-medium text-slate-500">Contoh resi:</span>
          {['BQ-2024-JKT-001', 'BQ-2024-BENTO-123'].map(ex => (
            <button
              key={ex}
              className="text-sm px-4 py-1.5 bg-slate-100 hover:bg-slate-200 border border-slate-200 text-slate-700 rounded-full cursor-pointer transition-colors font-medium"
              onClick={() => setInputVal(ex)}
            >
              {ex}
            </button>
          ))}
        </div>
      </div>

      {/* Error State */}
      {error && (
        <div className="w-full max-w-2xl bg-red-50 border-l-4 border-red-500 text-red-700 p-5 rounded-r-xl shadow-sm flex items-start gap-4 mb-8">
          <span className="text-2xl mt-0.5">⚠️</span>
          <div>
            <strong className="block font-bold mb-1">Paket tidak ditemukan</strong>
            <p className="text-red-600">{error}</p>
          </div>
        </div>
      )}

      {/* Loading Skeleton */}
      {loading && (
        <div className="w-full max-w-2xl bg-white border border-slate-200 rounded-2xl p-6 shadow-sm mb-8 animate-pulse">
          <div className="mb-8">
            <div className="h-6 bg-slate-200 rounded w-40 mb-3"></div>
            <div className="h-4 bg-slate-100 rounded w-24"></div>
          </div>
          <div className="space-y-6">
            {[1, 2, 3].map(i => (
              <div key={i} className="flex gap-4">
                <div className="w-4 h-4 rounded-full bg-slate-200 mt-1"></div>
                <div className="flex-1 space-y-3">
                  <div className="h-4 bg-slate-200 rounded w-3/4"></div>
                  <div className="h-3 bg-slate-100 rounded w-1/2"></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Result Container */}
      <div className="w-full max-w-2xl">
        {data && !loading && <TrackingResult awb={awb} data={data} />}
      </div>

      {/* Empty State */}
      {!data && !loading && !error && (
        <div className="w-full max-w-2xl flex flex-col items-center justify-center py-10 px-4">
          <div className="w-56 h-auto mb-8">
            <svg viewBox="0 0 200 180" fill="none" xmlns="http://www.w3.org/2000/svg" className="w-full h-full">
              {/* Road */}
              <rect x="20" y="130" width="160" height="20" rx="4" fill="#f1f5f9"/>
              <rect x="90" y="135" width="20" height="5" rx="2" fill="#e2e8f0"/>
              <rect x="55" y="135" width="15" height="5" rx="2" fill="#e2e8f0"/>
              <rect x="130" y="135" width="15" height="5" rx="2" fill="#e2e8f0"/>
              {/* Truck body */}
              <rect x="30" y="95" width="80" height="40" rx="6" fill="#64965a"/>
              <rect x="110" y="105" width="40" height="30" rx="4" fill="#79ae6f"/>
              {/* Truck cab window */}
              <rect x="114" y="108" width="18" height="14" rx="3" fill="#ffffff" opacity="0.9"/>
              {/* Truck wheels */}
              <circle cx="55" cy="138" r="10" fill="#334155" stroke="#9bc692" strokeWidth="3"/>
              <circle cx="55" cy="138" r="4" fill="#cbd5e1"/>
              <circle cx="130" cy="138" r="10" fill="#334155" stroke="#9bc692" strokeWidth="3"/>
              <circle cx="130" cy="138" r="4" fill="#cbd5e1"/>
              {/* Package on truck */}
              <rect x="45" y="75" width="55" height="25" rx="4" fill="#e2e8f0"/>
              <line x1="72" y1="75" x2="72" y2="100" stroke="#cbd5e1" strokeWidth="2"/>
              <line x1="45" y1="87" x2="100" y2="87" stroke="#cbd5e1" strokeWidth="2"/>
              {/* Signal waves */}
              <path d="M155 60 Q165 50 175 60" stroke="#9bc692" strokeWidth="2.5" strokeLinecap="round" fill="none"/>
              <path d="M150 50 Q165 35 180 50" stroke="#9bc692" strokeWidth="2" strokeLinecap="round" fill="none" opacity="0.6"/>
              <path d="M145 40 Q165 20 185 40" stroke="#9bc692" strokeWidth="1.5" strokeLinecap="round" fill="none" opacity="0.3"/>
              {/* Location pin */}
              <circle cx="165" cy="70" r="8" fill="#ef4444"/>
              <circle cx="165" cy="70" r="3" fill="white"/>
              <path d="M165 78 L162 85 L165 82 L168 85 Z" fill="#ef4444"/>
              {/* Stars / sparkles */}
              <circle cx="25" cy="60" r="2" fill="#79ae6f" opacity="0.5"/>
              <circle cx="15" cy="80" r="1.5" fill="#9bc692" opacity="0.4"/>
              <circle cx="185" cy="100" r="2" fill="#64965a" opacity="0.4"/>
            </svg>
          </div>
          <h3 className="text-xl font-bold text-slate-800 mb-2">Lacak pengiriman Anda</h3>
          <p className="text-slate-500 text-center max-w-md">
            Masukkan nomor resi di atas untuk melihat perjalanan paket secara real-time
          </p>
        </div>
      )}

    </div>
  );
}