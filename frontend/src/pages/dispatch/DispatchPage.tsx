import { StatCard } from '../../components/ui/StatCard';
import { Badge } from '../../components/ui/Badge';

// Mock data for presentation
const MOCK_DRIVERS = [
  { id: 'DRV-001', name: 'Budi Santoso', vehicle: 'Van (B 1234 CD)', status: 'ON_ROUTE', assignments: 12, location: 'Kebayoran, Jakarta' },
  { id: 'DRV-002', name: 'Agus Pratama', vehicle: 'Motor (D 5678 EF)', status: 'IDLE', assignments: 0, location: 'Hub Bandung' },
  { id: 'DRV-003', name: 'Siti Aminah', vehicle: 'Van (L 9012 GH)', status: 'ON_ROUTE', assignments: 8, location: 'Gubeng, Surabaya' },
  { id: 'DRV-004', name: 'Reza Rahadian', vehicle: 'Truck (B 3456 JK)', status: 'MAINTENANCE', assignments: 0, location: 'Pool Jakarta' },
];

export default function DispatchPage() {
  return (
    <div className="flex flex-col gap-5">
      <div className="page-header !mb-0">
        <h1 className="text-2xl font-bold">🚚 Dispatch & Armada</h1>
        <p className="text-slate-500 mt-1">Manajemen armada, penugasan kurir, dan pemantauan rute</p>
      </div>

      {/* Bento Stats - Reusable Components */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard 
          title="Total Armada Aktif" 
          value="124" 
          subtitle="92% Utility Rate" 
          color="default" 
        />
        <StatCard 
          title="Sedang Di Rute (On Route)" 
          value="86" 
          subtitle="Tersebar di 5 Region" 
          color="primary" 
        />
        <StatCard 
          title="Perlu Maintenance" 
          value="5" 
          subtitle="Jadwal servis minggu ini" 
          color="warning" 
          className="bg-amber-50/50"
        />
      </div>

      {/* Table Section */}
      <div className="card flex flex-col gap-4">
        <div className="flex justify-between items-center">
          <h3 className="text-base font-semibold">Status Kurir Real-time</h3>
          <div className="flex gap-2">
            <button className="btn btn-ghost">Lihat Peta Armada</button>
            <button className="btn btn-primary">Tugaskan Kurir</button>
          </div>
        </div>

        <div className="table-wrapper">
          <table className="w-full">
            <thead>
              <tr>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">ID Kurir</th>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">Nama</th>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">Kendaraan</th>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">Lokasi Terakhir</th>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">Tugas Aktif</th>
                <th className="text-left text-xs font-semibold text-slate-400 uppercase tracking-wider py-3 border-b border-slate-200">Status</th>
              </tr>
            </thead>
            <tbody>
              {MOCK_DRIVERS.map(drv => (
                <tr key={drv.id} className="hover:bg-slate-50 transition-colors">
                  <td className="py-3 font-medium border-b border-slate-100">{drv.id}</td>
                  <td className="py-3 border-b border-slate-100">{drv.name}</td>
                  <td className="py-3 border-b border-slate-100">{drv.vehicle}</td>
                  <td className="py-3 border-b border-slate-100">{drv.location}</td>
                  <td className="py-3 border-b border-slate-100">{drv.assignments > 0 ? `${drv.assignments} paket` : '-'}</td>
                  <td className="py-3 border-b border-slate-100">
                    <Badge variant={
                      drv.status === 'ON_ROUTE' ? 'success' : 
                      drv.status === 'MAINTENANCE' ? 'danger' : 'default'
                    }>
                      {drv.status.replace('_', ' ')}
                    </Badge>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
