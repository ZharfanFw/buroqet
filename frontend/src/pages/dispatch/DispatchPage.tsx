import { useState } from 'react';
import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';
import { StatCard } from '../../components/ui/StatCard';
import { Badge } from '../../components/ui/Badge';

interface Courier {
  id: string;
  name: string;
  current_location: {
    longitude: number;
    latitude: number;
  };
  status: string;
}

interface AssignResult {
  courier: Courier;
  distance_meters: number;
}

export default function DispatchPage() {
  // Stats
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Assign Form
  const [assignForm, setAssignForm] = useState({
    pickupLat: -6.175392, // Default Monas Jakarta
    pickupLon: 106.827153,
    radiusMeters: 5000 // 5 KM
  });
  const [assignedCourier, setAssignedCourier] = useState<AssignResult | null>(null);

  // Start Delivery Form
  const [deliveryForm, setDeliveryForm] = useState({
    awb: '',
    location: 'Hub Jakarta Pusat'
  });

  const handleAssignChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setAssignForm({
      ...assignForm,
      [e.target.name]: parseFloat(e.target.value) || 0
    });
  };

  const handleDeliveryChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setDeliveryForm({
      ...deliveryForm,
      [e.target.name]: e.target.value
    });
  };

  const handleAssignCourier = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);
      setSuccessMsg(null);
      setAssignedCourier(null);

      const res = await apiClient.post<AssignResult>(API_ENDPOINTS.dispatch.assign, {
        pickup_lat: assignForm.pickupLat,
        pickup_lon: assignForm.pickupLon,
        radius_meters: assignForm.radiusMeters
      });

      setAssignedCourier(res.data);
      setSuccessMsg(`Kurir ${res.data.courier.name} (${res.data.courier.id}) berhasil ditugaskan!`);
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.error || 'Tidak ada kurir aktif ditemukan di dalam radius pencarian.');
    } finally {
      setLoading(false);
    }
  };

  const handleStartDelivery = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!deliveryForm.awb) {
      setError('AWB harus diisi untuk memulai pengiriman.');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setSuccessMsg(null);

      const res = await apiClient.post<{ message: string }>(API_ENDPOINTS.dispatch.startDelivery, {
        awb: deliveryForm.awb,
        location: deliveryForm.location
      });

      setSuccessMsg(res.data.message || 'Pengiriman dimulai, event package.dispatched berhasil dipublish!');
      setDeliveryForm({ ...deliveryForm, awb: '' });
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.error || 'Gagal memulai pengiriman paket.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col gap-5">
      <div className="page-header !mb-0">
        <h1 className="text-2xl font-bold">🚚 Dispatch & Armada</h1>
        <p className="text-slate-500 mt-1">Manajemen armada, penugasan kurir terdekat via PostGIS, dan trigger pengiriman</p>
      </div>

      {successMsg && (
        <div style={{ padding: '12px 16px', background: '#e6f4ea', color: '#137333', borderRadius: '8px', fontSize: '14px', fontWeight: 500 }}>
          {successMsg}
        </div>
      )}

      {error && (
        <div style={{ padding: '12px 16px', background: '#fce8e6', color: '#c5221f', borderRadius: '8px', fontSize: '14px', fontWeight: 500 }}>
          {error}
        </div>
      )}

      {/* Bento Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard 
          title="Engine Lokasi" 
          value="PostGIS" 
          subtitle="ST_DWithin & ST_Distance" 
          color="primary" 
        />
        <StatCard 
          title="Radius Penugasan" 
          value={`${assignForm.radiusMeters / 1000} KM`} 
          subtitle="Dapat diatur dinamis" 
          color="default" 
        />
        <StatCard 
          title="Status Kafka" 
          value="Connected" 
          subtitle="Listening to order.created" 
          color="success" 
        />
      </div>

      {/* Main content grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
        
        {/* Card 1: Assign Courier */}
        <div className="card flex flex-col gap-4">
          <h3 className="text-base font-semibold">📍 Tugaskan Kurir Terdekat</h3>
          <p className="text-xs text-slate-500">Mencari kurir &quot;available&quot; terdekat dari koordinat gudang/toko menggunakan database spatial.</p>
          
          <form onSubmit={handleAssignCourier} className="flex flex-col gap-3">
            <div className="grid grid-cols-2 gap-3">
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-slate-400">Pickup Latitude</label>
                <input 
                  type="number" 
                  step="any"
                  name="pickupLat" 
                  value={assignForm.pickupLat} 
                  onChange={handleAssignChange} 
                  className="p-2 border border-slate-200 rounded-md text-sm"
                  required
                />
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-slate-400">Pickup Longitude</label>
                <input 
                  type="number" 
                  step="any"
                  name="pickupLon" 
                  value={assignForm.pickupLon} 
                  onChange={handleAssignChange} 
                  className="p-2 border border-slate-200 rounded-md text-sm"
                  required
                />
              </div>
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-slate-400">Radius Pencarian (Meter)</label>
              <input 
                type="number" 
                name="radiusMeters" 
                value={assignForm.radiusMeters} 
                onChange={handleAssignChange} 
                className="p-2 border border-slate-200 rounded-md text-sm"
                required
              />
            </div>

            <button type="submit" className="btn btn-primary w-full mt-2" disabled={loading}>
              {loading ? 'Mencari...' : 'Cari & Tugaskan Kurir'}
            </button>
          </form>

          {assignedCourier && (
            <div className="bg-slate-50 p-4 rounded-md border border-slate-200 mt-2 flex flex-col gap-2">
              <h4 className="font-semibold text-sm text-slate-700">Hasil Penugasan Kurir:</h4>
              <div className="text-xs text-slate-600 flex flex-col gap-1">
                <div><strong>ID Kurir:</strong> {assignedCourier.courier.id}</div>
                <div><strong>Nama:</strong> {assignedCourier.courier.name}</div>
                <div><strong>Jarak:</strong> {Math.round(assignedCourier.distance_meters)} meter dari lokasi pickup</div>
                <div className="mt-1">
                  <strong>Status Kurir:</strong>{' '}
                  <Badge variant="success">{assignedCourier.courier.status}</Badge>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Card 2: Start Delivery */}
        <div className="card flex flex-col gap-4">
          <h3 className="text-base font-semibold">🚀 Mulai Pengiriman Paket (Outbound)</h3>
          <p className="text-xs text-slate-500">Mempublikasikan event <code>package.dispatched</code> ke Kafka untuk mengupdate status tracking kurir.</p>

          <form onSubmit={handleStartDelivery} className="flex flex-col gap-3">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-slate-400">Nomor Resi (AWB)</label>
              <input 
                type="text" 
                name="awb" 
                placeholder="Contoh: BQ-123456" 
                value={deliveryForm.awb} 
                onChange={handleDeliveryChange} 
                className="p-2 border border-slate-200 rounded-md text-sm"
                required
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-slate-400">Lokasi Hub Pemindai</label>
              <input 
                type="text" 
                name="location" 
                value={deliveryForm.location} 
                onChange={handleDeliveryChange} 
                className="p-2 border border-slate-200 rounded-md text-sm"
                required
              />
            </div>

            <button type="submit" className="btn btn-primary w-full mt-2" style={{ backgroundColor: 'var(--success)' }} disabled={loading}>
              {loading ? 'Memproses...' : 'Mulai Pengiriman'}
            </button>
          </form>
        </div>

      </div>
    </div>
  );
}
