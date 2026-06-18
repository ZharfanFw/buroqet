import { useState, useEffect, useCallback } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';
import type { ApiError } from '../../types';

interface CommissionLog {
  id: string;
  courier_id: string;
  awb: string;
  amount: number;
  status: string;
  delivered_at: string;
  created_at: string;
}

interface Earnings {
  courier_id: string;
  total_deliveries: number;
  total_amount: number;
  pending_amount: number;
  paid_amount: number;
}

export default function SettlementPage() {
  const [courierIdSearch, setCourierIdSearch] = useState('');
  const [earnings, setEarnings] = useState<Earnings | null>(null);
  const [loadingEarnings, setLoadingEarnings] = useState(false);

  const [formCourierId, setFormCourierId] = useState('');
  const [formAwb, setFormAwb] = useState('');
  const [formServiceType, setFormServiceType] = useState('REGULER');
  const [loadingManual, setLoadingManual] = useState(false);
  const [manualMsg, setManualMsg] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  const [commissions, setCommissions] = useState<CommissionLog[]>([]);
  const [loadingCommissions, setLoadingCommissions] = useState(true);

  const fetchCommissions = useCallback(async () => {
    setLoadingCommissions(true);
    try {
      const { data } = await apiClient.get(API_ENDPOINTS.settlement.list);
      setCommissions(Array.isArray(data) ? data : []);
    } catch (err: any) {
      console.error('Gagal fetch commissions:', err.message);
      setCommissions([]);
    } finally {
      setLoadingCommissions(false);
    }
  }, []);

  useEffect(() => {
    fetchCommissions();
  }, [fetchCommissions]);

  const fetchEarnings = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!courierIdSearch) return;
    setLoadingEarnings(true);
    try {
      const { data } = await apiClient.get(API_ENDPOINTS.settlement.earnings(courierIdSearch));
      setEarnings(data);
    } catch (err: any) {
      setEarnings(null);
      alert(`Gagal mengambil data: ${err.response?.data?.error || err.message}`);
    } finally {
      setLoadingEarnings(false);
    }
  };

  const handleManualCommission = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formCourierId || !formAwb) return;
    setLoadingManual(true);
    setManualMsg(null);
    try {
      await apiClient.post(API_ENDPOINTS.settlement.base, {
        courier_id: formCourierId,
        awb: formAwb,
        service_type: formServiceType
      });
      setManualMsg({ text: 'Komisi berhasil dicatat secara manual.', type: 'success' });
      setFormAwb('');
      fetchCommissions(); // refresh table
    } catch (err: any) {
      const apiErr = err.response?.data as ApiError;
      setManualMsg({ text: `Gagal mencatat: ${apiErr?.message || apiErr?.error || err.message}`, type: 'error' });
    } finally {
      setLoadingManual(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>💳 Keuangan &amp; Settlement</h1>
        <p>Rekonsiliasi komisi kurir dan pencatatan pendapatan</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '24px' }}>

        {/* Panel Cek Penghasilan */}
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <h2 style={{ fontSize: '18px', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span>🔍</span> Cek Penghasilan Kurir
          </h2>
          <form onSubmit={fetchEarnings} style={{ display: 'flex', gap: '8px' }}>
            <input
              type="text"
              className="input"
              placeholder="Masukkan Courier ID (cth: KURIR-001)"
              value={courierIdSearch}
              onChange={e => setCourierIdSearch(e.target.value)}
              required
            />
            <button type="submit" className="btn btn-primary" disabled={loadingEarnings}>
              {loadingEarnings ? 'Mencari...' : 'Cari'}
            </button>
          </form>

          {earnings && (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginTop: '8px' }}>
              <div style={{ padding: '16px', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)', border: '1px solid var(--border)' }}>
                <div style={{ color: 'var(--text-muted)', fontSize: '12px', marginBottom: '4px' }}>Kurir ID</div>
                <div style={{ fontSize: '18px', fontWeight: 600 }}>{earnings.courier_id}</div>
                <div style={{ color: 'var(--text-secondary)', fontSize: '12px', marginTop: '4px' }}>Total Pengiriman: {earnings.total_deliveries} paket</div>
              </div>

              <div style={{ display: 'flex', gap: '12px' }}>
                <div style={{ flex: 1, padding: '16px', background: 'rgba(245, 158, 11, 0.1)', borderRadius: 'var(--radius-sm)', border: '1px solid rgba(245, 158, 11, 0.2)' }}>
                  <div style={{ color: 'var(--warning)', fontSize: '12px', fontWeight: 600 }}>Belum Cair (Pending)</div>
                  <div style={{ fontSize: '20px', fontWeight: 700, color: 'var(--warning)', marginTop: '4px' }}>
                    Rp {earnings.pending_amount.toLocaleString('id-ID')}
                  </div>
                </div>
                <div style={{ flex: 1, padding: '16px', background: 'rgba(16, 185, 129, 0.1)', borderRadius: 'var(--radius-sm)', border: '1px solid rgba(16, 185, 129, 0.2)' }}>
                  <div style={{ color: 'var(--success)', fontSize: '12px', fontWeight: 600 }}>Sudah Dibayar (Paid)</div>
                  <div style={{ fontSize: '20px', fontWeight: 700, color: 'var(--success)', marginTop: '4px' }}>
                    Rp {earnings.paid_amount.toLocaleString('id-ID')}
                  </div>
                </div>
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingTop: '16px', borderTop: '1px dashed var(--border)' }}>
                <span style={{ fontWeight: 600 }}>Total Penghasilan:</span>
                <span style={{ fontSize: '18px', fontWeight: 700, color: 'var(--primary-light)' }}>
                  Rp {earnings.total_amount.toLocaleString('id-ID')}
                </span>
              </div>
            </div>
          )}
        </div>

        {/* Panel Input Manual */}
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <h2 style={{ fontSize: '18px', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span>✍️</span> Catat Komisi Manual
          </h2>
          <p style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>
            Gunakan form ini jika event Kafka package.delivered gagal atau terlewat.
          </p>

          <form onSubmit={handleManualCommission} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>Courier ID</label>
              <input
                type="text"
                className="input"
                placeholder="Contoh: KURIR-001"
                value={formCourierId}
                onChange={e => setFormCourierId(e.target.value)}
                required
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>AWB / Resi</label>
              <input
                type="text"
                className="input"
                placeholder="Contoh: BQ-2024-JKT-001"
                value={formAwb}
                onChange={e => setFormAwb(e.target.value)}
                required
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontSize: '12px', color: 'var(--text-secondary)' }}>Tipe Layanan</label>
              <select
                className="input"
                value={formServiceType}
                onChange={e => setFormServiceType(e.target.value)}
              >
                <option value="REGULER">Reguler (Rp 3.500)</option>
                <option value="EXPRESS">Express (Rp 5.000)</option>
                <option value="SAME_DAY">Same Day (Rp 7.500)</option>
              </select>
            </div>
            <button type="submit" className="btn btn-primary" disabled={loadingManual} style={{ marginTop: '8px' }}>
              {loadingManual ? 'Memproses...' : 'Catat Komisi'}
            </button>
          </form>

          {manualMsg && (
            <div style={{
              marginTop: '8px',
              padding: '10px 12px',
              borderRadius: 'var(--radius-sm)',
              fontSize: '13px',
              background: manualMsg.type === 'error' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(16, 185, 129, 0.1)',
              color: manualMsg.type === 'error' ? 'var(--danger)' : 'var(--success)',
              border: `1px solid ${manualMsg.type === 'error' ? 'rgba(239, 68, 68, 0.2)' : 'rgba(16, 185, 129, 0.2)'}`
            }}>
              {manualMsg.text}
            </div>
          )}
        </div>
      </div>

      {/* Tabel Commission Logs */}
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <h2 style={{ fontSize: '16px' }}>📊 Log Komisi Keseluruhan ({commissions.length})</h2>
          <button
            onClick={fetchCommissions}
            className="btn"
            style={{ fontSize: '12px', padding: '6px 12px', background: 'var(--bg-surface)', border: '1px solid var(--border)', borderRadius: 'var(--radius-sm)', cursor: 'pointer', color: 'var(--text-secondary)' }}
          >
            🔄 Refresh
          </button>
        </div>

        {loadingCommissions ? (
          <div className="empty-state" style={{ padding: '20px' }}>
            <p>Memuat data...</p>
          </div>
        ) : commissions.length === 0 ? (
          <div className="empty-state" style={{ padding: '20px' }}>
            <p>Belum ada komisi tercatat. Coba tambah komisi manual di atas!</p>
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '13px' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--border)' }}>
                  {['AWB', 'Courier ID', 'Jumlah', 'Status', 'Tanggal Deliver'].map(h => (
                    <th key={h} style={{ textAlign: 'left', padding: '10px 12px', color: 'var(--text-muted)', fontWeight: 600, fontSize: '11px', textTransform: 'uppercase', letterSpacing: '0.05em' }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {commissions.map((log, idx) => (
                  <tr key={log.id} style={{ borderBottom: '1px solid var(--border)', background: idx % 2 === 0 ? 'transparent' : 'rgba(255,255,255,0.02)' }}>
                    <td style={{ padding: '10px 12px', fontWeight: 600, fontFamily: 'monospace', color: 'var(--primary-light)' }}>{log.awb}</td>
                    <td style={{ padding: '10px 12px', color: 'var(--text-secondary)' }}>{log.courier_id}</td>
                    <td style={{ padding: '10px 12px', fontWeight: 700, color: 'var(--success)' }}>
                      Rp {log.amount.toLocaleString('id-ID')}
                    </td>
                    <td style={{ padding: '10px 12px' }}>
                      <span style={{
                        padding: '3px 10px',
                        borderRadius: '999px',
                        fontSize: '11px',
                        fontWeight: 700,
                        background: log.status === 'PAID' ? 'rgba(16,185,129,0.15)' : 'rgba(245,158,11,0.15)',
                        color: log.status === 'PAID' ? '#10b981' : '#f59e0b'
                      }}>
                        {log.status}
                      </span>
                    </td>
                    <td style={{ padding: '10px 12px', color: 'var(--text-muted)', fontSize: '11px' }}>
                      {log.delivered_at ? new Date(log.delivered_at).toLocaleString('id-ID') : '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
