// TODO: Implement PricingPage — connect to PRICING service
export default function PricingPage() {
  return (
    <div>
      <div className="page-header">
        <h1>💰 Pricing</h1>
        <p>Kalkulasi harga ongkos kirim</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">💰</div>
          <p>Halaman Pricing dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/pricing</code>
          </p>
        </div>
      </div>
    </div>
  );
}
