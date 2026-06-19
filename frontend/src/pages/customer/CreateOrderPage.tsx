import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';
import { useAuthStore } from '../../store/auth.store';

export default function CreateOrderPage() {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const [form, setForm] = useState({
    sender_name: user?.name || '',
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
    payment_provider: '',
    use_insurance: false
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const payload = {
        ...form,
        weight_actual: Number(form.weight_actual),
        length: Number(form.length),
        width: Number(form.width),
        height: Number(form.height),
      };
      
      const response = await apiClient.post(API_ENDPOINTS.order.base, payload);
      if (response.data && response.data.success) {
        alert('Order berhasil dibuat!');
        navigate('/customer/orders');
      } else {
        setError('Gagal membuat order.');
      }
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.error || 'Gagal membuat order. Periksa koneksi dan kelengkapan data.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto pb-12 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-3xl font-extrabold text-slate-800 mb-2">📦 Buat Order Baru</h1>
        <p className="text-slate-600">Lengkapi detail pengiriman paket Anda di bawah ini.</p>
      </div>

      {error && (
        <div className="bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded-r-xl mb-8">
          <p className="font-bold">Gagal memproses order</p>
          <p>{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* SENDER SECTION */}
        <div className="bg-white p-6 md:p-8 rounded-2xl shadow-sm border border-slate-200">
          <h2 className="text-xl font-bold text-slate-800 mb-6 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-[#64965a]/10 text-[#64965a] flex items-center justify-center">👤</span>
            Informasi Pengirim
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Nama Pengirim *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.sender_name} 
                onChange={e => setForm({...form, sender_name: e.target.value})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">No. Telepon *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.sender_phone} 
                onChange={e => setForm({...form, sender_phone: e.target.value})}
              />
            </div>
          </div>
          <div className="mb-6">
            <label className="block text-sm font-semibold text-slate-700 mb-2">Alamat Lengkap *</label>
            <textarea 
              required 
              className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all min-h-[100px]"
              value={form.sender_address} 
              onChange={e => setForm({...form, sender_address: e.target.value})}
            />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Kode Pos *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.origin_postal} 
                onChange={e => setForm({...form, origin_postal: e.target.value})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Kota</label>
              <input 
                type="text" 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.origin_city} 
                onChange={e => setForm({...form, origin_city: e.target.value})}
              />
            </div>
          </div>
        </div>

        {/* RECEIVER SECTION */}
        <div className="bg-white p-6 md:p-8 rounded-2xl shadow-sm border border-slate-200">
          <h2 className="text-xl font-bold text-slate-800 mb-6 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-[#64965a]/10 text-[#64965a] flex items-center justify-center">📍</span>
            Informasi Penerima
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Nama Penerima *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.receiver_name} 
                onChange={e => setForm({...form, receiver_name: e.target.value})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">No. Telepon *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.receiver_phone} 
                onChange={e => setForm({...form, receiver_phone: e.target.value})}
              />
            </div>
          </div>
          <div className="mb-6">
            <label className="block text-sm font-semibold text-slate-700 mb-2">Alamat Lengkap *</label>
            <textarea 
              required 
              className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all min-h-[100px]"
              value={form.receiver_address} 
              onChange={e => setForm({...form, receiver_address: e.target.value})}
            />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Kode Pos *</label>
              <input 
                type="text" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.dest_postal} 
                onChange={e => setForm({...form, dest_postal: e.target.value})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Kota</label>
              <input 
                type="text" 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.dest_city} 
                onChange={e => setForm({...form, dest_city: e.target.value})}
              />
            </div>
          </div>
        </div>

        {/* PACKAGE & SHIPPING DETAILS */}
        <div className="bg-white p-6 md:p-8 rounded-2xl shadow-sm border border-slate-200">
          <h2 className="text-xl font-bold text-slate-800 mb-6 flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-[#64965a]/10 text-[#64965a] flex items-center justify-center">📦</span>
            Rincian Paket & Pengiriman
          </h2>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Berat (kg) *</label>
              <input 
                type="number" step="0.1" min="0.1" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.weight_actual} 
                onChange={e => setForm({...form, weight_actual: parseFloat(e.target.value) || 0})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">P (cm) *</label>
              <input 
                type="number" min="1" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.length} 
                onChange={e => setForm({...form, length: parseInt(e.target.value) || 0})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">L (cm) *</label>
              <input 
                type="number" min="1" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.width} 
                onChange={e => setForm({...form, width: parseInt(e.target.value) || 0})}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">T (cm) *</label>
              <input 
                type="number" min="1" required 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.height} 
                onChange={e => setForm({...form, height: parseInt(e.target.value) || 0})}
              />
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Tipe Layanan *</label>
              <select 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.service_type}
                onChange={e => setForm({...form, service_type: e.target.value})}
              >
                <option value="REG">Regular (REG)</option>
                <option value="EXP">Express (EXP)</option>
                <option value="CARGO">Cargo (CARGO)</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-2">Metode Pembayaran *</label>
              <select 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                value={form.payment_type}
                onChange={e => {
                  const val = e.target.value;
                  setForm({...form, payment_type: val, payment_provider: val === 'COD' ? '' : ''});
                }}
              >
                <option value="COD">Cash On Delivery (COD)</option>
                <option value="TRANSFER">Transfer Bank</option>
                <option value="EWALLET">E-Wallet</option>
                <option value="VA">Virtual Account</option>
              </select>
            </div>
          </div>

          {form.payment_type !== 'COD' && (
            <div className="mb-6">
              <label className="block text-sm font-semibold text-slate-700 mb-2">Penyedia Pembayaran *</label>
              <select 
                className="w-full px-4 py-3 rounded-xl bg-slate-50 border border-slate-200 focus:bg-white focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
                required
                value={form.payment_provider}
                onChange={e => setForm({...form, payment_provider: e.target.value})}
              >
                <option value="">-- Pilih Penyedia --</option>
                {(form.payment_type === 'TRANSFER' || form.payment_type === 'VA') && (
                  <>
                    <option value="BCA">BCA</option>
                    <option value="MANDIRI">Mandiri</option>
                    <option value="BNI">BNI</option>
                    <option value="BRI">BRI</option>
                  </>
                )}
                {form.payment_type === 'EWALLET' && (
                  <>
                    <option value="GOPAY">GoPay</option>
                    <option value="OVO">OVO</option>
                    <option value="DANA">Dana</option>
                    <option value="LINKAJA">LinkAja</option>
                  </>
                )}
              </select>
            </div>
          )}

          <div className="flex items-center gap-3 p-4 bg-slate-50 rounded-xl border border-slate-200">
            <input 
              type="checkbox" id="use_insurance" 
              checked={form.use_insurance}
              onChange={e => setForm({...form, use_insurance: e.target.checked})}
              className="w-5 h-5 text-[#64965a] focus:ring-[#64965a] rounded border-slate-300"
            />
            <label htmlFor="use_insurance" className="font-medium text-slate-700 cursor-pointer select-none">
              Lindungi paket dengan Asuransi Pengiriman
            </label>
          </div>
        </div>

        {/* ACTIONS */}
        <div className="flex gap-4 justify-end">
          <button 
            type="button" 
            onClick={() => navigate('/customer/orders')}
            className="px-6 py-3 font-semibold text-slate-600 bg-white border border-slate-300 rounded-xl hover:bg-slate-50 transition-colors"
            disabled={loading}
          >
            Batal
          </button>
          <button 
            type="submit" 
            className="px-8 py-3 font-bold text-white bg-[#64965a] rounded-xl hover:bg-[#53804a] shadow-lg hover:shadow-xl transition-all disabled:opacity-70 disabled:cursor-not-allowed flex items-center"
            disabled={loading}
          >
            {loading ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Memproses...
              </>
            ) : 'Buat Order Sekarang'}
          </button>
        </div>
      </form>
    </div>
  );
}
