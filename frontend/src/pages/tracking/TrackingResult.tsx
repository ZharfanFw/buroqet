import type { TrackingHistory, TrackingEvent } from '../../types';
import TrackingMap from './TrackingMap';
import './TrackingResult.css';

interface Props {
  awb: string;
  data: TrackingHistory;
}

// Mapping status ke warna, icon, dan label Indonesia
const STATUS_META: Record<string, { label: string; icon: string; color: string; bgColor: string }> = {
  INBOUND: {
    label: 'Diterima di Gudang',
    icon: '📥',
    color: '#B2D1A6',
    bgColor: 'rgba(121, 174, 111, 0.08)',
  },
  ON_TRANSIT: {
    label: 'Dalam Perjalanan',
    icon: '🚚',
    color: '#9CC38F',
    bgColor: 'rgba(121, 174, 111, 0.1)',
  },
  AT_HUB: {
    label: 'Tiba di Hub',
    icon: '🏢',
    color: '#8BB67D',
    bgColor: 'rgba(121, 174, 111, 0.12)',
  },
  OUT_FOR_DELIVERY: {
    label: 'Sedang Diantar',
    icon: '🛵',
    color: '#79AE6F',
    bgColor: 'rgba(121, 174, 111, 0.15)',
  },
  DELIVERED: {
    label: 'Terkirim',
    icon: '✅',
    color: '#5A8D50',
    bgColor: 'rgba(90, 141, 80, 0.15)',
  },
  FAILED: {
    label: 'Pengiriman Gagal',
    icon: '❌',
    color: '#ef4444',
    bgColor: 'rgba(239,68,68,0.12)',
  },
  RETURNED: {
    label: 'Dikembalikan',
    icon: '↩️',
    color: '#A0AAB5',
    bgColor: 'rgba(160, 170, 181, 0.12)',
  },
};

const getStatusMeta = (status: string) =>
  STATUS_META[status] ?? { label: status, icon: '📦', color: '#9898b0', bgColor: 'rgba(152,152,176,0.1)' };

function formatDate(ts: string) {
  const d = new Date(ts);
  return {
    date: d.toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' }),
    time: d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }),
  };
}

// Progress step labels untuk stepper atas (seperti Shopee)
const PROGRESS_STEPS = [
  { status: 'INBOUND', label: 'Diterima', icon: '📥' },
  { status: 'ON_TRANSIT', label: 'Transit', icon: '🚚' },
  { status: 'AT_HUB', label: 'Di Hub', icon: '🏢' },
  { status: 'OUT_FOR_DELIVERY', label: 'Diantar', icon: '🛵' },
  { status: 'DELIVERED', label: 'Terkirim', icon: '✅' },
];

function getProgressIndex(currentStatus: string): number {
  const idx = PROGRESS_STEPS.findIndex(s => s.status === currentStatus);
  return idx === -1 ? 0 : idx;
}

