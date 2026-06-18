import { useEffect } from 'react';
import {
  MapContainer,
  TileLayer,
  Marker,
  Polyline,
  Popup,
  useMap,
} from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { TrackingEvent } from '../../types';
import './TrackingMap.css';

// ─── Fix default marker icon path yang rusak di Vite ───────────────────────
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

delete (L.Icon.Default.prototype as unknown as Record<string, unknown>)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerIcon2x,
  shadowUrl: markerShadow,
});

// ─── Warna marker berdasarkan status ──────────────────────────────────────
const STATUS_COLORS: Record<string, string> = {
  INBOUND:          '#B2D1A6', // light sage
  ON_TRANSIT:       '#9CC38F', // medium sage
  AT_HUB:           '#8BB67D', // darker sage
  OUT_FOR_DELIVERY: '#79AE6F', // primary green
  DELIVERED:        '#5A8D50', // dark green
  FAILED:           '#ef4444', // keep red only for actual failures
  RETURNED:         '#A0AAB5', // grey
};

const STATUS_ICONS: Record<string, string> = {
  INBOUND:          '📥',
  ON_TRANSIT:       '🚚',
  AT_HUB:           '🏢',
  OUT_FOR_DELIVERY: '🛵',
  DELIVERED:        '✅',
  FAILED:           '❌',
  RETURNED:         '↩️',
};

const STATUS_LABELS: Record<string, string> = {
  INBOUND:          'Diterima di Gudang',
  ON_TRANSIT:       'Dalam Perjalanan',
  AT_HUB:           'Tiba di Hub',
  OUT_FOR_DELIVERY: 'Sedang Diantar',
  DELIVERED:        'Terkirim',
  FAILED:           'Pengiriman Gagal',
  RETURNED:         'Dikembalikan',
};

// ─── Buat custom circular SVG marker ──────────────────────────────────────
function makeIcon(status: string, isLatest: boolean) {
  const color = STATUS_COLORS[status] ?? '#9898b0';
  const emoji = STATUS_ICONS[status] ?? '📦';
  const size = isLatest ? 42 : 32;
  const pulse = isLatest
    ? `<circle cx="21" cy="21" r="20" fill="${color}" opacity="0.2">
         <animate attributeName="r" from="18" to="24" dur="1.5s" repeatCount="indefinite"/>
         <animate attributeName="opacity" from="0.3" to="0" dur="1.5s" repeatCount="indefinite"/>
       </circle>`
    : '';

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${size}" height="${size}" viewBox="0 0 42 42">
      ${pulse}
      <circle cx="21" cy="21" r="17" fill="${color}" stroke="white" stroke-width="3"/>
      <text x="21" y="27" text-anchor="middle" font-size="16">${emoji}</text>
    </svg>`;

  return L.divIcon({
    html: svg,
    className: '',
    iconSize: [size, size],
    iconAnchor: [size / 2, size / 2],
    popupAnchor: [0, -(size / 2 + 4)],
  });
}

// ─── Auto-fit map bounds ke semua marker ──────────────────────────────────
function FitBounds({ positions }: { positions: [number, number][] }) {
  const map = useMap();
  useEffect(() => {
    if (positions.length === 0) return;
    if (positions.length === 1) {
      map.setView(positions[0], 13, { animate: true });
    } else {
      map.fitBounds(positions, { padding: [50, 50], animate: true });
    }
  }, [map, positions]);
  return null;
}

// ─── Format tanggal ───────────────────────────────────────────────────────
function fmtTime(ts: string) {
  const d = new Date(ts);
  return d.toLocaleString('id-ID', {
    weekday: 'short', day: 'numeric', month: 'short',
    hour: '2-digit', minute: '2-digit',
  });
}

// ─── Props ────────────────────────────────────────────────────────────────
interface Props {
  events: TrackingEvent[];
}

export default function TrackingMap({ events }: Props) {
  // Ambil hanya events yang punya koordinat valid
  const geoEvents = events.filter(ev => ev.latitude !== 0 && ev.longitude !== 0);

  if (geoEvents.length === 0) {
    return (
      <div className="map-no-coords">
        <span>🗺️</span>
        <p>Data koordinat belum tersedia untuk resi ini.<br/>
        Koordinat ditambahkan secara otomatis untuk pengiriman baru.</p>
      </div>
    );
  }

  const positions: [number, number][] = geoEvents.map(ev => [ev.latitude, ev.longitude]);
  const latestIdx = geoEvents.length - 1;

  // Warna polyline mengikuti status terakhir
  const routeColor = STATUS_COLORS[geoEvents[latestIdx].status] ?? '#79AE6F';

  return (
    <div className="tracking-map-wrap">
      <MapContainer
        center={positions[0]}
        zoom={10}
        className="leaflet-map"
        zoomControl={true}
        scrollWheelZoom={true}
      >
        {/* ── OpenStreetMap tile layer (light bento theme) ── */}
        <TileLayer
          url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png"
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/attributions">CARTO</a>'
          maxZoom={19}
        />

        {/* ── Route polyline ── */}
        {positions.length > 1 && (
          <Polyline
            positions={positions}
            pathOptions={{
              color: routeColor,
              weight: 4,
              opacity: 0.85,
              dashArray: '10, 6',
              lineCap: 'round',
              lineJoin: 'round',
            }}
          />
        )}

        {/* ── Markers untuk setiap checkpoint ── */}
        {geoEvents.map((ev, idx) => (
          <Marker
            key={ev.id}
            position={[ev.latitude, ev.longitude]}
            icon={makeIcon(ev.status, idx === latestIdx)}
            zIndexOffset={idx === latestIdx ? 1000 : 0}
          >
            <Popup className="tracking-popup">
              <div className="popup-inner">
                <div className="popup-status" style={{ color: STATUS_COLORS[ev.status] ?? '#fff' }}>
                  {STATUS_ICONS[ev.status] ?? '📦'} {STATUS_LABELS[ev.status] ?? ev.status}
                </div>
                <div className="popup-location">📍 {ev.location}</div>
                {ev.hub_id && (
                  <div className="popup-hub">Hub: {ev.hub_id}</div>
                )}
                {ev.description && (
                  <div className="popup-desc">{ev.description}</div>
                )}
                <div className="popup-time">{fmtTime(ev.timestamp)}</div>
                {ev.source && (
                  <div className="popup-source">Sumber: {ev.source}</div>
                )}
              </div>
            </Popup>
          </Marker>
        ))}

        {/* ── Auto-fit ── */}
        <FitBounds positions={positions} />
      </MapContainer>

      {/* ── Legend ── */}
      <div className="map-legend">
        <div className="legend-item">
          <span className="legend-dot" style={{ background: '#B2D1A6' }} /> Asal
        </div>
        {positions.length > 2 && (
          <div className="legend-item">
            <span className="legend-dot" style={{ background: '#8BB67D' }} /> Transit
          </div>
        )}
        <div className="legend-item">
          <span className="legend-dot legend-pulse" style={{ background: routeColor }} /> Posisi Terakhir
        </div>
      </div>
    </div>
  );
}
