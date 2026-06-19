import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';

export interface BackendOrder {
  order_id: string;
  awb_number: string;
  sender_name: string;
  receiver_name: string;
  dest_city?: string;
  dest_postal_code?: string;
  origin_city?: string;
  origin_postal_code?: string;
  status: string;
  total_cost: number;
  service_id?: string;
}

export default function CustomerOrderList() {
  const navigate = useNavigate();
  const [orders, setOrders] = useState<BackendOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  const fetchOrders = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await apiClient.get(API_ENDPOINTS.order.base);
      if (response.data && response.data.success) {
        setOrders(response.data.data.orders || []);
      } else {
        setError('Format respons server tidak valid.');
      }
    } catch (err: any) {
      console.error(err);
      setError('Gagal memuat data order.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const filteredOrders = orders.filter(order => {
    const s = search.toLowerCase();
    return (
      order.awb_number.toLowerCase().includes(s) ||
      order.receiver_name.toLowerCase().includes(s) ||
      order.dest_city?.toLowerCase().includes(s)
    );
  });

  const getStatusDisplay = (status: string) => {
    switch (status) {
      case 'ORDER_CREATED': return { label: 'Order Dibuat', color: 'bg-slate-100 text-slate-700', icon: '📝' };
      case 'PAYMENT_PENDING': return { label: 'Menunggu Pembayaran', color: 'bg-orange-100 text-orange-800', icon: '💳' };
      case 'PAYMENT_CONFIRMED': return { label: 'Pembayaran Lunas', color: 'bg-emerald-100 text-emerald-800', icon: '✅' };
      case 'PICKED_UP': return { label: 'Paket Dijemput', color: 'bg-blue-100 text-blue-800', icon: '📦' };
      case 'ON_TRANSIT': return { label: 'Sedang Transit', color: 'bg-indigo-100 text-indigo-800', icon: '🚛' };
      case 'OUT_FOR_DELIVERY': return { label: 'Dalam Pengiriman', color: 'bg-yellow-100 text-yellow-800', icon: '🛵' };
      case 'DELIVERED': return { label: 'Selesai Terkirim', color: 'bg-green-100 text-green-800', icon: '🎉' };
      case 'FAILED': return { label: 'Pengiriman Gagal', color: 'bg-red-100 text-red-800', icon: '❌' };
      case 'RETURNED': return { label: 'Dikembalikan', color: 'bg-slate-200 text-slate-800', icon: '↩️' };
      default: return { label: status, color: 'bg-slate-100 text-slate-700', icon: '📌' };
    }
  };

  return (
    <div className="max-w-6xl mx-auto pb-12 animate-fade-in">
      {/* HEADER CENTERED & LARGE */}
      <div className="text-center mb-10">
        <h1 className="text-4xl font-extrabold text-slate-900 mb-4">📦 Manajemen Order</h1>
        <p className="text-lg text-slate-600 max-w-2xl mx-auto">
          Pantau riwayat dan status seluruh pengiriman paket Anda. Buat pesanan baru dengan mudah dan cepat.
        </p>
      </div>

      <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
        <div className="relative w-full md:w-96">
          <input 
            type="text" 
            placeholder="Cari Resi, Nama Penerima, atau Kota Tujuan..." 
            className="w-full pl-12 pr-4 py-3 bg-white border border-slate-200 rounded-xl shadow-sm focus:ring-2 focus:ring-[#64965a] focus:border-[#64965a] outline-none transition-all"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 text-lg">🔍</span>
        </div>
        
        <button 
          onClick={() => navigate('/customer/orders/create')}
          className="w-full md:w-auto px-6 py-3 bg-[#64965a] text-white font-bold rounded-xl shadow-lg hover:shadow-xl hover:bg-[#53804a] transition-all flex items-center justify-center gap-2"
        >
          <span>➕</span> Buat Order Baru
        </button>
      </div>

      {loading ? (
        <div className="text-center py-20">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#64965a] mx-auto mb-4"></div>
          <p className="text-slate-500 font-medium">Memuat data order...</p>
        </div>
      ) : error ? (
        <div className="bg-red-50 text-red-700 p-6 rounded-2xl text-center border border-red-100">
          <p className="text-2xl mb-2">⚠️</p>
          <p className="font-bold">{error}</p>
          <button onClick={fetchOrders} className="mt-4 px-4 py-2 bg-white text-red-600 rounded-lg shadow-sm hover:bg-red-50 text-sm font-semibold transition-colors">
            Coba Lagi
          </button>
        </div>
      ) : filteredOrders.length === 0 ? (
        <div className="bg-white rounded-3xl border border-slate-100 p-12 text-center shadow-sm">
          <div className="text-6xl mb-4">📭</div>
          <h3 className="text-2xl font-bold text-slate-800 mb-2">Belum ada Order</h3>
          <p className="text-slate-500 mb-6">Anda belum pernah membuat order pengiriman. Yuk buat order pertamamu sekarang!</p>
          <button 
            onClick={() => navigate('/customer/orders/create')}
            className="px-6 py-3 bg-[#64965a]/10 text-[#64965a] font-bold rounded-xl hover:bg-[#64965a]/20 transition-colors"
          >
            Buat Order Sekarang
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredOrders.map(order => {
            const statusInfo = getStatusDisplay(order.status);
            return (
              <div key={order.order_id} className="bg-white rounded-2xl border border-slate-200 p-6 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300 relative overflow-hidden group cursor-pointer">
                {/* Decorative background element */}
                <div className="absolute -right-6 -top-6 w-24 h-24 bg-slate-50 rounded-full group-hover:bg-[#64965a]/5 transition-colors -z-10"></div>
                
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <span className="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">NO. RESI</span>
                    <code className="text-slate-800 font-bold bg-slate-100 px-2 py-1 rounded text-sm group-hover:bg-[#64965a]/10 group-hover:text-[#64965a] transition-colors">
                      {order.awb_number}
                    </code>
                  </div>
                  <span className={`px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1 ${statusInfo.color}`}>
                    <span>{statusInfo.icon}</span> {statusInfo.label}
                  </span>
                </div>

                <div className="space-y-3 mb-5 border-y border-slate-100 py-4">
                  <div className="flex items-start gap-3">
                    <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-400 shrink-0 text-xs">📍</div>
                    <div>
                      <p className="text-xs text-slate-500 font-medium">Pengirim</p>
                      <p className="text-sm font-semibold text-slate-800">{order.sender_name}</p>
                      <p className="text-xs text-slate-500">{order.origin_city || order.origin_postal_code}</p>
                    </div>
                  </div>
                  <div className="w-0.5 h-4 bg-slate-200 ml-4"></div>
                  <div className="flex items-start gap-3">
                    <div className="w-8 h-8 rounded-full bg-[#64965a]/10 text-[#64965a] flex items-center justify-center shrink-0 text-xs">🏁</div>
                    <div>
                      <p className="text-xs text-[#64965a] font-bold">Penerima</p>
                      <p className="text-sm font-semibold text-slate-800">{order.receiver_name}</p>
                      <p className="text-xs text-slate-500">{order.dest_city || order.dest_postal_code}</p>
                    </div>
                  </div>
                </div>

                <div className="flex justify-between items-end">
                  <div>
                    <span className="text-xs font-bold text-slate-400 uppercase block mb-0.5">Total Biaya</span>
                    <span className="text-lg font-black text-slate-800">
                      Rp {order.total_cost.toLocaleString()}
                    </span>
                  </div>
                  <div className="text-right">
                    <span className="text-xs text-slate-500 block">Layanan</span>
                    <span className="text-sm font-bold text-[#64965a]">Buroqet {order.service_id === 'EXP' ? 'Express' : order.service_id === 'CARGO' ? 'Cargo' : 'Regular'}</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
