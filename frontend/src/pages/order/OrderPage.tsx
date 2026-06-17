// TODO: Implement OrderPage — connect to ORDER service
export default function OrderPage() {
  return (
    <div>
      <div className="page-header">
        <h1>🧾 Order</h1>
        <p>Kelola dan pantau semua order pengiriman</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">🧾</div>
          <p>Halaman Order dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/order</code>
          </p>
        </div>
      </div>
    </div>
  );
}
