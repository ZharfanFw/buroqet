import { useState } from 'react';

export default function PricingPage() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(false);

  const handleCalculate = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      setResult(true);
    }, 800);
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>💰 Cek Ongkir (Pricing)</h1>
        <p>Kalkulasi estimasi biaya pengiriman berdasarkan berat dan dimensi</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        {/* Calculator Form */}
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '4px' }}>Detail Pengiriman</h3>
          
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
            <div>
              <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '6px' }}>Kota Asal</label>
              <select className="input">
                <option>Jakarta</option>
                <option>Bandung</option>
                <option>Surabaya</option>
              </select>
            </div>
            <div>
              <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '6px' }}>Kota Tujuan</label>
              <select className="input">
                <option>Surabaya</option>
                <option>Bandung</option>
                <option>Jakarta</option>
              </select>
            </div>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
            <div>
              <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '6px' }}>Berat (kg)</label>
              <input type="number" className="input" defaultValue={1} min={1} />
            </div>
            <div>
              <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '6px' }}>Dimensi (PxLxT) cm</label>
              <input type="text" className="input" placeholder="Contoh: 10x10x10" />
            </div>
          </div>

          <button 
            className="btn btn-primary" 
            style={{ width: '100%', justifyContent: 'center', marginTop: '10px' }}
            onClick={handleCalculate}
            disabled={loading}
          >
            {loading ? 'Menghitung...' : 'Hitung Ongkos Kirim'}
          </button>
        </div>

        {/* Results */}
        <div className="card" style={{ display: 'flex', flexDirection: 'column', background: result ? 'var(--bg-surface)' : 'var(--bg-hover)' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Hasil Kalkulasi</h3>
          
          {!result && !loading && (
            <div className="empty-state" style={{ padding: '40px 20px', margin: 'auto' }}>
              <div className="empty-icon" style={{ opacity: 0.3 }}>🧾</div>
              <p>Masukkan detail pengiriman untuk melihat estimasi harga.</p>
            </div>
          )}

          {loading && (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', flex: 1 }}>
              <div className="spinner" style={{ borderColor: 'var(--border)', borderTopColor: 'var(--primary)' }}></div>
            </div>
          )}

          {result && !loading && (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <div style={{ border: '1px solid var(--primary)', borderRadius: '12px', padding: '16px', background: 'rgba(121, 174, 111, 0.05)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                  <div style={{ fontWeight: 600, color: 'var(--primary-dark)' }}>Reguler (REG)</div>
                  <div style={{ fontSize: '18px', fontWeight: 700 }}>Rp 15.000</div>
                </div>
                <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Estimasi tiba: 2-3 Hari</div>
              </div>
              
              <div style={{ border: '1px solid var(--border)', borderRadius: '12px', padding: '16px', background: 'var(--bg-surface)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                  <div style={{ fontWeight: 600 }}>Next Day (YES)</div>
                  <div style={{ fontSize: '18px', fontWeight: 700 }}>Rp 22.000</div>
                </div>
                <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Estimasi tiba: 1 Hari (Besok)</div>
              </div>
              
              <div style={{ border: '1px solid var(--border)', borderRadius: '12px', padding: '16px', background: 'var(--bg-surface)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                  <div style={{ fontWeight: 600 }}>Kargo (CARGO)</div>
                  <div style={{ fontSize: '18px', fontWeight: 700 }}>Rp 45.000</div>
                </div>
                <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>Estimasi tiba: 4-7 Hari • Min 10kg</div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