export default function TrackingResult({ awb, data }: Props) {
  const events = [...data.events].reverse(); // Terbaru di atas
  const latestEvent = data.events[data.events.length - 1] as TrackingEvent | undefined;
  const currentStatus = latestEvent?.status ?? 'INBOUND';
  const statusMeta = getStatusMeta(currentStatus);
  const progressIdx = getProgressIndex(currentStatus);
  const isDelivered = currentStatus === 'DELIVERED';
  const isFailed = currentStatus === 'FAILED' || currentStatus === 'RETURNED';

  return (
    <div className="tracking-result">
      {/* ── Status Hero Card ── */}
      <div className="status-hero card" style={{ borderColor: statusMeta.color + '40' }}>
        <div className="status-hero-left">
          <div className="status-icon-wrap" style={{ background: statusMeta.bgColor }}>
            <span className="status-icon-big">{statusMeta.icon}</span>
          </div>
          <div>
            <div className="status-label" style={{ color: statusMeta.color }}>
              {statusMeta.label}
            </div>
            <div className="status-awb">Resi: <strong>{awb}</strong></div>
            {latestEvent && (
              <div className="status-location">
                📍 {latestEvent.location}
                {latestEvent.hub_id && <span className="hub-tag">{latestEvent.hub_id}</span>}
              </div>
            )}
          </div>
        </div>
        {latestEvent && (
          <div className="status-hero-right">
            {(() => {
              const { date, time } = formatDate(latestEvent.timestamp);
              return (
                <>
                  <div className="status-time-big">{time}</div>
                  <div className="status-date">{date}</div>
                </>
              );
            })()}
          </div>
        )}
      </div>

      {/* ── Progress Stepper (Shopee-style) ── */}
      {!isFailed && (
        <div className="progress-stepper card">
          {PROGRESS_STEPS.map((step, idx) => {
            const isDone = idx <= progressIdx;
            const isCurrent = idx === progressIdx;
            return (
              <div key={step.status} className={`step ${isDone ? 'done' : ''} ${isCurrent ? 'current' : ''}`}>
                <div className="step-icon-wrap">
                  <div
                    className="step-circle"
                    style={isDone ? { background: statusMeta.color, borderColor: statusMeta.color } : {}}
                  >
                    {isDone ? (isCurrent ? step.icon : '✓') : step.icon}
                  </div>
                  {idx < PROGRESS_STEPS.length - 1 && (
                    <div
                      className="step-line"
                      style={idx < progressIdx ? { background: statusMeta.color } : {}}
                    />
                  )}
                </div>
                <div className={`step-label ${isCurrent ? 'current-label' : ''}`}
                  style={isCurrent ? { color: statusMeta.color } : {}}
                >
                  {step.label}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* ── Real Map ── */}
      <div className="map-card card">
        <div className="map-header">
          <span>🗺️ Peta Perjalanan Paket</span>
          <span className="map-total">{data.total} checkpoint</span>
        </div>
        <TrackingMap events={data.events} />
      </div>

      {/* ── Timeline Events ── */}
      <div className="timeline-card card">
        <div className="timeline-header">
          <h3>📋 Riwayat Perjalanan</h3>
          <span className="timeline-count">{data.total} event</span>
        </div>

        <div className="timeline">
          {events.map((ev, idx) => {
            const meta = getStatusMeta(ev.status);
            const { date, time } = formatDate(ev.timestamp);
            const isFirst = idx === 0;
            const isLast = idx === events.length - 1;
            return (
              <div key={ev.id} className={`timeline-item ${isFirst ? 'latest' : ''}`}>
                {/* Line */}
                <div className="tl-line-wrap">
                  <div
                    className="tl-dot"
                    style={{
                      background: isFirst ? meta.color : 'var(--bg-card)',
                      borderColor: isFirst ? meta.color : 'var(--border)',
                      boxShadow: isFirst ? `0 0 0 4px ${meta.color}22` : 'none',
                    }}
                  >
                    {isFirst && <span className="tl-dot-icon">{meta.icon}</span>}
                  </div>
                  {!isLast && <div className="tl-connector" />}
                </div>

                {/* Content */}
                <div className={`tl-content ${isFirst ? 'tl-content-latest' : ''}`}>
                  <div className="tl-top">
                    <div className="tl-left">
                      <span
                        className="tl-status-badge"
                        style={{ background: meta.bgColor, color: meta.color }}
                      >
                        {meta.icon} {meta.label}
                      </span>
                      <div className="tl-location">
                        📍 {ev.location}
                        {ev.hub_id && <span className="hub-tag">{ev.hub_id}</span>}
                      </div>
                      {ev.description && (
                        <div className="tl-description">{ev.description}</div>
                      )}
                      {ev.source && (
                        <div className="tl-source">Sumber: {ev.source}</div>
                      )}
                    </div>
                    <div className="tl-right">
                      <div className="tl-time">{time}</div>
                      <div className="tl-date">{date}</div>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* ── Delivered Celebration ── */}
      {isDelivered && (
        <div className="delivered-banner card">
          <div className="delivered-icon">🎉</div>
          <div>
            <strong>Paket berhasil diterima!</strong>
            <p>Terima kasih telah menggunakan layanan Buroqet</p>
          </div>
        </div>
      )}
    </div>
  );
}
