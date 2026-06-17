// TODO: Implement TrackingPage — connect to TRACKING service
export default function TrackingPage() {
  return (
    <div>
      <div className="page-header">
        <h1>📦 Tracking</h1>
        <p>Lacak status pengiriman paket secara real-time</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">📦</div>
          <p>Halaman Tracking dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/tracking</code>
          </p>
        </div>
      </div>
    </div>
  );
}
