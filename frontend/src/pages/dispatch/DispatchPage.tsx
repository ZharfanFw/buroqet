// TODO: Implement DispatchPage — connect to DISPATCH service
export default function DispatchPage() {
  return (
    <div>
      <div className="page-header">
        <h1>🚚 Dispatch</h1>
        <p>Manajemen armada dan penugasan kurir</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">🚚</div>
          <p>Halaman Dispatch dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/dispatch</code>
          </p>
        </div>
      </div>
    </div>
  );
}
