// TODO: Implement WarehousePage — connect to WAREHOUSE service
export default function WarehousePage() {
  return (
    <div>
      <div className="page-header">
        <h1>🏭 Warehouse</h1>
        <p>Manajemen inventori dan gudang</p>
      </div>
      <div className="card">
        <div className="empty-state">
          <div className="empty-icon">🏭</div>
          <p>Halaman Warehouse dalam pengembangan.</p>
          <p style={{ fontSize: '12px', marginTop: '8px', opacity: 0.6 }}>
            Hubungkan ke endpoint: <code>/api/warehouse</code>
          </p>
        </div>
      </div>
    </div>
  );
}
