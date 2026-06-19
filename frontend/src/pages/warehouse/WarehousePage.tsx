// Mock data for presentation
const MOCK_INVENTORY = [
  { id: 'PKG-8821', type: 'INBOUND', awb: 'BQ-2024-JKT-001', location: 'Rack A-12', time: '10:42 AM', operator: 'Ahmad' },
  { id: 'PKG-8822', type: 'OUTBOUND', awb: 'BQ-2024-SBY-042', location: 'Gate 3', time: '10:15 AM', operator: 'Siti' },
  { id: 'PKG-8823', type: 'INBOUND', awb: 'BQ-2024-MDN-007', location: 'Rack B-05', time: '09:30 AM', operator: 'Ahmad' },
  { id: 'PKG-8824', type: 'OUTBOUND', awb: 'BQ-2024-BENTO-123', location: 'Gate 1', time: '08:45 AM', operator: 'Joko' },
];

export default function WarehousePage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>🏭 Manajemen Gudang</h1>
        <p>Pantau arus barang masuk (inbound) dan keluar (outbound) secara real-time</p>
      </div>

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Kapasitas Gudang</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>78%</div>
          <div style={{ width: '100%', height: '6px', background: 'var(--bg-hover)', borderRadius: '3px', marginTop: '8px', overflow: 'hidden' }}>
            <div style={{ width: '78%', height: '100%', background: 'var(--primary)' }}></div>
          </div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Inbound (Hari Ini)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--info)', marginTop: '8px' }}>1,402</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Paket diterima di hub</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Outbound (Hari Ini)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--success)', marginTop: '8px' }}>1,385</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Paket diberangkatkan</div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Aktivitas Pemindaian Terakhir</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button className="btn btn-ghost">Export Log</button>
            <button className="btn btn-primary">Scan Barcode</button>
          </div>
        </div>

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
              {MOCK_INVENTORY.map(inv => (
                <tr key={inv.id}>
                  <td style={{ fontWeight: 500 }}>{inv.id}</td>
                  <td>
                    <span className={`badge ${inv.type === 'INBOUND' ? 'badge-info' : 'badge-success'}`}>
                      {inv.type === 'INBOUND' ? '📥 INBOUND' : '📤 OUTBOUND'}
                    </span>
                  </td>
                  <td><code style={{ background: 'var(--bg-hover)', padding: '2px 6px', borderRadius: '4px' }}>{inv.awb}</code></td>
                  <td>{inv.location}</td>
                  <td>{inv.time}</td>
                  <td>{inv.operator}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
