// Mock data for presentation
const MOCK_SETTLEMENT = [
  { id: 'STL-2024-0991', date: '18 Jun 2024', amount: 'Rp 14.500.000', method: 'Bank Transfer (BCA)', status: 'COMPLETED' },
  { id: 'STL-2024-0992', date: '18 Jun 2024', amount: 'Rp 2.150.000', method: 'E-Wallet (OVO)', status: 'PENDING' },
  { id: 'STL-2024-0993', date: '17 Jun 2024', amount: 'Rp 8.900.000', method: 'Virtual Account (Mandiri)', status: 'COMPLETED' },
  { id: 'STL-2024-0994', date: '16 Jun 2024', amount: 'Rp 450.000', method: 'COD Collection', status: 'FAILED' },
];

export default function SettlementPage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>💳 Keuangan & Settlement</h1>
        <p>Rekonsiliasi, pembayaran COD, dan penyelesaian tagihan</p>
      </div>

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center', background: 'rgba(121, 174, 111, 0.05)', borderColor: 'var(--primary)' }}>
          <div style={{ color: 'var(--primary-dark)', fontSize: '13px', fontWeight: 600 }}>Total Pendapatan (Bulan Ini)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>Rp 124.5M</div>
          <div style={{ fontSize: '12px', color: 'var(--success)', marginTop: '4px' }}>↑ 8.4% vs bulan lalu</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Dana Belum Cair (Pending)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--warning)', marginTop: '8px' }}>Rp 4.2M</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Proses kliring bank (T+1)</div>
        </div>
        <div className="card" style={{ padding: '20px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Setoran COD Hari Ini</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--info)', marginTop: '8px' }}>Rp 18.5M</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Dari 142 kurir aktif</div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Riwayat Transaksi Terbaru</h3>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button className="btn btn-ghost">Download Report</button>
            <button className="btn btn-primary">Tarik Dana</button>
          </div>
        </div>

        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>ID Transaksi</th>
                <th>Tanggal</th>
                <th>Metode Pembayaran</th>
                <th>Nominal</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {MOCK_SETTLEMENT.map(trx => (
                <tr key={trx.id}>
                  <td style={{ fontWeight: 500 }}>{trx.id}</td>
                  <td>{trx.date}</td>
                  <td>{trx.method}</td>
                  <td style={{ fontWeight: 600 }}>{trx.amount}</td>
                  <td>
                    <span className={`badge ${
                      trx.status === 'COMPLETED' ? 'badge-success' : 
                      trx.status === 'PENDING' ? 'badge-warning' : 'badge-danger'
                    }`}>
                      {trx.status}
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
