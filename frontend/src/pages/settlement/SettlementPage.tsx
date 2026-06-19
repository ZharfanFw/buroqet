import { useState, useEffect } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';

interface CommissionLog {
  id: string;
  courier_id: string;
  awb: string;
  amount: number;
  status: string;
  created_at: string;
}

export default function SettlementPage() {
  const [commissions, setCommissions] = useState<CommissionLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCommissions = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await apiClient.get<CommissionLog[]>(API_ENDPOINTS.settlement.commissions);
      setCommissions(res.data || []);
    } catch (err: any) {
      console.error(err);
      setError('Gagal memuat data komisi & settlement keuangan.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCommissions();
  }, []);

  // Format currency
  const formatRupiah = (num: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(num);
  };

  // Compute stats
  const totalEarnings = commissions.reduce((sum, item) => sum + item.amount, 0);
  const pendingAmount = commissions
    .filter(item => item.status === 'PENDING')
    .reduce((sum, item) => sum + item.amount, 0);
  const paidAmount = commissions
    .filter(item => item.status === 'PAID' || item.status === 'COMPLETED')
    .reduce((sum, item) => sum + item.amount, 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>💳 Keuangan & Settlement</h1>
        <p>Rekonsiliasi, pembayaran COD, dan penyelesaian tagihan komisi kurir</p>
      </div>

      {error && (
        <div style={{ padding: '12px 16px', background: '#fce8e6', color: '#c5221f', borderRadius: '8px', fontSize: '14px', fontWeight: 500 }}>
          {error}
        </div>
      )}

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center', background: 'rgba(121, 174, 111, 0.05)', borderColor: 'var(--primary)' }}>
          <div style={{ color: 'var(--primary-dark)', fontSize: '13px', fontWeight: 600 }}>Total Komisi Terhitung</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>
            {formatRupiah(totalEarnings)}
          </div>
          <div style={{ fontSize: '12px', color: 'var(--success)', marginTop: '4px' }}>Akumulasi keseluruhan rute</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Dana Belum Cair (Pending)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--warning)', marginTop: '8px' }}>
            {formatRupiah(pendingAmount)}
          </div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Menunggu approval admin</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Komisi Sudah Dibayarkan</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--info)', marginTop: '8px' }}>
            {formatRupiah(paidAmount)}
          </div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Telah ditransfer ke rekening kurir</div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Log Komisi Pengiriman Terbaru</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button className="btn btn-ghost" onClick={fetchCommissions}>Refresh</button>
          </div>
        </div>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>Memuat data...</div>
        ) : commissions.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>Belum ada log komisi terbaru.</div>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>ID Komisi</th>
                  <th>ID Kurir</th>
                  <th>Nomor Resi (AWB)</th>
                  <th>Nominal</th>
                  <th>Status</th>
                  <th>Tanggal Terbentuk</th>
                </tr>
              </thead>
              <tbody>
                {commissions.map((trx) => (
                  <tr key={trx.id}>
                    <td style={{ fontWeight: 500 }}>{trx.id.substring(0, 8).toUpperCase()}</td>
                    <td>{trx.courier_id}</td>
                    <td><code style={{ background: 'var(--bg-hover)', padding: '2px 6px', borderRadius: '4px' }}>{trx.awb}</code></td>
                    <td style={{ fontWeight: 600 }}>{formatRupiah(trx.amount)}</td>
                    <td>
                      <span className={`badge ${
                        trx.status === 'PAID' || trx.status === 'COMPLETED' ? 'badge-success' : 
                        trx.status === 'PENDING' ? 'badge-warning' : 'badge-danger'
                      }`}>
                        {trx.status}
                      </span>
                    </td>
                    <td>{new Date(trx.created_at || new Date()).toLocaleString('id-ID')}</td>
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
