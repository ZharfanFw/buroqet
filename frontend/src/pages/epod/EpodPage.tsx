// Mock data for presentation
const MOCK_EPOD = [
  { id: 'POD-1021', awb: 'BQ-2024-JKT-001', receiver: 'Satpam (Bpk. Yanto)', date: '18 Jun, 14:30', status: 'VERIFIED' },
  { id: 'POD-1022', awb: 'BQ-2024-BENTO-123', receiver: 'Budi Santoso', date: '17 Jun, 10:15', status: 'VERIFIED' },
  { id: 'POD-1023', awb: 'BQ-2024-SBY-042', receiver: 'Ibu RT', date: '16 Jun, 16:45', status: 'PENDING_REVIEW' },
  { id: 'POD-1024', awb: 'BQ-2024-MDN-007', receiver: 'Tidak Diketahui', date: '15 Jun, 11:20', status: 'REJECTED' },
];

export default function EpodPage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>✍️ Electronic Proof of Delivery (ePOD)</h1>
        <p>Verifikasi bukti pengiriman berupa foto dan tanda tangan penerima</p>
      </div>

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total ePOD Hari Ini</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>845</div>
          <div style={{ fontSize: '12px', color: 'var(--success)', marginTop: '4px' }}>100% tersinkronisasi</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Menunggu Verifikasi (Pending)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--warning)', marginTop: '8px' }}>12</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Harap ditinjau oleh Admin</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Tingkat Keberhasilan Verifikasi</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--primary)', marginTop: '8px' }}>99.8%</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Bulan berjalan</div>
        </div>
      </div>

      {/* Gallery Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Bukti Pengiriman Terbaru</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <input type="text" className="input" placeholder="Cari Resi (AWB)..." style={{ width: '200px' }} />
            <button className="btn btn-ghost">Filter Status</button>
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: '16px', marginTop: '10px' }}>
          {MOCK_EPOD.map(pod => (
            <div key={pod.id} style={{ border: '1px solid var(--border)', borderRadius: '12px', overflow: 'hidden' }}>
              <div style={{ height: '140px', background: 'var(--bg-hover)', display: 'flex', alignItems: 'center', justifyContent: 'center', position: 'relative' }}>
                <span style={{ fontSize: '40px', opacity: 0.2 }}>📸</span>
                <div style={{ position: 'absolute', bottom: '8px', right: '8px', background: 'rgba(0,0,0,0.6)', color: 'white', fontSize: '10px', padding: '2px 6px', borderRadius: '4px' }}>
                  Lat: -6.2, Lng: 106.8
                </div>
              </div>
              <div style={{ padding: '12px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '8px' }}>
                  <code style={{ fontSize: '11px', fontWeight: 600, color: 'var(--primary-dark)' }}>{pod.awb}</code>
                  <span className={`badge ${
                    pod.status === 'VERIFIED' ? 'badge-success' : 
                    pod.status === 'PENDING_REVIEW' ? 'badge-warning' : 'badge-danger'
                  }`} style={{ fontSize: '9px', padding: '2px 6px' }}>
                    {pod.status}
                  </span>
                </div>
                <div style={{ fontSize: '13px', fontWeight: 500, marginBottom: '2px' }}>{pod.receiver}</div>
                <div style={{ fontSize: '11px', color: 'var(--text-muted)' }}>{pod.date}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
