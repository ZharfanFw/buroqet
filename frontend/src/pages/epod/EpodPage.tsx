import { useEffect, useState, useCallback } from "react";

const API_BASE_URL = "/api/epod";

interface Epod {
  id: string;
  awb: string;
  receiver: string;
  photo_url: string;
  signature_url?: string;
  lat: number;
  lng: number;
  status: "VERIFIED" | "PENDING_REVIEW" | "REJECTED";
  reject_reason?: string;
  created_at: string;
  verified_at?: string;
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return (
    d.toLocaleDateString("id-ID", { day: "numeric", month: "short" }) +
    ", " +
    d.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" })
  );
}

function badgeClass(status: string) {
  if (status === "VERIFIED") return "badge-success";
  if (status === "PENDING_REVIEW") return "badge-warning";
  return "badge-danger";
}

export default function EpodPage() {
  const [epods, setEpods] = useState<Epod[]>([]);
  const [stats, setStats] = useState<{ total: number; pending: number }>({
    total: 0,
    pending: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [showUpload, setShowUpload] = useState(false);

  const fetchEpods = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams();
      if (statusFilter) params.set("status", statusFilter);
      params.set("page_size", "50");

      const res = await fetch(`${API_BASE_URL}/epod?${params.toString()}`);
      if (!res.ok) throw new Error("Gagal memuat data ePOD");
      const json = await res.json();

      const data: Epod[] = json.data || [];
      setEpods(data);

      const pendingCount = data.filter(
        (e) => e.status === "PENDING_REVIEW",
      ).length;
      setStats({ total: json.total ?? data.length, pending: pendingCount });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal memuat data ePOD");
    } finally {
      setLoading(false);
    }
  }, [statusFilter]);

  useEffect(() => {
    fetchEpods();
  }, [fetchEpods]);

  const handleVerify = async (id: string, status: "VERIFIED" | "REJECTED") => {
    try {
      const res = await fetch(`${API_BASE_URL}/epod/${id}/verify`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });
      if (!res.ok) {
        const errBody = await res.json().catch(() => ({}));
        throw new Error(errBody.error || "Gagal memperbarui status");
      }
      fetchEpods();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Gagal memperbarui status");
    }
  };

  const filtered = epods.filter(
    (e) => !search || e.awb.toLowerCase().includes(search.toLowerCase()),
  );

  const verifiedRate =
    stats.total > 0
      ? (((stats.total - stats.pending) / stats.total) * 100).toFixed(1)
      : "0.0";

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
      <div
        className="page-header"
        style={{
          marginBottom: 0,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "flex-start",
        }}
      >
        <div>
          <h1>✍️ Electronic Proof of Delivery (ePOD)</h1>
          <p>
            Verifikasi bukti pengiriman berupa foto dan tanda tangan penerima
          </p>
        </div>
        <button className="btn btn-primary" onClick={() => setShowUpload(true)}>
          + Upload ePOD
        </button>
      </div>

      {error && (
        <div
          className="card"
          style={{
            padding: "16px",
            borderColor: "var(--danger)",
            color: "var(--danger)",
          }}
        >
          {error}. Pastikan backend berjalan di {API_BASE_URL}.
        </div>
      )}

      {/* Bento Stats */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
          gap: "16px",
        }}
      >
        <div
          className="card"
          style={{
            padding: "20px",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          <div
            style={{
              color: "var(--text-muted)",
              fontSize: "13px",
              fontWeight: 600,
            }}
          >
            Total ePOD
          </div>
          <div
            style={{
              fontSize: "32px",
              fontWeight: 700,
              color: "var(--text-primary)",
              marginTop: "8px",
            }}
          >
            {stats.total}
          </div>
          <div
            style={{
              fontSize: "12px",
              color: "var(--success)",
              marginTop: "4px",
            }}
          >
            Tersinkronisasi dari server
          </div>
        </div>
        <div
          className="card"
          style={{
            padding: "20px",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          <div
            style={{
              color: "var(--text-muted)",
              fontSize: "13px",
              fontWeight: 600,
            }}
          >
            Menunggu Verifikasi (Pending)
          </div>
          <div
            style={{
              fontSize: "32px",
              fontWeight: 700,
              color: "var(--warning)",
              marginTop: "8px",
            }}
          >
            {stats.pending}
          </div>
          <div
            style={{
              fontSize: "12px",
              color: "var(--text-secondary)",
              marginTop: "4px",
            }}
          >
            Harap ditinjau oleh Admin
          </div>
        </div>
        <div
          className="card"
          style={{
            padding: "20px",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          <div
            style={{
              color: "var(--text-muted)",
              fontSize: "13px",
              fontWeight: 600,
            }}
          >
            Tingkat Keberhasilan Verifikasi
          </div>
          <div
            style={{
              fontSize: "32px",
              fontWeight: 700,
              color: "var(--primary)",
              marginTop: "8px",
            }}
          >
            {verifiedRate}%
          </div>
          <div
            style={{
              fontSize: "12px",
              color: "var(--text-secondary)",
              marginTop: "4px",
            }}
          >
            Dari total data yang dimuat
          </div>
        </div>
      </div>

      {/* Gallery Section */}
      <div
        className="card"
        style={{ display: "flex", flexDirection: "column", gap: "16px" }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <h3 style={{ fontSize: "16px", fontWeight: 600 }}>
            Bukti Pengiriman
          </h3>
          <div style={{ display: "flex", gap: "10px" }}>
            <input
              type="text"
              className="input"
              placeholder="Cari Resi (AWB)..."
              style={{ width: "200px" }}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <select
              className="input"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <option value="">Semua Status</option>
              <option value="PENDING_REVIEW">Pending Review</option>
              <option value="VERIFIED">Verified</option>
              <option value="REJECTED">Rejected</option>
            </select>
          </div>
        </div>

        {loading ? (
          <div
            style={{
              padding: "40px",
              textAlign: "center",
              color: "var(--text-muted)",
            }}
          >
            Memuat data...
          </div>
        ) : filtered.length === 0 ? (
          <div
            style={{
              padding: "40px",
              textAlign: "center",
              color: "var(--text-muted)",
            }}
          >
            Belum ada ePOD yang cocok dengan filter ini.
          </div>
        ) : (
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fill, minmax(240px, 1fr))",
              gap: "16px",
              marginTop: "10px",
            }}
          >
            {filtered.map((pod) => (
              <div
                key={pod.id}
                style={{
                  border: "1px solid var(--border)",
                  borderRadius: "12px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    height: "140px",
                    background: "var(--bg-hover)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    position: "relative",
                  }}
                >
                  {pod.photo_url ? (
                    <img
                      src={`${API_BASE_URL}${pod.photo_url}`}
                      alt={`Bukti pengiriman ${pod.awb}`}
                      style={{
                        width: "100%",
                        height: "100%",
                        objectFit: "cover",
                      }}
                    />
                  ) : (
                    <span style={{ fontSize: "40px", opacity: 0.2 }}>📸</span>
                  )}
                  <div
                    style={{
                      position: "absolute",
                      bottom: "8px",
                      right: "8px",
                      background: "rgba(0,0,0,0.6)",
                      color: "white",
                      fontSize: "10px",
                      padding: "2px 6px",
                      borderRadius: "4px",
                    }}
                  >
                    Lat: {pod.lat?.toFixed(1)}, Lng: {pod.lng?.toFixed(1)}
                  </div>
                </div>
                <div style={{ padding: "12px" }}>
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "flex-start",
                      marginBottom: "8px",
                    }}
                  >
                    <code
                      style={{
                        fontSize: "11px",
                        fontWeight: 600,
                        color: "var(--primary-dark)",
                      }}
                    >
                      {pod.awb}
                    </code>
                    <span
                      className={`badge ${badgeClass(pod.status)}`}
                      style={{ fontSize: "9px", padding: "2px 6px" }}
                    >
                      {pod.status}
                    </span>
                  </div>
                  <div
                    style={{
                      fontSize: "13px",
                      fontWeight: 500,
                      marginBottom: "2px",
                    }}
                  >
                    {pod.receiver}
                  </div>
                  <div
                    style={{
                      fontSize: "11px",
                      color: "var(--text-muted)",
                      marginBottom: "10px",
                    }}
                  >
                    {formatDate(pod.created_at)}
                  </div>

                  {pod.status === "PENDING_REVIEW" && (
                    <div style={{ display: "flex", gap: "6px" }}>
                      <button
                        className="btn btn-ghost"
                        style={{
                          flex: 1,
                          fontSize: "11px",
                          padding: "6px",
                          color: "var(--success)",
                        }}
                        onClick={() => handleVerify(pod.id, "VERIFIED")}
                      >
                        Verifikasi
                      </button>
                      <button
                        className="btn btn-ghost"
                        style={{
                          flex: 1,
                          fontSize: "11px",
                          padding: "6px",
                          color: "var(--danger)",
                        }}
                        onClick={() => handleVerify(pod.id, "REJECTED")}
                      >
                        Tolak
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showUpload && (
        <UploadModal
          onClose={() => setShowUpload(false)}
          onUploaded={() => {
            setShowUpload(false);
            fetchEpods();
          }}
        />
      )}
    </div>
  );
}

function UploadModal({
  onClose,
  onUploaded,
}: {
  onClose: () => void;
  onUploaded: () => void;
}) {
  const [awb, setAwb] = useState("");
  const [receiver, setReceiver] = useState("");
  const [lat, setLat] = useState("");
  const [lng, setLng] = useState("");
  const [photo, setPhoto] = useState<File | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!photo) {
      setError("Foto bukti pengiriman wajib diunggah.");
      return;
    }
    setSubmitting(true);
    setError(null);

    const form = new FormData();
    form.append("awb", awb);
    form.append("receiver", receiver);
    form.append("lat", lat || "0");
    form.append("lng", lng || "0");
    form.append("photo", photo);

    try {
      const res = await fetch(`${API_BASE_URL}/upload`, {
        method: "POST",
        body: form,
      });
      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.error || "Gagal mengunggah ePOD");
      }
      onUploaded();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal mengunggah ePOD");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        background: "rgba(0,0,0,0.4)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 50,
      }}
      onClick={onClose}
    >
      <form
        className="card"
        style={{
          width: "420px",
          padding: "24px",
          display: "flex",
          flexDirection: "column",
          gap: "14px",
        }}
        onClick={(e) => e.stopPropagation()}
        onSubmit={handleSubmit}
      >
        <h3 style={{ fontSize: "16px", fontWeight: 600, marginBottom: "4px" }}>
          Upload ePOD Baru
        </h3>

        <label
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "var(--text-muted)",
          }}
        >
          Nomor AWB
          <input
            className="input"
            style={{ width: "100%", marginTop: "4px" }}
            required
            value={awb}
            onChange={(e) => setAwb(e.target.value)}
            placeholder="BQ-2024-JKT-001"
          />
        </label>

        <label
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "var(--text-muted)",
          }}
        >
          Nama Penerima
          <input
            className="input"
            style={{ width: "100%", marginTop: "4px" }}
            required
            value={receiver}
            onChange={(e) => setReceiver(e.target.value)}
            placeholder="Budi Santoso"
          />
        </label>

        <div style={{ display: "flex", gap: "10px" }}>
          <label
            style={{
              fontSize: "12px",
              fontWeight: 600,
              color: "var(--text-muted)",
              flex: 1,
            }}
          >
            Latitude
            <input
              className="input"
              style={{ width: "100%", marginTop: "4px" }}
              type="number"
              step="any"
              value={lat}
              onChange={(e) => setLat(e.target.value)}
              placeholder="-6.2"
            />
          </label>
          <label
            style={{
              fontSize: "12px",
              fontWeight: 600,
              color: "var(--text-muted)",
              flex: 1,
            }}
          >
            Longitude
            <input
              className="input"
              style={{ width: "100%", marginTop: "4px" }}
              type="number"
              step="any"
              value={lng}
              onChange={(e) => setLng(e.target.value)}
              placeholder="106.8"
            />
          </label>
        </div>

        <label
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "var(--text-muted)",
          }}
        >
          Foto Bukti Pengiriman
          <input
            type="file"
            accept="image/*"
            required
            style={{ display: "block", marginTop: "4px", fontSize: "13px" }}
            onChange={(e) => setPhoto(e.target.files?.[0] || null)}
          />
        </label>

        {error && (
          <div style={{ color: "var(--danger)", fontSize: "12px" }}>
            {error}
          </div>
        )}

        <div style={{ display: "flex", gap: "8px", marginTop: "8px" }}>
          <button
            type="button"
            className="btn btn-ghost"
            style={{ flex: 1 }}
            onClick={onClose}
          >
            Batal
          </button>
          <button
            type="submit"
            className="btn btn-primary"
            style={{ flex: 1 }}
            disabled={submitting}
          >
            {submitting ? "Mengunggah..." : "Upload"}
          </button>
        </div>
      </form>
    </div>
  );
}
