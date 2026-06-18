import { useState } from 'react';

// Mock data for presentation
const MOCK_ORDERS = [
  { id: 'ORD-2024-001', awb: 'BQ-2024-JKT-001', customer: 'PT. Teknologi Maju', destination: 'Jakarta Pusat', status: 'ON_TRANSIT', date: '2024-06-18' },
  { id: 'ORD-2024-002', awb: 'BQ-2024-BENTO-123', customer: 'Budi Santoso', destination: 'Bandung', status: 'DELIVERED', date: '2024-06-17' },
  { id: 'ORD-2024-003', awb: 'BQ-2024-SBY-042', customer: 'CV. Karya Abadi', destination: 'Surabaya', status: 'INBOUND', date: '2024-06-18' },
  { id: 'ORD-2024-004', awb: 'BQ-2024-MDN-007', customer: 'Siti Aminah', destination: 'Medan', status: 'FAILED', date: '2024-06-15' },
];

export default function OrderPage() {
  const [search, setSearch] = useState('');

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>🧾 Manajemen Order</h1>
        <p>Kelola dan pantau semua order pengiriman di sistem Buroqet</p>
      </div>

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Order (Bulan Ini)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>1,248</div>
          <div style={{ fontSize: '12px', color: 'var(--success)', marginTop: '4px' }}>↑ 12% dari bulan lalu</div>
        </div>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Menunggu Penjemputan</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--warning)', marginTop: '8px' }}>42</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Butuh tindakan segera</div>
        </div>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Selesai Terkirim</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--success)', marginTop: '8px' }}>1,102</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Success rate: 98.2%</div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Daftar Order Terbaru</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <input 
              type="text" 
              className="input" 
              placeholder="Cari order atau AWB..." 
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              style={{ width: '250px' }}
            />
            <button className="btn btn-primary">Buat Order Baru</button>
          </div>
        </div>

        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Order ID</th>
                <th>Resi (AWB)</th>
                <th>Pelanggan</th>
                <th>Tujuan</th>
                <th>Tanggal</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {MOCK_ORDERS.map(order => (
                <tr key={order.id}>
                  <td style={{ fontWeight: 500 }}>{order.id}</td>
                  <td><code style={{ background: 'var(--bg-hover)', padding: '2px 6px', borderRadius: '4px' }}>{order.awb}</code></td>
                  <td>{order.customer}</td>
                  <td>{order.destination}</td>
                  <td>{order.date}</td>
                  <td>
                    <span className={`badge ${
                      order.status === 'DELIVERED' ? 'badge-success' : 
                      order.status === 'FAILED' ? 'badge-danger' : 
                      order.status === 'ON_TRANSIT' ? 'badge-warning' : 'badge-info'
                    }`}>
                      {order.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
