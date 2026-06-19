import type { TrackingHistory } from '../../types';
import TrackingMap from './TrackingMap';
import './TrackingResult.css';

interface Props {
  awb: string;
  data: TrackingHistory;
}

export default function TrackingResult({ awb, data }: Props) {
  const events = data.events || [];

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 md:p-8 animate-fade-in">
      <div className="border-b border-slate-100 pb-4 mb-6">
        <h2 className="text-sm font-bold text-slate-400 uppercase tracking-wider mb-1">Hasil Pencarian Resi</h2>
        <p className="text-2xl font-black text-slate-800">{awb}</p>
      </div>

      <div className="mb-8 rounded-2xl overflow-hidden border border-slate-200 shadow-sm relative z-0">
        <TrackingMap events={events} />
      </div>

      <div className="relative">
        {events.length === 0 ? (
          <p className="text-slate-500 text-center py-4">Belum ada riwayat perjalanan paket.</p>
        ) : (
          <div className="space-y-6 relative before:absolute before:inset-0 before:ml-5 before:-translate-x-px md:before:mx-auto md:before:translate-x-0 before:h-full before:w-0.5 before:bg-gradient-to-b before:from-transparent before:via-slate-200 before:to-transparent">
            {events.map((evt, idx) => (
              <div key={evt.id || idx} className="relative flex items-center justify-between md:justify-normal md:odd:flex-row-reverse group is-active">
                {/* Icon */}
                <div className="flex items-center justify-center w-10 h-10 rounded-full border-4 border-white bg-slate-200 text-slate-500 shadow shrink-0 md:order-1 md:group-odd:-translate-x-1/2 md:group-even:translate-x-1/2 z-10 group-first:bg-[#64965a] group-first:text-white transition-colors">
                  {idx === 0 ? '📍' : '✅'}
                </div>
                {/* Card */}
                <div className="w-[calc(100%-4rem)] md:w-[calc(50%-2.5rem)] p-4 rounded-xl border border-slate-100 bg-white shadow-sm group-hover:shadow-md transition-shadow group-first:border-[#64965a]/30 group-first:bg-[#64965a]/5">
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-1">
                    <h3 className="font-bold text-slate-800 text-sm">{evt.status} - {evt.location}</h3>
                    <time className="text-xs text-slate-400 font-medium">
                      {new Date(evt.timestamp).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })}
                    </time>
                  </div>
                  <p className="text-sm text-slate-600">{evt.description}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
