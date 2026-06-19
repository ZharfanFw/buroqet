import { useState, useEffect } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';

interface PackageActivity {
  id: string;
  awb: string;
  status: string;
  location: string;
  operator: string;
  created_at: string;
}

export default function WarehousePage() {
  const [activities, setActivities] = useState<PackageActivity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [scanType, setScanType] = useState<'INBOUND' | 'OUTBOUND'>('INBOUND');
  const [formData, setFormData] = useState({
    awb: '',
    location: '',
    operator: '',
    destination: ''
  });
  const [submitting, setSubmitting] = useState(false);

  const fetchActivities = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await apiClient.get<PackageActivity[]>(API_ENDPOINTS.warehouse.packages);
      // Sort by newest
      const sorted = (res.data || []).sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );
      setActivities(sorted);
    } catch (err: any) {
      console.error(err);
      setError('Gagal memuat aktivitas inventaris gudang.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchActivities();
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.awb || !formData.location || !formData.operator) {
      setError('AWB, Lokasi, dan Operator harus diisi.');
      return;
    }

    try {
      setSubmitting(true);
      setError(null);
      setSuccessMsg(null);

      if (scanType === 'INBOUND') {
        await apiClient.post(API_ENDPOINTS.warehouse.inbound, {
          awb: formData.awb,
          location: formData.location,
          operator: formData.operator
        });
        setSuccessMsg(`Sukses melakukan Inbound untuk AWB ${formData.awb}`);
      } else {
        await apiClient.post(API_ENDPOINTS.warehouse.dispatch, {
          awb: formData.awb,
          location: formData.location,
          operator: formData.operator,
          destination: formData.destination || 'Hub Tujuan'
        });
        setSuccessMsg(`Sukses melakukan Outbound/Dispatch untuk AWB ${formData.awb}`);
      }

      // Reset Form & Close Modal
      setFormData({ awb: '', location: '', operator: '', destination: '' });
      setIsModalOpen(false);
      fetchActivities();
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.error || 'Gagal memproses pemindaian barcode.');
    } finally {
      setSubmitting(false);
    }
  };

  // Count stats
  const totalInbound = activities.filter(a => a.status === 'INBOUND').length;
  const totalOutbound = activities.filter(a => a.status === 'DISPATCHED' || a.status === 'OUTBOUND').length;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>🏭 Manajemen Gudang</h1>
        <p>Pantau arus barang masuk (inbound) dan keluar (outbound) secara real-time</p>
      </div>

      {successMsg && (
        <div style={{ padding: '12px 16px', background: '#e6f4ea', color: '#137333', borderRadius: '8px', fontSize: '14px', fontWeight: 500 }}>
          {successMsg}
        </div>
      )}

      {error && (
        <div style={{ padding: '12px 16px', background: '#fce8e6', color: '#c5221f', borderRadius: '8px', fontSize: '14px', fontWeight: 500 }}>
          {error}
        </div>
      )}

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Kapasitas Gudang</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>
            {Math.min(100, Math.max(10, activities.length * 5))}%
          </div>
          <div style={{ width: '100%', height: '6px', background: 'var(--bg-hover)', borderRadius: '3px', marginTop: '8px', overflow: 'hidden' }}>
            <div style={{ width: `${Math.min(100, Math.max(10, activities.length * 5))}%`, height: '100%', background: 'var(--primary)' }}></div>
          </div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Inbound</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--info)', marginTop: '8px' }}>{totalInbound}</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Paket diterima di hub</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Outbound</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--success)', marginTop: '8px' }}>{totalOutbound}</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Paket diberangkatkan</div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Aktivitas Pemindaian Terakhir</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button className="btn btn-ghost" onClick={fetchActivities}>Refresh</button>
            <button className="btn btn-primary" onClick={() => { setIsModalOpen(true); setScanType('INBOUND'); }}>Scan Inbound</button>
            <button className="btn btn-primary" style={{ backgroundColor: 'var(--info)' }} onClick={() => { setIsModalOpen(true); setScanType('OUTBOUND'); }}>Scan Outbound</button>
          </div>
        </div>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>Memuat data...</div>
        ) : activities.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>Belum ada aktivitas inventaris gudang.</div>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>ID Paket</th>
                  <th>Jenis Aktivitas</th>
                  <th>Resi (AWB)</th>
                  <th>Lokasi Rack/Gate</th>
                  <th>Waktu</th>
                  <th>Operator</th>
                </tr>
              </thead>
              <tbody>
                {activities.map((inv) => (
                  <tr key={inv.id}>
                    <td style={{ fontWeight: 500 }}>{inv.id.substring(0, 8).toUpperCase()}</td>
                    <td>
                      <span className={`badge ${inv.status === 'INBOUND' ? 'badge-info' : 'badge-success'}`}>
                        {inv.status === 'INBOUND' ? '📥 INBOUND' : '📤 OUTBOUND'}
                      </span>
                    </td>
                    <td><code style={{ background: 'var(--bg-hover)', padding: '2px 6px', borderRadius: '4px' }}>{inv.awb}</code></td>
                    <td>{inv.location}</td>
                    <td>{new Date(inv.created_at).toLocaleString('id-ID')}</td>
                    <td>{inv.operator}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Interactive Modal Scan */}
      {isModalOpen && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 999
        }}>
          <div className="card" style={{ width: '400px', padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <h3 style={{ fontSize: '18px', fontWeight: 600 }}>
              {scanType === 'INBOUND' ? '📥 Scan Barang Masuk (Inbound)' : '📤 Scan Barang Keluar (Outbound)'}
            </h3>

            <form onSubmit={handleFormSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)' }}>Nomor Resi (AWB)</label>
                <input
                  type="text"
                  name="awb"
                  value={formData.awb}
                  onChange={handleInputChange}
                  placeholder="Contoh: BQ-123456"
                  required
                  style={{ padding: '8px 12px', border: '1px solid var(--border-color)', borderRadius: '6px' }}
                />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)' }}>
                  {scanType === 'INBOUND' ? 'Lokasi Rak' : 'Gate Keberangkatan'}
                </label>
                <input
                  type="text"
                  name="location"
                  value={formData.location}
                  onChange={handleInputChange}
                  placeholder={scanType === 'INBOUND' ? 'Contoh: Rak A-12' : 'Contoh: Gate 3'}
                  required
                  style={{ padding: '8px 12px', border: '1px solid var(--border-color)', borderRadius: '6px' }}
                />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)' }}>Nama Operator</label>
                <input
                  type="text"
                  name="operator"
                  value={formData.operator}
                  onChange={handleInputChange}
                  placeholder="Masukkan nama Anda"
                  required
                  style={{ padding: '8px 12px', border: '1px solid var(--border-color)', borderRadius: '6px' }}
                />
              </div>

              {scanType === 'OUTBOUND' && (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                  <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)' }}>Hub Tujuan</label>
                  <input
                    type="text"
                    name="destination"
                    value={formData.destination}
                    onChange={handleInputChange}
                    placeholder="Contoh: Hub Bandung"
                    style={{ padding: '8px 12px', border: '1px solid var(--border-color)', borderRadius: '6px' }}
                  />
                </div>
              )}

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '12px' }}>
                <button type="button" className="btn btn-ghost" onClick={() => setIsModalOpen(false)}>Batal</button>
                <button type="submit" className="btn btn-primary" disabled={submitting}>
                  {submitting ? 'Memproses...' : 'Submit Scan'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
