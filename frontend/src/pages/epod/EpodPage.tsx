// TODO: Implement EpodPage — connect to EPOD service
export default function EpodPage() {
  return (
    <div>
      <div className="page-header">
        <h1>✍️ Epod</h1>
        <p>Electronic Proof of Delivery</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">✍️</div>
          <p>Halaman Epod dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/epod</code>
          </p>
        </div>
      </div>
    </div>
  );
}
