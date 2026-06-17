// TODO: Implement SettlementPage — connect to SETTLEMENT service
export default function SettlementPage() {
  return (
    <div>
      <div className="page-header">
        <h1>💳 Settlement</h1>
        <p>Rekonsiliasi dan penyelesaian pembayaran</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">💳</div>
          <p>Halaman Settlement dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/settlement</code>
          </p>
        </div>
      </div>
    </div>
  );
}
