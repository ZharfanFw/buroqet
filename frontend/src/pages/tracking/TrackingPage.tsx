import { useState } from 'react';
import TrackingResult from './TrackingResult';
import type { TrackingHistory } from '../../types';
import { API_ENDPOINTS } from '../../utils/api-config';
import './TrackingPage.css';

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
      // If we didn't return early from mock
      if (trimmed !== 'BQ-2024-BENTO-123' && trimmed !== 'BQ-2024-JKT-001') {
        setLoading(false);
      }
    }
  };

  return (
    <div className="tracking-page">
      {/* Header */}
      <div className="page-header">
        <h1>📦 Lacak Paket</h1>
        <p>Masukkan nomor resi untuk melihat status pengiriman secara real-time</p>
      </div>

      {/* Search Box */}
      <div className="tracking-search-card card">
        <form className="tracking-search-form" onSubmit={handleSearch}>
          <div className="search-input-wrap">
            <span className="search-icon">🔍</span>
            <input
              id="tracking-awb-input"
              type="text"
              className="search-input"
              placeholder="Masukkan nomor resi, contoh: BQ-2024-JKT-001"
              value={inputVal}
              onChange={(e) => setInputVal(e.target.value)}
              autoComplete="off"
              spellCheck={false}
            />
            {inputVal && (
              <button
                type="button"
                className="search-clear"
                onClick={() => { setInputVal(''); setData(null); setError(''); }}
              >✕</button>
            )}
          </div>
          <button
            id="tracking-search-btn"
            type="submit"
            className="btn btn-primary search-btn"
            disabled={loading || !inputVal.trim()}
          >
            {loading ? <span className="spinner" /> : 'Lacak'}
          </button>
        </form>

        {/* Quick AWB examples */}
        <div className="quick-examples">
          <span className="quick-label">Contoh resi:</span>
          {['BQ-2024-JKT-001', 'BQ-2024-BENTO-123'].map(ex => (
            <button
              key={ex}
              className="quick-chip"
              onClick={() => setInputVal(ex)}
            >{ex}</button>
          ))}
        </div>
      </div>

      {/* Error State */}
      {error && (
        <div className="tracking-error card">
          <span className="error-icon">⚠️</span>
          <div>
            <strong>Paket tidak ditemukan</strong>
            <p>{error}</p>
          </div>
        </div>
      )}

      {/* Loading Skeleton */}
      {loading && (
        <div className="tracking-skeleton card">
          <div className="skeleton-header">
            <div className="skeleton-block w-40" />
            <div className="skeleton-block w-24" />
          </div>
          <div className="skeleton-timeline">
            {[1, 2, 3].map(i => (
              <div key={i} className="skeleton-event">
                <div className="skeleton-dot" />
                <div className="skeleton-lines">
                  <div className="skeleton-block w-60" />
                  <div className="skeleton-block w-40" />
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Result */}
      {data && !loading && <TrackingResult awb={awb} data={data} />}

      {/* Empty State */}
      {!data && !loading && !error && (
        <div className="tracking-empty">
          <div className="tracking-illustration">
            <svg viewBox="0 0 200 180" fill="none" xmlns="http://www.w3.org/2000/svg">
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
          <h3>Lacak pengiriman Anda</h3>
          <p>Masukkan nomor resi di atas untuk melihat perjalanan paket secara real-time</p>
        </div>
      )}
    </div>
  );
}
