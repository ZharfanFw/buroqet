import { useState, useEffect, useCallback } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';
import type { ApiError } from '../../types';

interface Package {
  id: string;
  awb: string;
  hub_id: string;
  manifest_id: string | null;
  status: string;
  scanned_at: string;
  created_at: string;
}

const STATUS_COLORS: Record<string, { bg: string; color: string }> = {
  INBOUND:     { bg: 'rgba(59,130,246,0.15)',  color: '#60a5fa' },
  ON_TRANSIT:  { bg: 'rgba(245,158,11,0.15)', color: '#f59e0b' },
  OUTBOUND:    { bg: 'rgba(139,92,246,0.15)', color: '#a78bfa' },
  DELIVERED:   { bg: 'rgba(16,185,129,0.15)', color: '#10b981' },
};

export default function WarehousePage() {
  const [inboundAwb, setInboundAwb] = useState('');
  const [hubId, setHubId] = useState('HUB-JKT-01');
  const [manifestId, setManifestId] = useState('');

  const [loadingInbound, setLoadingInbound] = useState(false);
  const [loadingDispatch, setLoadingDispatch] = useState(false);
  const [loadingPackages, setLoadingPackages] = useState(true);

  const [packages, setPackages] = useState<Package[]>([]);
  const [logs, setLogs] = useState<{ id: string; time: string; msg: string; type: 'success' | 'error' }[]>([]);

  const addLog = (msg: string, type: 'success' | 'error') => {
    setLogs(prev => [{ id: Math.random().toString(), time: new Date().toLocaleTimeString(), msg, type }, ...prev].slice(0, 10));
  };

  const fetchPackages = useCallback(async () => {
    setLoadingPackages(true);
    try {
      const { data } = await apiClient.get(API_ENDPOINTS.warehouse.packages);
      setPackages(Array.isArray(data) ? data : []);
    } catch (err: any) {
      console.error('Gagal fetch packages:', err.message);
      setPackages([]);
    } finally {
      setLoadingPackages(false);
    }
  }, []);

  useEffect(() => {
    fetchPackages();
  }, [fetchPackages]);

  const handleInbound = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inboundAwb) return;
    setLoadingInbound(true);
    try {
      await apiClient.post(API_ENDPOINTS.warehouse.inbound, { awb: inboundAwb, hub_id: hubId });
      addLog(`📦 Inbound sukses: AWB ${inboundAwb} masuk ke ${hubId}`, 'success');
      setInboundAwb('');
      fetchPackages(); // refresh table
    } catch (err: any) {
      const apiErr = err.response?.data as ApiError;
      addLog(`❌ Inbound gagal: ${apiErr?.message || apiErr?.error || err.message}`, 'error');
    } finally {
      setLoadingInbound(false);
    }
  };

  const handleDispatch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!manifestId) return;
    setLoadingDispatch(true);
    try {
      await apiClient.post(API_ENDPOINTS.warehouse.dispatch, { manifest_id: manifestId });
      addLog(`🚚 Dispatch sukses: Manifest ${manifestId} dikirim`, 'success');
      setManifestId('');
      fetchPackages(); // refresh table
    } catch (err: any) {
      const apiErr = err.response?.data as ApiError;
      addLog(`❌ Dispatch gagal: ${apiErr?.message || apiErr?.error || err.message}`, 'error');
    } finally {
      setLoadingDispatch(false);
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>🏭 Warehouse Management</h1>
        <p>Manajemen operasional inbound paket dan dispatch manifest.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '24px', marginBottom: '24px' }}>

        {/* Panel Inbound */}
        <div className="card">
          <h2 style={{ fontSize: '18px', marginBottom: '16px', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span>📥</span> Scan Inbound
          </h2>
          <form onSubmit={handleInbound} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>Hub ID</label>
              <input
                type="text"
                className="input"
                value={hubId}
                onChange={e => setHubId(e.target.value)}
                required
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>AWB Scanner</label>
              <input
                type="text"
                className="input"
                placeholder="Contoh: BQ-2024-JKT-001"
                value={inboundAwb}
                onChange={e => setInboundAwb(e.target.value)}
                required
              />
            </div>
            <button type="submit" className="btn btn-primary" disabled={loadingInbound} style={{ marginTop: '8px' }}>
              {loadingInbound ? 'Processing...' : 'Simpan Inbound'}
            </button>
          </form>
        </div>

        {/* Panel Dispatch */}
        <div className="card">
          <h2 style={{ fontSize: '18px', marginBottom: '16px', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span>📤</span> Dispatch Manifest
          </h2>
          <form onSubmit={handleDispatch} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>Manifest ID</label>
              <input
                type="text"
                className="input"
                placeholder="Contoh: MNF-JKT-SBY-001"
                value={manifestId}
                onChange={e => setManifestId(e.target.value)}
                required
              />
            </div>
            <button type="submit" className="btn btn-primary" disabled={loadingDispatch} style={{ marginTop: '8px', background: 'var(--warning)', color: '#000' }}>
              {loadingDispatch ? 'Processing...' : 'Dispatch Sekarang'}
            </button>
          </form>
        </div>
      </div>

      {/* Tabel Daftar Paket */}
      <div className="card" style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <h2 style={{ fontSize: '16px' }}>📋 Daftar Paket di Gudang ({packages.length})</h2>
          <button
            onClick={fetchPackages}
            className="btn"
            style={{ fontSize: '12px', padding: '6px 12px', background: 'var(--bg-surface)', border: '1px solid var(--border)', borderRadius: 'var(--radius-sm)', cursor: 'pointer', color: 'var(--text-secondary)' }}
          >
            🔄 Refresh
          </button>
        </div>

        {loadingPackages ? (
          <div className="empty-state" style={{ padding: '20px' }}>
            <p>Memuat data...</p>
          </div>
        ) : packages.length === 0 ? (
          <div className="empty-state" style={{ padding: '20px' }}>
            <p>Belum ada paket di gudang. Coba scan inbound paket pertama Anda!</p>
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '13px' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--border)' }}>
                  {['AWB', 'Hub', 'Status', 'Manifest ID', 'Waktu Scan'].map(h => (
                    <th key={h} style={{ textAlign: 'left', padding: '10px 12px', color: 'var(--text-muted)', fontWeight: 600, fontSize: '11px', textTransform: 'uppercase', letterSpacing: '0.05em' }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {packages.map((pkg, idx) => {
                  const sc = STATUS_COLORS[pkg.status] || { bg: 'rgba(100,100,100,0.1)', color: 'var(--text-secondary)' };
                  return (
                    <tr key={pkg.id} style={{ borderBottom: '1px solid var(--border)', background: idx % 2 === 0 ? 'transparent' : 'rgba(255,255,255,0.02)' }}>
                      <td style={{ padding: '10px 12px', fontWeight: 600, fontFamily: 'monospace', color: 'var(--primary-light)' }}>{pkg.awb}</td>
                      <td style={{ padding: '10px 12px', color: 'var(--text-secondary)' }}>{pkg.hub_id}</td>
                      <td style={{ padding: '10px 12px' }}>
                        <span style={{ padding: '3px 10px', borderRadius: '999px', fontSize: '11px', fontWeight: 700, background: sc.bg, color: sc.color }}>
                          {pkg.status}
                        </span>
                      </td>
                      <td style={{ padding: '10px 12px', color: 'var(--text-muted)', fontFamily: 'monospace', fontSize: '11px' }}>{pkg.manifest_id || '—'}</td>
                      <td style={{ padding: '10px 12px', color: 'var(--text-muted)', fontSize: '11px' }}>
                        {pkg.scanned_at ? new Date(pkg.scanned_at).toLocaleString('id-ID') : '—'}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Operation Logs */}
      <div className="card">
        <h2 style={{ fontSize: '16px', marginBottom: '16px' }}>📝 Riwayat Operasi Sesi Ini</h2>
        {logs.length === 0 ? (
          <div className="empty-state" style={{ padding: '20px' }}>
            <p>Belum ada operasi yang dilakukan di sesi ini.</p>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            {logs.map(log => (
              <div key={log.id} style={{
                padding: '12px 16px',
                borderRadius: 'var(--radius-sm)',
                background: log.type === 'error' ? 'rgba(239, 68, 68, 0.1)' : 'var(--bg-surface)',
                borderLeft: `4px solid ${log.type === 'error' ? 'var(--danger)' : 'var(--success)'}`,
                display: 'flex',
                gap: '12px',
                alignItems: 'center'
              }}>
                <span style={{ fontSize: '12px', color: 'var(--text-muted)', whiteSpace: 'nowrap' }}>{log.time}</span>
                <span style={{ color: log.type === 'error' ? 'var(--danger)' : 'var(--text-primary)' }}>{log.msg}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
