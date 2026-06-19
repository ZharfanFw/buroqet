import { useState, useEffect } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';

interface BackendOrder {
  order_id: string;
  awb_number: string;
  customer_id?: string;
  service_id?: string;
  sender_name: string;
  sender_phone: string;
  sender_address: string;
  origin_postal_code: string;
  origin_city?: string;
  origin_lat?: number;
  origin_lng?: number;
  receiver_name: string;
  receiver_phone: string;
  receiver_address: string;
  dest_postal_code: string;
  dest_city?: string;
  dest_lat?: number;
  dest_lng?: number;
  actual_weight_kg: number;
  length_cm: number;
  width_cm: number;
  height_cm: number;
  pricing_method?: string;
  base_tariff: number;
  insurance_fee: number;
  discount: number;
  total_cost: number;
  use_insurance: boolean;
  payment_type: 'COD' | 'TRANSFER' | 'EWALLET' | 'VA';
  payment_ref?: string;
  status: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

const AVAILABLE_STATUSES = [
  'ORDER_CREATED',
  'PAYMENT_PENDING',
  'PAYMENT_CONFIRMED',
  'PICKED_UP',
  'ON_TRANSIT',
  'AT_DESTINATION_HUB',
  'OUT_FOR_DELIVERY',
  'DELIVERED',
  'FAILED',
  'RETURNED'
];

export default function OrderPage() {
  const [orders, setOrders] = useState<BackendOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  
  // Modals state
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<BackendOrder | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [actionError, setActionError] = useState('');

  // Create Form state
  const [form, setForm] = useState({
    sender_name: '',
    sender_phone: '',
    sender_address: '',
    origin_postal: '',
    origin_city: 'Jakarta',
    origin_lat: -6.2088,
    origin_lng: 106.8456,
    receiver_name: '',
    receiver_phone: '',
    receiver_address: '',
    dest_postal: '',
    dest_city: 'Bandung',
    dest_lat: -6.9175,
    dest_lng: 107.6191,
    weight_actual: 1.0,
    length: 10,
    width: 10,
    height: 10,
    service_type: 'REG',
    payment_type: 'COD',
    use_insurance: false
  });

  const fetchOrders = async () => {
    setLoading(true);
    setError("");
    try {
      const response = await apiClient.get(API_ENDPOINTS.order.base);

      // 💡 Tambahkan log ini untuk melihat bentuk asli dari Golang di Inspect Element -> Console
      console.log("📦 RAW Response dari Backend:", response.data);

      // Cek berbagai kemungkinan format JSON dari Golang
      if (Array.isArray(response.data)) {
        // Jika Golang mengembalikan array langsung: [...]
        setOrders(response.data);
      } else if (response.data && response.data.data) {
        // Jika Golang mengembalikan format: { data: [...] } atau { data: { orders: [...] } }
        const dataPayload = response.data.data;

        if (Array.isArray(dataPayload)) {
          setOrders(dataPayload);
        } else if (dataPayload.orders && Array.isArray(dataPayload.orders)) {
          setOrders(dataPayload.orders);
        } else {
          setOrders([]); // Fallback kalau datanya kosong
        }
      } else {
        // Jika formatnya benar-benar tidak dikenali
        setError(
          "Format respons server tidak valid. Cek Console (F12) untuk detailnya.",
        );
      }
    } catch (err: any) {
      console.error("Gagal fetch:", err);
      setError("Gagal memuat data order. Pastikan backend server berjalan.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const filteredOrders = orders.filter(order => {
    const searchLower = search.toLowerCase();
    return (
      order.awb_number.toLowerCase().includes(searchLower) ||
      order.sender_name.toLowerCase().includes(searchLower) ||
      order.receiver_name.toLowerCase().includes(searchLower) ||
      order.order_id.toLowerCase().includes(searchLower)
    );
  });

  // Stats
  const totalOrdersCount = orders.length;
  const pendingPickupCount = orders.filter(o => 
    o.status === 'ORDER_CREATED' || o.status === 'PAYMENT_PENDING' || o.status === 'PAYMENT_CONFIRMED'
  ).length;
  const deliveredCount = orders.filter(o => o.status === 'DELIVERED').length;

  const handleCreateOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    setActionLoading(true);
    setActionError("");
    try {
      // 🛠️ MAPPING PAYLOAD:
      // Sesuaikan nama variabel dari form state ke nama variabel yang diminta Backend Golang
      const payload = {
        sender_name: form.sender_name,
        sender_phone: form.sender_phone,
        sender_address: form.sender_address,
        origin_postal_code: form.origin_postal,
        origin_city: form.origin_city,
        origin_lat: form.origin_lat,
        origin_lng: form.origin_lng,
        receiver_name: form.receiver_name,
        receiver_phone: form.receiver_phone,
        receiver_address: form.receiver_address,
        dest_postal_code: form.dest_postal,
        dest_city: form.dest_city,
        dest_lat: form.dest_lat,
        dest_lng: form.dest_lng,
        actual_weight_kg: Number(form.weight_actual),
        length_cm: Number(form.length),
        width_cm: Number(form.width),
        height_cm: Number(form.height),
        service_id: form.service_type, // Menyimpan nilai 'REG' / 'EXP' dsb
        payment_type: form.payment_type,
        use_insurance: form.use_insurance,
      };

      const response = await apiClient.post(API_ENDPOINTS.order.base, payload);

      if (response.data && response.data.success) {
        // 1. Kasih notifikasi sukses! Ambil AWB / Order ID dari respons backend
        const newResi =
          response.data.data?.awb_number ||
          response.data.data?.order_id ||
          "Sukses";
        alert(`🎉 Yey! Order berhasil dibuat dengan Resi: ${newResi}`);

        // 2. Tutup Modal
        setShowCreateModal(false);

        // 3. Reset form ke default
        setForm({
          sender_name: "",
          sender_phone: "",
          sender_address: "",
          origin_postal: "",
          origin_city: "Jakarta",
          origin_lat: -6.2088,
          origin_lng: 106.8456,
          receiver_name: "",
          receiver_phone: "",
          receiver_address: "",
          dest_postal: "",
          dest_city: "Bandung",
          dest_lat: -6.9175,
          dest_lng: 107.6191,
          weight_actual: 1.0,
          length: 10,
          width: 10,
          height: 10,
          service_type: "REG",
          payment_type: "COD",
          use_insurance: false,
        });

        // 4. Refresh data tabel
        fetchOrders();
      } else {
        setActionError("Gagal membuat order.");
      }
    } catch (err: any) {
      console.error(err);
      // Tangkap pesan eror spesifik dari Golang jika ada
      const serverMessage =
        err.response?.data?.message || err.response?.data?.error;
      setActionError(
        serverMessage || "Gagal membuat order. Cek kelengkapan data.",
      );
    } finally {
      setActionLoading(false);
    }
  };

  const handleUpdateStatus = async (status: string) => {
    if (!selectedOrder) return;
    setActionLoading(true);
    setActionError('');
    try {
      const response = await apiClient.patch(API_ENDPOINTS.order.status(selectedOrder.awb_number), {
        status,
        notes: `Status diupdate via dashboard frontend ke ${status}`
      });
      if (response.data && response.data.success) {
        setSelectedOrder(response.data.data);
        fetchOrders();
      }
    } catch (err: any) {
      console.error(err);
      setActionError(err.response?.data?.error || 'Gagal memperbarui status.');
    } finally {
      setActionLoading(false);
    }
  };

  const handleCancelOrder = async () => {
    if (!selectedOrder) return;
    if (!window.confirm(`Apakah Anda yakin ingin membatalkan order ${selectedOrder.awb_number}?`)) return;
    
    setActionLoading(true);
    setActionError('');
    try {
      const response = await apiClient.delete(API_ENDPOINTS.order.byId(selectedOrder.awb_number));
      if (response.data && response.data.success) {
        setSelectedOrder(null);
        fetchOrders();
      }
    } catch (err: any) {
      console.error(err);
      setActionError(err.response?.data?.error || 'Gagal membatalkan order. Order hanya bisa dibatalkan jika statusnya masih ORDER_CREATED.');
    } finally {
      setActionLoading(false);
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'DELIVERED': return 'badge-success';
      case 'FAILED': return 'badge-danger';
      case 'ORDER_CREATED':
      case 'PAYMENT_PENDING':
      case 'PAYMENT_CONFIRMED': return 'badge-info';
      case 'PICKED_UP':
      case 'ON_TRANSIT':
      case 'AT_DESTINATION_HUB':
      case 'OUT_FOR_DELIVERY': return 'badge-warning';
      default: return 'badge-default';
    }
  };

  

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>🧾 Manajemen Order</h1>
        <p>Kelola dan pantau semua order pengiriman di sistem Buroqet</p>
      </div>

      {/* Bento Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Total Order (Semua)</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '8px' }}>{totalOrdersCount}</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Jumlah order dalam database</div>
        </div>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Menunggu / Diproses</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--warning)', marginTop: '8px' }}>{pendingPickupCount}</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>Butuh penjemputan & proses</div>
        </div>
        <div className="card" style={{ padding: '20px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '13px', fontWeight: 600 }}>Selesai Terkirim</div>
          <div style={{ fontSize: '32px', fontWeight: 700, color: 'var(--success)', marginTop: '8px' }}>{deliveredCount}</div>
          <div style={{ fontSize: '12px', color: 'var(--text-secondary)', marginTop: '4px' }}>
            Success rate: {totalOrdersCount > 0 ? ((deliveredCount / totalOrdersCount) * 100).toFixed(1) : '0'}%
          </div>
        </div>
      </div>

      {/* Table Section */}
      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px' }}>
          <h3 style={{ fontSize: '16px', fontWeight: 600 }}>Daftar Order Terbaru</h3>
          <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
            <input 
              type="text" 
              className="input" 
              placeholder="Cari pengirim, penerima, AWB..." 
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              style={{ width: '250px' }}
            />
            <button className="btn btn-primary" onClick={() => setShowCreateModal(true)}>Buat Order Baru</button>
          </div>
        </div>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--text-muted)' }}>
            Memuat data order...
          </div>
        ) : error ? (
          <div style={{ textAlign: 'center', padding: '40px 0', color: 'var(--danger)' }}>
            ⚠️ {error}
            <div style={{ marginTop: '10px' }}>
              <button className="btn btn-ghost" onClick={fetchOrders}>Coba Lagi</button>
            </div>
          </div>
        ) : filteredOrders.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📦</div>
            <h3>Belum ada order</h3>
            <p>Tidak ada order yang ditemukan. Buat order baru untuk memulai pengujian.</p>
          </div>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Order ID</th>
                  <th>Resi (AWB)</th>
                  <th>Pengirim</th>
                  <th>Penerima</th>
                  <th>Tujuan</th>
                  <th>Total Biaya</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {filteredOrders.map(order => (
                  <tr 
                    key={order.order_id} 
                    onClick={() => { setSelectedOrder(order); setActionError(''); }}
                    style={{ cursor: 'pointer' }}
                  >
                    <td style={{ fontWeight: 500, fontSize: '12px' }}>{order.order_id.substring(0, 8)}...</td>
                    <td>
                      <code style={{ background: 'var(--bg-hover)', padding: '2px 6px', borderRadius: '4px', fontSize: '12px' }}>
                        {order.awb_number}
                      </code>
                    </td>
                    <td>{order.sender_name}</td>
                    <td>{order.receiver_name}</td>
                    <td>{order.receiver_address.substring(0, 20)}...</td>
                    <td style={{ fontWeight: 600 }}>Rp {order.total_cost.toLocaleString()}</td>
                    <td>
                      <span className={`badge ${getStatusBadgeClass(order.status)}`}>
                        {order.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* CREATE ORDER MODAL */}
      {showCreateModal && (
        <div style={modalOverlayStyle}>
          <div style={modalContentStyle}>
            <button style={closeButtonStyle} onClick={() => setShowCreateModal(false)}>✕</button>
            <h2 style={{ marginBottom: '20px', fontWeight: 700, fontSize: '20px' }}>📦 Buat Order Baru</h2>
            
            {actionError && (
              <div style={{ background: 'rgba(226,106,106,0.1)', color: 'var(--danger)', padding: '12px', borderRadius: '8px', marginBottom: '16px', fontSize: '13px' }}>
                {actionError}
              </div>
            )}

            <form onSubmit={handleCreateOrder} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              
              {/* SENDER SECTION */}
              <div style={sectionStyle}>
                <h4 style={sectionTitleStyle}>👤 Informasi Pengirim</h4>
                <div style={grid2Style}>
                  <div>
                    <label style={labelStyle}>Nama Pengirim *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.sender_name} 
                      onChange={e => setForm({...form, sender_name: e.target.value})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>No. Telepon Pengirim *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.sender_phone} 
                      onChange={e => setForm({...form, sender_phone: e.target.value})}
                    />
                  </div>
                </div>
                <div style={{ marginTop: '8px' }}>
                  <label style={labelStyle}>Alamat Pengirim *</label>
                  <textarea 
                    required className="input" style={{ minHeight: '60px', resize: 'vertical' }}
                    value={form.sender_address} 
                    onChange={e => setForm({...form, sender_address: e.target.value})}
                  />
                </div>
                <div style={grid2Style}>
                  <div>
                    <label style={labelStyle}>Kode Pos Asal *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.origin_postal} 
                      onChange={e => setForm({...form, origin_postal: e.target.value})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>Kota Asal</label>
                    <input 
                      type="text" className="input" 
                      value={form.origin_city} 
                      onChange={e => setForm({...form, origin_city: e.target.value})}
                    />
                  </div>
                </div>
              </div>

              {/* RECEIVER SECTION */}
              <div style={sectionStyle}>
                <h4 style={sectionTitleStyle}>📍 Informasi Penerima</h4>
                <div style={grid2Style}>
                  <div>
                    <label style={labelStyle}>Nama Penerima *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.receiver_name} 
                      onChange={e => setForm({...form, receiver_name: e.target.value})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>No. Telepon Penerima *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.receiver_phone} 
                      onChange={e => setForm({...form, receiver_phone: e.target.value})}
                    />
                  </div>
                </div>
                <div style={{ marginTop: '8px' }}>
                  <label style={labelStyle}>Alamat Penerima *</label>
                  <textarea 
                    required className="input" style={{ minHeight: '60px', resize: 'vertical' }}
                    value={form.receiver_address} 
                    onChange={e => setForm({...form, receiver_address: e.target.value})}
                  />
                </div>
                <div style={grid2Style}>
                  <div>
                    <label style={labelStyle}>Kode Pos Tujuan *</label>
                    <input 
                      type="text" required className="input" 
                      value={form.dest_postal} 
                      onChange={e => setForm({...form, dest_postal: e.target.value})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>Kota Tujuan</label>
                    <input 
                      type="text" className="input" 
                      value={form.dest_city} 
                      onChange={e => setForm({...form, dest_city: e.target.value})}
                    />
                  </div>
                </div>
              </div>

              {/* PACKAGE & SHIPPING DETAILS */}
              <div style={sectionStyle}>
                <h4 style={sectionTitleStyle}>📦 Rincian Paket & Pengiriman</h4>
                <div style={grid4Style}>
                  <div>
                    <label style={labelStyle}>Berat (kg) *</label>
                    <input 
                      type="number" step="0.1" min="0.1" required className="input" 
                      value={form.weight_actual} 
                      onChange={e => setForm({...form, weight_actual: parseFloat(e.target.value) || 0})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>P (cm) *</label>
                    <input 
                      type="number" min="1" required className="input" 
                      value={form.length} 
                      onChange={e => setForm({...form, length: parseInt(e.target.value) || 0})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>L (cm) *</label>
                    <input 
                      type="number" min="1" required className="input" 
                      value={form.width} 
                      onChange={e => setForm({...form, width: parseInt(e.target.value) || 0})}
                    />
                  </div>
                  <div>
                    <label style={labelStyle}>T (cm) *</label>
                    <input 
                      type="number" min="1" required className="input" 
                      value={form.height} 
                      onChange={e => setForm({...form, height: parseInt(e.target.value) || 0})}
                    />
                  </div>
                </div>
                
                <div style={grid2Style}>
                  <div>
                    <label style={labelStyle}>Tipe Layanan *</label>
                    <select 
                      className="input" 
                      value={form.service_type}
                      onChange={e => setForm({...form, service_type: e.target.value})}
                    >
                      <option value="REG">Regular (REG)</option>
                      <option value="EXP">Express (EXP)</option>
                      <option value="CARGO">Cargo (CARGO)</option>
                    </select>
                  </div>
                  <div>
                    <label style={labelStyle}>Metode Pembayaran *</label>
                    <select 
                      className="input" 
                      value={form.payment_type}
                      onChange={e => setForm({...form, payment_type: e.target.value})}
                    >
                      <option value="COD">COD</option>
                      <option value="TRANSFER">Transfer Bank</option>
                      <option value="EWALLET">E-Wallet</option>
                      <option value="VA">Virtual Account</option>
                    </select>
                  </div>
                </div>

                <div style={{ marginTop: '12px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <input 
                    type="checkbox" id="use_insurance" 
                    checked={form.use_insurance}
                    onChange={e => setForm({...form, use_insurance: e.target.checked})}
                    style={{ cursor: 'pointer', width: '16px', height: '16px' }}
                  />
                  <label htmlFor="use_insurance" style={{ cursor: 'pointer', fontWeight: 500 }}>Gunakan Asuransi Pengiriman</label>
                </div>
              </div>

              <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginTop: '10px' }}>
                <button type="button" className="btn btn-ghost" onClick={() => setShowCreateModal(false)} disabled={actionLoading}>Batal</button>
                <button type="submit" className="btn btn-primary" disabled={actionLoading}>
                  {actionLoading ? 'Menyimpan...' : 'Buat Order'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ORDER DETAIL MODAL */}
      {selectedOrder && (
        <div style={modalOverlayStyle}>
          <div style={modalContentStyle}>
            <button style={closeButtonStyle} onClick={() => setSelectedOrder(null)}>✕</button>
            <h2 style={{ marginBottom: '4px', fontWeight: 700, fontSize: '20px' }}>🧾 Detail Order</h2>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '20px' }}>
              <code style={{ background: 'var(--bg-hover)', padding: '3px 8px', borderRadius: '4px', fontSize: '13px', fontWeight: 600 }}>
                {selectedOrder.awb_number}
              </code>
              <span className={`badge ${getStatusBadgeClass(selectedOrder.status)}`}>
                {selectedOrder.status}
              </span>
            </div>

            {actionError && (
              <div style={{ background: 'rgba(226,106,106,0.1)', color: 'var(--danger)', padding: '12px', borderRadius: '8px', marginBottom: '16px', fontSize: '13px' }}>
                {actionError}
              </div>
            )}

            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', maxHeight: '60vh', overflowY: 'auto', paddingRight: '4px' }}>
              <div style={grid2Style}>
                <div style={sectionStyle}>
                  <h4 style={sectionTitleStyle}>👤 Pengirim</h4>
                  <div style={{ fontSize: '13px' }}>
                    <strong>{selectedOrder.sender_name}</strong>
                    <p style={{ color: 'var(--text-secondary)' }}>{selectedOrder.sender_phone}</p>
                    <p style={{ marginTop: '4px' }}>{selectedOrder.sender_address}</p>
                    <p style={{ color: 'var(--text-muted)', fontSize: '11px', marginTop: '4px' }}>
                      Kode Pos: {selectedOrder.origin_postal_code} {selectedOrder.origin_city && `| ${selectedOrder.origin_city}`}
                    </p>
                  </div>
                </div>
                <div style={sectionStyle}>
                  <h4 style={sectionTitleStyle}>📍 Penerima</h4>
                  <div style={{ fontSize: '13px' }}>
                    <strong>{selectedOrder.receiver_name}</strong>
                    <p style={{ color: 'var(--text-secondary)' }}>{selectedOrder.receiver_phone}</p>
                    <p style={{ marginTop: '4px' }}>{selectedOrder.receiver_address}</p>
                    <p style={{ color: 'var(--text-muted)', fontSize: '11px', marginTop: '4px' }}>
                      Kode Pos: {selectedOrder.dest_postal_code} {selectedOrder.dest_city && `| ${selectedOrder.dest_city}`}
                    </p>
                  </div>
                </div>
              </div>

              <div style={sectionStyle}>
                <h4 style={sectionTitleStyle}>📦 Detail Paket & Biaya</h4>
                <div style={grid3Style}>
                  <div>
                    <span style={{ color: 'var(--text-muted)', fontSize: '11px' }}>Dimensi & Berat</span>
                    <p style={{ fontWeight: 600, fontSize: '13px' }}>
                      {selectedOrder.actual_weight_kg} kg ({selectedOrder.length_cm}x{selectedOrder.width_cm}x{selectedOrder.height_cm} cm)
                    </p>
                  </div>
                  <div>
                    <span style={{ color: 'var(--text-muted)', fontSize: '11px' }}>Metode Bayar</span>
                    <p style={{ fontWeight: 600, fontSize: '13px' }}>{selectedOrder.payment_type} {selectedOrder.use_insurance && '(Insured)'}</p>
                  </div>
                  <div>
                    <span style={{ color: 'var(--text-muted)', fontSize: '11px' }}>ID Transaksi / Ref</span>
                    <p style={{ fontWeight: 500, fontSize: '12px' }}>{selectedOrder.payment_ref || '-'}</p>
                  </div>
                </div>

                <div style={{ marginTop: '16px', borderTop: '1px solid var(--border)', paddingTop: '12px' }}>
                  <div style={costRowStyle}>
                    <span>Tarif Dasar (Base)</span>
                    <span>Rp {selectedOrder.base_tariff.toLocaleString()}</span>
                  </div>
                  <div style={costRowStyle}>
                    <span>Biaya Asuransi</span>
                    <span>Rp {selectedOrder.insurance_fee.toLocaleString()}</span>
                  </div>
                  <div style={costRowStyle}>
                    <span>Diskon</span>
                    <span style={{ color: 'var(--danger)' }}>- Rp {selectedOrder.discount.toLocaleString()}</span>
                  </div>
                  <div style={{ ...costRowStyle, fontWeight: 700, fontSize: '15px', color: 'var(--primary)', marginTop: '8px', borderTop: '1px dashed var(--border)', paddingTop: '8px' }}>
                    <span>Total Biaya</span>
                    <span>Rp {selectedOrder.total_cost.toLocaleString()}</span>
                  </div>
                </div>
              </div>

              {/* ADMIN CONTROL SECTION FOR TESTING TRANSITIONS */}
              <div style={{ ...sectionStyle, border: '1px solid var(--primary-light)', background: 'rgba(121,174,111,0.05)' }}>
                <h4 style={{ ...sectionTitleStyle, color: 'var(--primary-dark)' }}>⚙️ Simulasi Status (Pengujian Backend)</h4>
                <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
                  <label style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text-secondary)', minWidth: '80px' }}>Ubah Status:</label>
                  <select 
                    className="input" 
                    value={selectedOrder.status}
                    onChange={e => handleUpdateStatus(e.target.value)}
                    style={{ flex: 1, height: '36px', padding: '0 10px' }}
                    disabled={actionLoading}
                  >
                    {AVAILABLE_STATUSES.map(status => (
                      <option key={status} value={status}>{status}</option>
                    ))}
                  </select>
                </div>
                <p style={{ fontSize: '11px', color: 'var(--text-muted)', marginTop: '6px' }}>
                  * Gunakan dropdown ini untuk memperbarui lifecycle status order langsung ke database lokal.
                </p>
              </div>

              <div style={{ display: 'flex', gap: '10px', justifyContent: 'space-between', marginTop: '10px', alignItems: 'center' }}>
                <div>
                  {selectedOrder.status === 'ORDER_CREATED' ? (
                    <button 
                      type="button" 
                      className="btn btn-ghost" 
                      onClick={handleCancelOrder}
                      disabled={actionLoading}
                      style={{ color: 'var(--danger)', borderColor: 'var(--danger)' }}
                    >
                      🗑️ Batalkan Order
                    </button>
                  ) : (
                    <span style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
                      * Hanya order bertipe <strong>ORDER_CREATED</strong> yang bisa dibatalkan.
                    </span>
                  )}
                </div>
                <button type="button" className="btn btn-ghost" onClick={() => setSelectedOrder(null)} disabled={actionLoading}>Tutup</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Inline Styles to avoid writing external CSS file and guarantee compatibility
const modalOverlayStyle: React.CSSProperties = {
  position: 'fixed',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: 'rgba(0, 0, 0, 0.4)',
  backdropFilter: 'blur(4px)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  zIndex: 1000,
  padding: '20px'
};

const modalContentStyle: React.CSSProperties = {
  backgroundColor: 'var(--bg-surface)',
  border: '1px solid var(--border)',
  borderRadius: 'var(--radius)',
  padding: '24px 30px',
  maxWidth: '650px',
  width: '100%',
  maxHeight: '90vh',
  overflowY: 'auto',
  boxShadow: '0 20px 40px rgba(0, 0, 0, 0.12)',
  position: 'relative',
  display: 'flex',
  flexDirection: 'column'
};

const closeButtonStyle: React.CSSProperties = {
  position: 'absolute',
  top: '20px',
  right: '20px',
  background: 'none',
  border: 'none',
  fontSize: '18px',
  cursor: 'pointer',
  color: 'var(--text-secondary)',
  padding: '4px'
};

const sectionStyle: React.CSSProperties = {
  border: '1px solid var(--border)',
  borderRadius: 'var(--radius-sm)',
  padding: '16px',
  background: 'var(--bg-surface)'
};

const sectionTitleStyle: React.CSSProperties = {
  fontSize: '13px',
  fontWeight: 700,
  marginBottom: '12px',
  color: 'var(--text-primary)',
  textTransform: 'uppercase',
  letterSpacing: '0.5px'
};

const labelStyle: React.CSSProperties = {
  display: 'block',
  fontSize: '11px',
  fontWeight: 600,
  color: 'var(--text-secondary)',
  marginBottom: '4px'
};

const grid2Style: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
  gap: '12px'
};

const grid3Style: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(130px, 1fr))',
  gap: '12px'
};

const grid4Style: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(4, 1fr)',
  gap: '10px'
};

const costRowStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  fontSize: '13px',
  marginBottom: '4px',
  color: 'var(--text-secondary)'
};
